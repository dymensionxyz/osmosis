package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	pooltypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
	v2 "github.com/osmosis-labs/osmosis/v15/x/txfees/types/migrations/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	k Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		k: keeper,
	}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	baseDenom := m.k.MustGetBaseDenom(ctx)
	for _, oldFeeToken := range m.GetOldFeeTokens(ctx) {
		newFeeToken := types.FeeToken{
			Denom: oldFeeToken.Denom,
			Route: []pooltypes.SwapAmountInRoute{
				{
					PoolId:        oldFeeToken.PoolID,
					TokenOutDenom: baseDenom,
				},
			},
		}
		err := m.k.SetFeeToken(ctx, newFeeToken)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Migrator) GetOldFeeTokens(ctx sdk.Context) (feetokens []v2.FeeToken) {
	prefixStore := m.k.getFeeTokensStore(ctx)

	// this entire store just contains FeeTokens, so iterate over all entries.
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	feeTokens := []v2.FeeToken{}

	for ; iterator.Valid(); iterator.Next() {
		feeToken := v2.FeeToken{}

		err := proto.Unmarshal(iterator.Value(), &feeToken)
		if err != nil {
			panic(err)
		}

		feeTokens = append(feeTokens, feeToken)
	}

	return feeTokens
}
