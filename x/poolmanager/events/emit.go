package events

import (
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
)

func EmitSwapEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, input, output sdk.Coins, spotPrice, takerFee, swapFee math.LegacyDec) {
	ctx.EventManager().EmitEvents(sdk.Events{
		newSwapEvent(sender, poolId, input, output, spotPrice, takerFee, swapFee),
	})
}

func newSwapEvent(sender sdk.AccAddress, poolId uint64, input, output sdk.Coins, spotPrice, takerFee, swapFee math.LegacyDec) sdk.Event {
	return sdk.NewEvent(
		types.TypeEvtTokenSwapped,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(types.AttributeKeyTokensIn, input.String()),
		sdk.NewAttribute(types.AttributeKeyTokensOut, output.String()),
		sdk.NewAttribute(types.AttributeKeyClosingPrice, spotPrice.String()),
		sdk.NewAttribute(types.AttributeKeyTakerFee, takerFee.String()),
		sdk.NewAttribute(types.AttributeKeySwapFee, swapFee.String()),
	)
}

func EmitAddLiquidityEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		newAddLiquidityEvent(sender, poolId, liquidity),
	})
}

func newAddLiquidityEvent(sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		types.TypeEvtPoolJoined,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(types.AttributeKeyTokensIn, liquidity.String()),
	)
}

func EmitRemoveLiquidityEvent(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) {
	ctx.EventManager().EmitEvents(sdk.Events{
		newRemoveLiquidityEvent(sender, poolId, liquidity),
	})
}

func newRemoveLiquidityEvent(sender sdk.AccAddress, poolId uint64, liquidity sdk.Coins) sdk.Event {
	return sdk.NewEvent(
		types.TypeEvtPoolExited,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		sdk.NewAttribute(types.AttributeKeyPoolId, strconv.FormatUint(poolId, 10)),
		sdk.NewAttribute(types.AttributeKeyTokensOut, liquidity.String()),
	)
}
