package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/sdk-utils/utils/uevent"

	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

var (
	errUnswappableFeeToken = fmt.Errorf("fee token cannot be swapped to base denom")
)

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
	takerFeeBaseDenom, err := k.swapFeeToBaseDenom(ctx, takerFee)
	// Send unknown fee tokens to the community pool
	if errors.Is(err, errUnswappableFeeToken) {
		k.Logger(ctx).With("fee", takerFee.String(), "payer", payer).
			Debug("Cannot swap fee to base denom. Sending it to the community pool.")

		err = k.communityPool.FundCommunityPool(ctx, sdk.NewCoins(takerFee), k.accountKeeper.GetModuleAddress(types.ModuleName))
		if err != nil {
			return fmt.Errorf("fund community pool: %w", err)
		}

		err = uevent.EmitTypedEvent(ctx, &types.EventChargeFee{
			Payer:         payer,
			TakerFee:      takerFee.String(),
			CommunityPool: true,
		})
		if err != nil {
			k.Logger(ctx).Error("emit event", "event", "EventChargeFee", "error", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("swap fee to base denom: %w", err)
	}

	// Nothing to charge after swapping (not supposed to happen actually)
	if takerFeeBaseDenom.IsNil() || takerFeeBaseDenom.IsZero() {
		k.Logger(ctx).With("fee", takerFee.String(), "payer", payer).
			Error("Fee after swapping to base denom is zero. Nothing to charge.")
		return nil
	}

	// Send 50% of the base denom fee to the beneficiary if presented
	beneficiaryFee := sdk.Coins{}
	if beneficiary != nil {
		// beneficiaryFee = takerFee / 2
		beneficiaryFee = sdk.Coins{sdk.NewCoin(takerFeeBaseDenom.Denom, takerFeeBaseDenom.Amount.QuoRaw(2))}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, *beneficiary, beneficiaryFee)
		if err != nil {
			return fmt.Errorf("send coins from fee payer to beneficiary: %w", err)
		}

		// update remaining taker fee
		takerFeeBaseDenom = takerFeeBaseDenom.Sub(beneficiaryFee[0])
	}

	// Burn the remaining base denom fee
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(takerFeeBaseDenom))
	if err != nil {
		return fmt.Errorf("burn coins: %w", err)
	}

	err = uevent.EmitTypedEvent(ctx, &types.EventChargeFee{
		Payer:              payer,
		TakerFee:           takerFee.String(),
		SwappedTakerFee:    takerFeeBaseDenom.String(),
		Beneficiary:        ValueFromPtr(beneficiary).String(),
		BeneficiaryRevenue: beneficiaryFee.String(),
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
// Returns error if the fee token is unknown or if swapping fails.
// The fee must be sent to the txfees module account beforehand.
func (k Keeper) swapFeeToBaseDenom(
	ctx sdk.Context,
	takerFeeCoin sdk.Coin,
) (sdk.Coin, error) {
	baseDenom, err := k.GetBaseDenom(ctx)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("get base denom: %w", err)
	}
	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)

	// The fee is already in the base denom
	if takerFeeCoin.Denom == baseDenom {
		return takerFeeCoin, nil
	}

	// Get a fee token for the coin
	feetoken, err := k.GetFeeToken(ctx, takerFeeCoin.Denom)
	if err != nil {
		return sdk.Coin{}, errUnswappableFeeToken
	}

	// Swap the coin to base denom
	var (
		route = []poolmanagertypes.SwapAmountInRoute{{
			PoolId:        feetoken.PoolID,
			TokenOutDenom: baseDenom,
		}}
	)

	cacheCtx, write := ctx.CacheContext()
	tokenOutAmount, err := k.poolManager.RouteExactAmountIn(cacheCtx, moduleAddr, route, takerFeeCoin, sdk.ZeroInt())
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("failed to swap fee token: %w", err)
	}
	cacheCtx.WithEventManager(sdk.NewEventManager()) // This clears any events that were emitted
	write()                                          // This writes the state changes but with cleared events

	return sdk.NewCoin(baseDenom, tokenOutAmount), nil
}
