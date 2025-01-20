package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/sdk-utils/utils/uevent"

	"github.com/osmosis-labs/osmosis/v15/osmoutils"
	gammtypes "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

const AttributeKeyTakerFee = "fees_swap"

// ChargeFeesFromPayer charges the specified taker fee from the payer's account and
// processes it according to the fee token's properties.
// Wrapper for ChargeFees that sends the fee to x/txfees in advance.
func (k Keeper) ChargeFeesFromPayer(
	ctx sdk.Context,
	payer sdk.AccAddress,
	takerFeeCoin sdk.Coin,
	beneficiary *sdk.AccAddress,
) error {
	if takerFeeCoin.IsZero() {
		// Nothing to charge
		return nil
	}
	// Charge the fee from the payer to x/txfees
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, payer, types.ModuleName, sdk.NewCoins(takerFeeCoin))
	if err != nil {
		return fmt.Errorf("send coins to txfees account: %w", err)
	}

	return k.ChargeFees(ctx, takerFeeCoin, beneficiary, payer.String())
}

// ChargeFees processes the specified taker fee according to the fee token's properties.
// The fee must be sent to the module account beforehand.
// Payer field if optional and is only used for the event.
//
// If a beneficiary is provided, half of the fee is sent to the beneficiary.
// The remaining fee is sent to the txfees module account.
// If the fee token is the base denomination, it is burned.
// If the fee token is a registered fee token, it is swapped to the base denomination and then burned.
// If the fee token is unknown, it is sent to the community pool.
func (k Keeper) ChargeFees(
	ctx sdk.Context,
	takerFee sdk.Coin,
	beneficiary *sdk.AccAddress,
	payer string, // optional, only used for the event
) error {
	if takerFee.IsZero() {
		// Nothing to charge
		return nil
	}

	// Swap the taker fee to the base denom
	baseDenomFee, communityPool, err := k.swapFeeToBaseDenom(ctx, takerFee)
	if err != nil {
		return fmt.Errorf("swap fee to base denom: %w", err)
	}

	// If the fee token is unknown or the swap is unsuccessful, the fee is sent to the community pool.
	if !communityPool.Empty() {
		// Send unknown fee tokens to the community pool
		err = k.communityPool.FundCommunityPool(ctx, communityPool, k.accountKeeper.GetModuleAddress(types.ModuleName))
		if err != nil {
			return fmt.Errorf("fund community pool: %w", err)
		}

		k.Logger(ctx).With("fee", communityPool.String()).Debug("Sent fees to the community pool.")

		err = uevent.EmitTypedEvent(ctx, &types.EventChargeFee{
			Payer:         payer,
			TakerFee:      communityPool.String(),
			CommunityPool: true,
		})
		if err != nil {
			k.Logger(ctx).Error("Failed to emit event", "event", "EventChargeFee", "error", err)
		}
		return nil
	}

	// Send 50% of the base denom fee to the beneficiary if presented
	var beneficiaryCoins sdk.Coins
	if beneficiary != nil {
		fee := baseDenomFee[0]
		// beneficiaryCoin = takerFeeCoin / 2
		// note that beneficiaryCoin * 2 != takerFeeCoin because of the integer division rounding
		beneficiaryCoins = sdk.Coins{{Denom: fee.Denom, Amount: fee.Amount.QuoRaw(2)}}
		// takerFeeCoin = takerFeeCoin - beneficiaryCoin
		baseDenomFee = baseDenomFee.Sub(beneficiaryCoins...)

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, *beneficiary, beneficiaryCoins)
		if err != nil {
			return fmt.Errorf("send coins from fee payer to beneficiary: %w", err)
		}
	}

	// Burn the remaining base denom fee
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, baseDenomFee)
	if err != nil {
		return fmt.Errorf("burn coins: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventChargeFee{
		Payer:              payer,
		TakerFee:           baseDenomFee.String(),
		Beneficiary:        ValueFromPtr(beneficiary).String(),
		BeneficiaryRevenue: beneficiaryCoins.String(),
	})
	if err != nil {
		k.Logger(ctx).Error("Failed to emit event", "event", "EventChargeFee", "error", err)
	}

	return nil
}

func ValueFromPtr[T any](ptr *T) (zero T) {
	if ptr == nil {
		return zero
	}
	return *ptr
}

// swapFeeToBaseDenom swaps the taker fee coin to the base denom.
// If the fee token is unknown, it is sent to the community pool.
// The fee must be sent to the txfees module account beforehand.
func (k Keeper) swapFeeToBaseDenom(
	ctx sdk.Context,
	takerFeeCoin sdk.Coin,
) (baseDenomFee, communityPoolFee sdk.Coins, err error) {
	baseDenom, err := k.GetBaseDenom(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get base denom: %w", err)
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	// The fee is already in the base denom
	if takerFeeCoin.Denom == baseDenom {
		return sdk.Coins{takerFeeCoin}, nil, nil
	}

	// Get a fee token for the coin
	feetoken, err := k.GetFeeToken(ctx, takerFeeCoin.Denom)
	if err != nil {
		return nil, sdk.Coins{takerFeeCoin}, nil
	}

	// Swap the coin to base denom
	var (
		tokenOutAmount = math.ZeroInt() // Token amount in base denom
		route          = []poolmanagertypes.SwapAmountInRoute{{
			PoolId:        feetoken.PoolID,
			TokenOutDenom: baseDenom,
		}}
	)

	err = osmoutils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		tokenOutAmount, err = k.poolManager.RouteExactAmountIn(ctx, moduleAddr, route, takerFeeCoin, math.ZeroInt())
		return err
	})
	if err != nil {
		k.Logger(ctx).Error("swap fee to base denom", "error", err)
		return nil, sdk.Coins{takerFeeCoin}, nil
	}

	// If the swap is successful, we add the taker fee attribute to the emitted event
	k.appendTakerFeeAttribute(ctx)

	return sdk.Coins{{Denom: baseDenom, Amount: tokenOutAmount}}, nil, nil
}

// appendTakerFeeAttribute modifies the last token swap event to include the taker fee attribute
func (k Keeper) appendTakerFeeAttribute(ctx sdk.Context) {
	// we can modify the event in place, directly on the emittedEvents slice
	emittedEvents := ctx.EventManager().Events()

	for i := len(emittedEvents) - 1; i >= 0; i-- {
		if emittedEvents[i].Type == gammtypes.TypeEvtTokenSwapped {
			ev := emittedEvents[i].AppendAttributes(sdk.NewAttribute(AttributeKeyTakerFee, "true"))
			emittedEvents[i] = ev
			break
		}
	}
}
