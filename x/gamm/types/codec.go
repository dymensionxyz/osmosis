package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*CFMMPoolI)(nil), nil)
	cdc.RegisterConcrete(&MsgJoinPool{}, "osmosis/gamm/JoinPool", nil)
	cdc.RegisterConcrete(&MsgExitPool{}, "osmosis/gamm/ExitPool", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountIn{}, "osmosis/gamm/SwapExactAmountIn", nil)
	cdc.RegisterConcrete(&MsgSwapExactAmountOut{}, "osmosis/gamm/SwapExactAmountOut", nil)
	cdc.RegisterConcrete(&MsgJoinSwapExternAmountIn{}, "osmosis/gamm/JoinSwapExternAmountIn", nil)
	cdc.RegisterConcrete(&MsgJoinSwapShareAmountOut{}, "osmosis/gamm/JoinSwapShareAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapExternAmountOut{}, "osmosis/gamm/ExitSwapExternAmountOut", nil)
	cdc.RegisterConcrete(&MsgExitSwapShareAmountIn{}, "osmosis/gamm/ExitSwapShareAmountIn", nil)
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

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	// Register all Amino interfaces and concrete types on the authz Amino codec so that this can later be
	// used to properly serialize MsgGrant and MsgExec instances
	sdk.RegisterLegacyAminoCodec(amino)
	RegisterLegacyAminoCodec(authzcodec.Amino)

	amino.Seal()
}
