package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ChargeTakerFee charges the taker fee to the sender
// If the taker fee coin is the base denom, send it to the txfees module
// If the taker fee coin is a registered fee token, send it to the txfees module
// If the taker fee coin is not supported, swap it to the base denom on the first pool, then send it to the txfees module
// Send some portion of the taker fee to the provided beneficiary
func (k Keeper) chargeTakerFee(
	ctx sdk.Context,
	takerFeeCoin sdk.Coin,
	sender sdk.AccAddress,
	beneficiary *sdk.AccAddress,
) error {
	if takerFeeCoin.IsZero() {
		return nil
	}
	err := k.txfeeKeeper.ChargeFeesFromPayer(ctx, sender, takerFeeCoin, beneficiary)
	if err != nil {
		return fmt.Errorf("charge fees: sender: %s: fee: %s: %w", sender, takerFeeCoin, err)
	}
	return nil

}

// While charging taker fee, we reward the owner of the rollapp involved in swap. In that case,
// the owner is called the beneficiary. The following cases are possible:
//
//	 No | In Denom    | Out Denom   | Result
//	----|-------------|-------------|------------------------------
//	 1  | RollApp     | RollApp     | Reward the IN RollApp owner
//	 2  | RollApp     | Non-RollApp | Reward the IN RollApp owner
//	 3  | Non-RollApp | RollApp     | Reward the OUT RollApp owner
//	 4  | Non-RollApp | Non-RollApp | No one is rewarded
//
// Return nil beneficiary address if no one is rewarded: case (4) or error.
func (k Keeper) getTakerFeeBeneficiary(ctx sdk.Context, inDenom, outDenom string) *sdk.AccAddress {
	// This keeper is set to nil in osmosis repo to avoid circular dependency.
	// Should be non-nil in the dymension repo.
	if k.rollappKeeper == nil {
		return nil
	}
	// First, try cases (1) and (2)
	ownerIn, errIn := k.rollappKeeper.GetRollappOwnerByDenom(ctx, inDenom)
	if errIn == nil {
		return &ownerIn
	}
	// Try case (3)
	ownerOut, errOut := k.rollappKeeper.GetRollappOwnerByDenom(ctx, outDenom)
	if errOut == nil {
		return &ownerOut
	}
	// Case (4) or error while parsing denoms
	ctx.Logger().With("in_denom", inDenom, "out_denom", outDenom, "parse_err_in", errIn, "parse_err_out", errOut).
		Debug("swap without beneficiary: either two non-rollapp tokens or error when determining beneficiary")
	return nil
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
