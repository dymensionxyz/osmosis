package keeper

import (
	"fmt"

	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func permContains(perms []string, perm string) bool {
	for _, v := range perms {
		if v == perm {
			return true
		}
	}

	return false
}

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	paramSpace paramtypes.Subspace
	hooks      types.GammHooks

	// keepers
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	communityPoolKeeper types.CommunityPoolKeeper
	poolManager         types.PoolManager
	txfeeKeeper         types.TxFeeKeeper
	rollappKeeper       types.RollappKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	communityPoolKeeper types.CommunityPoolKeeper,
) Keeper {
	// Ensure that the module account are set.
	moduleAddr, perms := accountKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	if moduleAddr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	if !permContains(perms, authtypes.Minter) {
		panic(fmt.Sprintf("%s module account should have the minter permission", types.ModuleName))
	}
	if !permContains(perms, authtypes.Burner) {
		panic(fmt.Sprintf("%s module account should have the burner permission", types.ModuleName))
	}
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	if accountKeeper == nil {
		panic("account keeper is nil")
	}
	if bankKeeper == nil {
		panic("bank keeper is nil")
	}
	if communityPoolKeeper == nil {
		panic("community pool keeper is nil")
	}

	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramSpace: paramSpace,
		// keepers
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		communityPoolKeeper: communityPoolKeeper,
	}
}

// SetHooks sets the gamm hooks.
func (k *Keeper) SetHooks(gh types.GammHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set gamm hooks twice")
	}

	k.hooks = gh

	return k
}

// SetPoolManager sets the pool manager.
// must be called when initializing the keeper.
func (k *Keeper) SetPoolManager(poolManager types.PoolManager) {
	k.poolManager = poolManager
}

// SetTxFees sets the tx fees keeper.
// must be called when initializing the keeper.
func (k *Keeper) SetTxFees(txfees types.TxFeeKeeper) {
	k.txfeeKeeper = txfees
}

// SetRollapp sets the tx fees keeper.
// must be called when initializing the keeper.
func (k *Keeper) SetRollapp(rollapp types.RollappKeeper) {
	k.rollappKeeper = rollapp
}

// GetParams returns the total set params.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
