package keeper

import (
	"fmt"

	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the txfees module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	recipientAcc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	if recipientAcc == nil {
		panic(fmt.Sprintf("module account %s does not exist", types.ModuleName))
	}

	k.SetParams(ctx, genState.Params)

	err := k.SetBaseDenom(ctx, genState.Basedenom)
	if err != nil {
		panic(err)
	}
	err = k.SetFeeTokens(ctx, genState.Feetokens)
	if err != nil {
		panic(err)
	}

	epochIdentifier := k.GetParams(ctx).EpochIdentifier
	info := k.epochKeeper.GetEpochInfo(ctx, epochIdentifier)
	if info.Identifier == "" {
		panic(fmt.Sprintf("epoch info for identifier %s does not exist", epochIdentifier))
	}
}

// ExportGenesis returns the txfees module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Basedenom, _ = k.GetBaseDenom(ctx)
	genesis.Feetokens = k.GetFeeTokens(ctx)
	genesis.Params = k.GetParams(ctx)
	return genesis
}
