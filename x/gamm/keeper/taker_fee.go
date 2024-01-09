package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

// ChargeTakerFee charges the taker fee to the sender
// If the taker fee coin is the base denom, send it to the txfees module
// If the taker fee coin is a registered fee token, send it to the txfees module
// If the taker fee coin is not supported, swap it to the base denom on the first pool, then send it to the txfees module
func (k Keeper) chargeTakerFee(ctx sdk.Context, takerFeeCoin sdk.Coin, sender sdk.AccAddress, route poolmanagertypes.SwapAmountInRoute) error {
	// Check if the taker fee coin is the base denom
	denom, err := k.txfeeKeeper.GetBaseDenom(ctx)
	if err != nil {
		return err
	}
	if takerFeeCoin.Denom == denom {
		return k.sendToTxFees(ctx, sender, takerFeeCoin)
	}

	// Check if the taker fee coin is a registered fee token
	_, err = k.txfeeKeeper.GetFeeToken(ctx, takerFeeCoin.Denom)
	if err == nil {
		return k.sendToTxFees(ctx, sender, takerFeeCoin)
	}

	// If not supported denom, swap on the first pool to get some pool base denom, which has liquidity with DYM
	ctx.Logger().Debug("taker fee coin is not supported by txfee module, requires swap", "takerFeeCoin", takerFeeCoin)
	swappedTakerFee, err := k.swapTakerFee(ctx, sender, route, takerFeeCoin)
	if err != nil {
		return err
	}

	return k.sendToTxFees(ctx, sender, swappedTakerFee)
}

// swapTakerFee swaps the taker fee coin to the base denom on the first pool
func (k Keeper) swapTakerFee(ctx sdk.Context, sender sdk.AccAddress, route poolmanagertypes.SwapAmountInRoute, tokenIn sdk.Coin) (sdk.Coin, error) {
	minAmountOut := sdk.ZeroInt()
	swapRoutes := poolmanagertypes.SwapAmountInRoutes{route}
	out, err := k.poolManager.RouteExactAmountIn(ctx, sender, swapRoutes, tokenIn, minAmountOut)
	if err != nil {
		return sdk.Coin{}, err
	}
	coin := sdk.NewCoin(route.TokenOutDenom, out)
	return coin, nil
}

// sendToTxFees sends the taker fee coin to the txfees module
func (k Keeper) sendToTxFees(ctx sdk.Context, sender sdk.AccAddress, takerFeeCoin sdk.Coin) error {
	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, txfeestypes.ModuleName, sdk.NewCoins(takerFeeCoin))
}

/* ---------------------------------- Utils --------------------------------- */
// Returns remaining amount in to swap, and takerFeeCoins.
// returns (1 - takerFee) * tokenIn, takerFee * tokenIn
func (k Keeper) SubTakerFee(tokenIn sdk.Coin, takerFee sdk.Dec) (sdk.Coin, sdk.Coin) {
	amountInAfterSubTakerFee := sdk.NewDecFromInt(tokenIn.Amount).MulTruncate(sdk.OneDec().Sub(takerFee))
	tokenInAfterSubTakerFee := sdk.NewCoin(tokenIn.Denom, amountInAfterSubTakerFee.TruncateInt())
	takerFeeCoin := sdk.NewCoin(tokenIn.Denom, tokenIn.Amount.Sub(tokenInAfterSubTakerFee.Amount))
	return tokenInAfterSubTakerFee, takerFeeCoin
}

// here we need the output to be (tokenIn / (1 - takerFee), takerFee * tokenIn)
func (k Keeper) AddTakerFee(tokenIn sdk.Coin, takerFee sdk.Dec) (sdk.Coin, sdk.Coin) {
	amountInAfterAddTakerFee := sdk.NewDecFromInt(tokenIn.Amount).Quo(sdk.OneDec().Sub(takerFee))
	tokenInAfterAddTakerFee := sdk.NewCoin(tokenIn.Denom, amountInAfterAddTakerFee.Ceil().TruncateInt())
	takerFeeCoin := sdk.NewCoin(tokenIn.Denom, tokenInAfterAddTakerFee.Amount.Sub(tokenIn.Amount))
	return tokenInAfterAddTakerFee, takerFeeCoin
}
