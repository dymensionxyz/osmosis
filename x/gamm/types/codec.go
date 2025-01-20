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
	cdc.RegisterConcrete(&MsgJoinPool{}, "dymensionxyz/dymension/gamm/JoinPool", nil)
	cdc.RegisterConcrete(&MsgExitPool{}, "dymensionxyz/dymension/gamm/ExitPool", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "dymensionxyz/dymension/gamm/SwapExactAmountIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "dymensionxyz/dymension/gamm/SwapExactAmountOut", nil)
	cdc.RegisterConcrete(&MsgJoinSwapExternAmountIn{}, "dymensionxyz/dymension/gamm/JoinSwapExternAmountIn", nil)
	cdc.RegisterConcrete(&MsgJoinSwapShareAmountOut{}, "dymensionxyz/dymension/gamm/JoinSwapShareAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapExternAmountOut{}, "dymensionxyz/dymension/gamm/ExitSwapExternAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapShareAmountIn{}, "dymensionxyz/dymension/gamm/ExitSwapShareAmountIn", nil)
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
