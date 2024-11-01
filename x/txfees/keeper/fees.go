package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/osmoutils"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

// ChargeFees charges the specified taker fee from the payer's account and
// processes it according to the fee token's properties.
// If a beneficiary is provided, half of the fee is sent to the beneficiary.
// The remaining fee is sent to the txfees module account.
// If the fee token is the base denomination, it is burned.
// If the fee token is a registered fee token, it is swapped to the base denomination and then burned.
// If the fee token is unknown, it is burned directly.
func (k Keeper) ChargeFees(
	ctx sdk.Context,
	payer sdk.AccAddress,
	takerFeeCoin sdk.Coin,
	beneficiary *sdk.AccAddress,
) error {
	if takerFeeCoin.Amount.IsZero() {
		// Nothing to charge
		return nil
	}

	// Send half of the fee to beneficiary if presented
	if beneficiary != nil {
		// beneficiaryCoin = takerFeeCoin / 2
		// note that beneficiaryCoin * 2 != takerFeeCoin because of the integer division rounding
		beneficiaryCoin := sdk.Coin{Denom: takerFeeCoin.Denom, Amount: takerFeeCoin.Amount.QuoRaw(2)}
		// takerFeeCoin = takerFeeCoin - beneficiaryCoin
		takerFeeCoin = takerFeeCoin.Sub(beneficiaryCoin)

		err := k.bankKeeper.SendCoins(ctx, payer, *beneficiary, sdk.NewCoins(beneficiaryCoin))
		if err != nil {
			return fmt.Errorf("send coins from fee payer to beneficiary: %w", err)
		}
	}

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, payer, types.ModuleName, sdk.NewCoins(takerFeeCoin))
	if err != nil {
		return fmt.Errorf("send coins to txfees account: %w", err)
	}

	return k.ChargeFeesFromModuleAcc(ctx, takerFeeCoin)
}

// ChargeFeesFromModuleAcc processes the specified taker fee according to the fee token's properties.
// The main difference from ChargeFees is that the specified fee must be already present in the module account.
// If the fee token is the base denomination, it is burned.
// If the fee token is a registered fee token, it is swapped to the base denomination and then burned.
// If the fee token is unknown, it is burned directly.
func (k Keeper) ChargeFeesFromModuleAcc(
	ctx sdk.Context,
	takerFeeCoin sdk.Coin,
) error {
	if takerFeeCoin.Amount.IsZero() {
		// Nothing to charge
		return nil
	}

	baseDenom, err := k.GetBaseDenom(ctx)
	if err != nil {
		return fmt.Errorf("get base denom: %w", err)
	}

	// If the coin is in base denom, just burn
	if takerFeeCoin.Denom == baseDenom {
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(takerFeeCoin))
		if err != nil {
			return fmt.Errorf("burn coins: %w", err)
		}
		return nil
	}

	// Get a fee token for the coin
	feetoken, err := k.GetFeeToken(ctx, takerFeeCoin.Denom)
	if err != nil {
		// This should never happen in practice
		k.Logger(ctx).Error("Unknown fee token", "denom", takerFeeCoin.Denom, "error", err)

		// Burn unknown fee tokens
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(takerFeeCoin))
		if err != nil {
			return fmt.Errorf("unknown fee token: burn coins: %w", err)
		}

		return nil
	}

	// Swap the coin to base denom
	var (
		tokenOutAmount = sdk.ZeroInt() // Token amount in base denom
		moduleAddr     = k.accountKeeper.GetModuleAddress(types.ModuleName)
		route          = []poolmanagertypes.SwapAmountInRoute{{
			PoolId:        feetoken.PoolID,
			TokenOutDenom: baseDenom,
		}}
	)
	err = osmoutils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		tokenOutAmount, err = k.poolManager.RouteExactAmountIn(ctx, moduleAddr, route, takerFeeCoin, sdk.ZeroInt())
		return err
	})
	if err != nil {
		k.Logger(ctx).Error("Failed to swap fee token to base token. Trying to burn the tokens", "denom", takerFeeCoin.Denom, "error", err)

		// Burn unknown fee tokens
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(takerFeeCoin))
		if err != nil {
			return fmt.Errorf("unknown fee token: burn coins: %w", err)
		}

		return nil
	}

	// Burn the coin swapped to base denom
	takerFeeBaseDenom := sdk.NewCoins(sdk.NewCoin(baseDenom, tokenOutAmount))
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, takerFeeBaseDenom)
	if err != nil {
		return fmt.Errorf("burn coins: %w", err)
	}

	return nil
}
