package keeper_test

import (
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/gogoproto/proto"

	pooltypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/keeper"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
	v2 "github.com/osmosis-labs/osmosis/v15/x/txfees/types/migrations/v2"
)

func (suite *KeeperTestSuite) TestMigration_v1tov2() {
	suite.SetupTest()

	k := suite.App.TxFeesKeeper
	store := prefix.NewStore(suite.Ctx.KVStore(suite.App.GetKey(types.ModuleName)), types.FeeTokensStorePrefix)
	oldFeeToken := v2.FeeToken{
		Denom:  "uion",
		PoolID: 10,
	}

	bz, err := proto.Marshal(&oldFeeToken)
	suite.Require().NoError(err)
	store.Set([]byte("uion"), bz)

	// not marshaled correctly
	token, err := k.GetFeeToken(suite.Ctx, "uion")
	suite.Require().NoError(err)
	suite.Require().Len(token.Route, 0)

	m := keeper.NewMigrator(*k)
	err = m.Migrate1to2(suite.Ctx)
	suite.Require().NoError(err)

	feeToken, err := k.GetFeeToken(suite.Ctx, "uion")
	suite.Require().NoError(err)
	suite.Require().Equal(feeToken, types.FeeToken{
		Denom: "uion",
		Route: []pooltypes.SwapAmountInRoute{
			{
				PoolId:        10,
				TokenOutDenom: k.MustGetBaseDenom(suite.Ctx),
			},
		},
	})

}
