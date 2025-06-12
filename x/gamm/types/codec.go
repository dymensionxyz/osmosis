package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*CFMMPoolI)(nil), nil)
	cdc.RegisterConcrete(&MsgJoinPool{}, "dymension/gamm/JoinPool", nil)
	cdc.RegisterConcrete(&MsgExitPool{}, "dymension/gamm/ExitPool", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "dymension/gamm/SwapExactAmountIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "dymension/gamm/SwapExactAmountOut", nil)
	cdc.RegisterConcrete(&MsgJoinSwapExternAmountIn{}, "dymension/gamm/JoinExternAmountIn", nil)
	cdc.RegisterConcrete(&MsgJoinSwapShareAmountOut{}, "dymension/gamm/JoinShareAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapExternAmountOut{}, "dymension/gamm/ExitExternAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapShareAmountIn{}, "dymension/gamm/ExitShareAmountIn", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface(
		"osmosis.gamm.v1beta1.PoolI", // N.B.: the old proto-path is preserved for backwards-compatibility.
		(*CFMMPoolI)(nil),
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgJoinPool{},
		&MsgExitPool{},
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
		&MsgJoinSwapExternAmountIn{},
		&MsgJoinSwapShareAmountOut{},
		&MsgExitSwapExternAmountOut{},
		&MsgExitSwapShareAmountIn{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
