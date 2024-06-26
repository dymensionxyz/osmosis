package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v15/testutils"
	"github.com/osmosis-labs/osmosis/v15/testutils/apptesting"
	"github.com/osmosis-labs/osmosis/v15/x/lockup/keeper"
)

type KeeperTestSuite struct {
	apptesting.KeeperTestHelper

	querier keeper.Querier
	cleanup func()
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.Setup()
	suite.querier = keeper.NewQuerier(*suite.App.LockupKeeper)
	unbondingDuration := suite.App.StakingKeeper.GetParams(suite.Ctx).UnbondingTime
	suite.App.IncentivesKeeper.SetLockableDurations(suite.Ctx, []time.Duration{
		time.Hour * 24 * 14,
		time.Hour,
		time.Hour * 3,
		time.Hour * 7,
		unbondingDuration,
	})
}

func (suite *KeeperTestSuite) SetupTestWithLevelDb() {
	suite.App, suite.cleanup = testutils.SetupTestingAppWithLevelDb(false)
	suite.Ctx = suite.App.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "osmosis_100-1", Time: time.Now().UTC()})
}

func (suite *KeeperTestSuite) Cleanup() {
	suite.cleanup()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
