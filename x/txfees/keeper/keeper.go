package keeper

import (
	"fmt"

	"cosmossdk.io/log"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper types.AccountKeeper
	epochKeeper   types.EpochKeeper
	bankKeeper    types.BankKeeper
	poolManager   types.PoolManager
	gammKeeper    types.GAMMKeeper
	communityPool types.CommunityPoolKeeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	accountKeeper types.AccountKeeper,
	epochKeeper types.EpochKeeper,
	bankKeeper types.BankKeeper,
	poolManager types.PoolManager,
	spotPriceCalculator types.GAMMKeeper,
	communityPool types.CommunityPoolKeeper,
) Keeper {

	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		epochKeeper:   epochKeeper,
		poolManager:   poolManager,
		gammKeeper:    spotPriceCalculator,
		communityPool: communityPool,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) getFeeTokensStore(ctx sdk.Context) storetypes.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.FeeTokensStorePrefix)
}

// func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
// 	k.paramSpace.GetParamSet(ctx, &params)
// 	return params
// }

// // SetParams sets the total set of params.
// func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
// 	k.paramSpace.SetParamSet(ctx, &params)
// }

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ParamsKey)
	var params types.Params
	k.cdc.MustUnmarshal(b, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, b)
}
