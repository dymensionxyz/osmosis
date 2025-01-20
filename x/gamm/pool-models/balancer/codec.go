package balancer

import (
	types "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	proto "github.com/cosmos/gogoproto/proto"
)

// RegisterLegacyAminoCodec registers the necessary x/gamm interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&Pool{}, "dymensionxyz/dymension/gamm/BalancerPool", nil)
	cdc.RegisterConcrete(&MsgCreateBalancerPool{}, "dymensionxyz/dymension/gamm/CreateBalancerPool", nil)
	cdc.RegisterConcrete(&PoolParams{}, "dymensionxyz/dymension/gamm/BalancerPoolParams", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	//Registering old proto-path for backwards-compatibility
	proto.RegisterType((*Pool)(nil), "osmosis.gamm.v1beta1.Pool")

	registry.RegisterInterface(
		"osmosis.poolmanager.v1beta1.PoolI",
		(*poolmanagertypes.PoolI)(nil),
		&Pool{},
	)
	registry.RegisterInterface(
		"osmosis.gamm.v1beta1.PoolI", // N.B.: the old proto-path is preserved for backwards-compatibility.
		(*types.CFMMPoolI)(nil),
		&Pool{},
	)
	registry.RegisterImplementations(
		(*proto.Message)(nil),
		&PoolParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
