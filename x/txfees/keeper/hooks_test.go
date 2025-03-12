package keeper_test

import (
	"time"

	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

func (suite *KeeperTestSuite) TestTxFeesAfterEpochEnd() {
	uion := "uion"
	atom := "atom"

	baseDenom := sdk.DefaultBondDenom

	tests := []struct {
		name        string
		coins       sdk.Coins
		burnedDenom string
	}{
		{
			name:        "DYM is burned",
			coins:       sdk.Coins{sdk.NewInt64Coin(baseDenom, 100000)},
			burnedDenom: baseDenom,
		},
		{
			name:        "One non-dym fee token (uion)",
			coins:       sdk.Coins{sdk.NewInt64Coin(uion, 1000)},
			burnedDenom: baseDenom,
		},
		{
			name:        "Multiple non-dym fee token",
			coins:       sdk.Coins{sdk.NewInt64Coin(baseDenom, 2000), sdk.NewInt64Coin(uion, 30000)},
			burnedDenom: baseDenom,
		},
		{
			name:        "unknown fee token is burned as well",
			coins:       sdk.Coins{sdk.NewInt64Coin(atom, 2000)},
			burnedDenom: atom,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()

		// create pools for three separate fee tokens
		suite.PrepareBalancerPoolWithCoins(sdk.NewCoin(baseDenom, math.NewInt(1000000000000)), sdk.NewCoin(uion, math.NewInt(5000)))

		moduleAddrFee := suite.App.AccountKeeper.GetModuleAddress(types.ModuleName)
		suite.FundModuleAcc(types.ModuleName, tc.coins)
		balances := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)
		suite.Assert().Equal(balances, tc.coins, tc.name)

		totalSupplyBefore := suite.App.BankKeeper.GetSupply(suite.Ctx, tc.burnedDenom).Amount

		// End of epoch, so all the non-dym fee amount should be swapped to dym and burned
		futureCtx := suite.Ctx.WithBlockTime(time.Now().Add(time.Minute))
		suite.App.TxFeesKeeper.AfterEpochEnd(futureCtx, "day", int64(1))

		// check the balance of the native-basedenom in module
		balances = suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)
		totalSupplyAfter := suite.App.BankKeeper.GetSupply(suite.Ctx, tc.burnedDenom).Amount

		//Check for token burn
		suite.Assert().True(balances.IsZero(), tc.name)
		suite.Require().True(totalSupplyAfter.LT(totalSupplyBefore), tc.name)
	}
}

// TestPoolCreationHooks validates fee token registration through sequential pool creation
func (suite *KeeperTestSuite) TestPoolCreationHooks() {
	baseDenom := suite.App.TxFeesKeeper.MustGetBaseDenom(suite.Ctx)

	gammParams := suite.App.GAMMKeeper.GetParams(suite.Ctx)
	gammParams.AllowedPoolCreationDenoms = []string{baseDenom, "tokenBase2"}
	suite.App.GAMMKeeper.SetParams(suite.Ctx, gammParams)

	// Initial pool: base <-> tokenA
	pool1 := suite.PrepareBalancerPoolWithCoins(
		sdk.NewCoin(baseDenom, math.NewInt(1e18)),
		sdk.NewCoin("tokenA", math.NewInt(1e18)),
	)

	// Verify tokenA registration
	feeTokenA, err := suite.App.TxFeesKeeper.GetFeeToken(suite.Ctx, "tokenA")
	suite.Require().NoError(err)
	suite.Require().Equal(pool1, feeTokenA.Route[0].PoolId)

	// Second pool: base <-> tokenBase2
	pool2 := suite.PrepareBalancerPoolWithCoins(
		sdk.NewCoin(baseDenom, math.NewInt(1e18)),
		sdk.NewCoin("tokenBase2", math.NewInt(1e18)),
	)

	// Verify tokenBase2 registration
	feeTokenBase2, err := suite.App.TxFeesKeeper.GetFeeToken(suite.Ctx, "tokenBase2")
	suite.Require().NoError(err)
	suite.Require().Equal(pool2, feeTokenBase2.Route[0].PoolId)

	// Third pool: tokenBase2 <-> tokenB
	pool3 := suite.PrepareBalancerPoolWithCoins(
		sdk.NewCoin("tokenBase2", math.NewInt(1e18)),
		sdk.NewCoin("tokenB", math.NewInt(1e18)),
	)

	// Verify tokenB has composite route
	feeTokenB, err := suite.App.TxFeesKeeper.GetFeeToken(suite.Ctx, "tokenB")
	suite.Require().NoError(err)
	suite.Require().Len(feeTokenB.Route, 2)
	suite.Require().Equal(pool3, feeTokenB.Route[0].PoolId)
	suite.Require().Equal(pool2, feeTokenB.Route[1].PoolId)

	// FIXME: create tokenB <-> basedenom and assert it updates to this route
}
