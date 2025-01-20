package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"

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
			coins:       sdk.Coins{math.NewInt64Coin(baseDenom, 100000)},
			burnedDenom: baseDenom,
		},
		{
			name:        "One non-dym fee token (uion)",
			coins:       sdk.Coins{math.NewInt64Coin(uion, 1000)},
			burnedDenom: baseDenom,
		},
		{
			name:        "Multiple non-dym fee token",
			coins:       sdk.Coins{math.NewInt64Coin(baseDenom, 2000), math.NewInt64Coin(uion, 30000)},
			burnedDenom: baseDenom,
		},
		{
			name:        "unknown fee token is burned as well",
			coins:       sdk.Coins{math.NewInt64Coin(atom, 2000)},
			burnedDenom: atom,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()

		// create pools for three separate fee tokens
		suite.PrepareBalancerPoolWithCoins(sdk.NewCoin(baseDenom, math.NewInt(1000000000000)), sdk.NewCoin(uion, math.NewInt(5000)))

		moduleAddrFee := suite.App.AccountKeeper.GetModuleAddress(types.ModuleName)
		err := bankutil.FundModuleAccount(suite.App.BankKeeper, suite.Ctx, types.ModuleName, tc.coins)
		suite.Require().NoError(err)
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

//TODO: pool hooks
