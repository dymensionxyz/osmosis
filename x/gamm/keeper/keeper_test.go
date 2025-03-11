package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/testutils/apptesting"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/pool-models/balancer"
	balancertypes "github.com/osmosis-labs/osmosis/v15/x/gamm/pool-models/balancer"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	queryClient types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()
	suite.queryClient = types.NewQueryClient(suite.QueryHelper)

	params := suite.App.GAMMKeeper.GetParams(suite.Ctx)
	params.PoolCreationFee[0].Denom = "adym"
	params.AllowedPoolCreationDenoms = []string{"adym"}
	suite.App.GAMMKeeper.SetParams(suite.Ctx, params)

	// fund account for pool creation fee
	suite.FundAcc(suite.TestAccs[0], suite.App.GAMMKeeper.GetParams(suite.Ctx).PoolCreationFee)
}

func (suite *KeeperTestSuite) prepareCustomBalancerPool(
	balances sdk.Coins,
	poolAssets []balancertypes.PoolAsset,
	poolParams balancer.PoolParams,
) uint64 {
	suite.fundAllAccountsWith(balances)

	poolID, err := suite.App.PoolManagerKeeper.CreatePool(
		suite.Ctx,
		balancer.NewMsgCreateBalancerPool(suite.TestAccs[0], poolParams, poolAssets, ""),
	)
	suite.Require().NoError(err)

	return poolID
}

func (suite *KeeperTestSuite) fundAllAccountsWith(balances sdk.Coins) {
	for _, acc := range suite.TestAccs {
		suite.FundAcc(acc, balances)
	}
}
