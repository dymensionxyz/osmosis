package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

type Keeper struct {
	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	accountKeeper       types.AccountKeeper
	epochKeeper         types.EpochKeeper
	bankKeeper          types.BankKeeper
	poolManager         types.PoolManager
	spotPriceCalculator types.SpotPriceCalculator
	communityPool       types.CommunityPoolKeeper
}

var _ types.TxFeesKeeper = (*Keeper)(nil)

func NewKeeper(
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	epochKeeper types.EpochKeeper,
	bankKeeper types.BankKeeper,
	poolManager types.PoolManager,
	spotPriceCalculator types.SpotPriceCalculator,
	communityPool types.CommunityPoolKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:            storeKey,
		paramSpace:          paramSpace,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		epochKeeper:         epochKeeper,
		poolManager:         poolManager,
		spotPriceCalculator: spotPriceCalculator,
		communityPool:       communityPool,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetFeeTokensStore(ctx sdk.Context) sdk.KVStore {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.FeeTokensStorePrefix)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
