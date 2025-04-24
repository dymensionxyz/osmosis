package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestBaseDenom() {
	suite.SetupTest()

	// Test getting basedenom (should be default from genesis)
	baseDenom, err := suite.App.TxFeesKeeper.GetBaseDenom(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.DefaultBondDenom, baseDenom)

	converted, err := suite.App.TxFeesKeeper.CalcCoinInBaseDenom(suite.Ctx, sdk.NewInt64Coin(sdk.DefaultBondDenom, 10))
	suite.Require().True(converted.IsEqual(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)))
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestCalcCoinInBaseDenom() {
	baseDenom := sdk.DefaultBondDenom
	equalPoolAssets := sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 10_000_000), sdk.NewInt64Coin("uion", 10_000_000))
	diffPoolAssets := sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 10_000_000), sdk.NewInt64Coin("uion", 20_000_000))

	tests := []struct {
		name                string
		poolAssets          sdk.Coins
		inputFee            sdk.Coin
		expectedConvertable bool
		expectedOutput      sdk.Coin
	}{
		{
			name:                "equal value",
			poolAssets:          equalPoolAssets,
			inputFee:            sdk.NewInt64Coin("uion", 1000),
			expectedOutput:      sdk.NewInt64Coin(baseDenom, 999), // truncated
			expectedConvertable: true,
		},
		{
			name:       "unequal value",
			poolAssets: diffPoolAssets,
			inputFee:   sdk.NewInt64Coin("uion", 1000),
			// expected to get approximately 5 base denom (truncated to 4)
			expectedOutput:      sdk.NewInt64Coin(baseDenom, 499),
			expectedConvertable: true,
		},
		{
			name:                "basedenom value",
			poolAssets:          equalPoolAssets,
			inputFee:            sdk.NewInt64Coin(baseDenom, 1000),
			expectedOutput:      sdk.NewInt64Coin(baseDenom, 1000),
			expectedConvertable: true,
		},
		{
			name:                "convert non-existent",
			poolAssets:          equalPoolAssets,
			inputFee:            sdk.NewInt64Coin("foo", 1000),
			expectedConvertable: false,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()

		_ = suite.PrepareBalancerPoolWithCoins(
			tc.poolAssets...,
		)

		converted, err := suite.App.TxFeesKeeper.CalcCoinInBaseDenom(suite.Ctx, tc.inputFee)
		if tc.expectedConvertable {
			suite.Require().NoError(err, "test: %s", tc.name)
			suite.Require().Equal(tc.expectedOutput, converted)
		} else {
			suite.Require().Error(err, "test: %s", tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestBaseInCoinConversions() {
	baseDenom := sdk.DefaultBondDenom
	equalPoolAssets := sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 10_000_000), sdk.NewInt64Coin("uion", 10_000_000))
	diffPoolAssets := sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 10_000_000), sdk.NewInt64Coin("uion", 20_000_000))

	tests := []struct {
		name                string
		poolAssets          sdk.Coins
		inputBaseCoin       sdk.Coin
		targetDenom         string
		expectedConvertable bool
		expectedOutput      sdk.Coin
	}{
		{
			name:                "equal value",
			poolAssets:          equalPoolAssets,
			inputBaseCoin:       sdk.NewInt64Coin(baseDenom, 1000),
			targetDenom:         "uion",
			expectedOutput:      sdk.NewInt64Coin("uion", 999),
			expectedConvertable: true,
		},
		{
			name:                "unequal value",
			poolAssets:          diffPoolAssets,
			inputBaseCoin:       sdk.NewInt64Coin(baseDenom, 500),
			targetDenom:         "uion",
			expectedOutput:      sdk.NewInt64Coin("uion", 999),
			expectedConvertable: true,
		},
		{
			name:                "non-base denom input",
			poolAssets:          equalPoolAssets,
			inputBaseCoin:       sdk.NewInt64Coin("nonbase", 1000),
			targetDenom:         "uion",
			expectedOutput:      sdk.Coin{},
			expectedConvertable: false,
		},
		{
			name:                "convert to non-existent denom",
			poolAssets:          equalPoolAssets,
			inputBaseCoin:       sdk.NewInt64Coin(baseDenom, 1000),
			targetDenom:         "nonexistent",
			expectedOutput:      sdk.Coin{},
			expectedConvertable: false,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()

		_ = suite.PrepareBalancerPoolWithCoins(
			tc.poolAssets...,
		)

		// Test CalcBaseInCoin
		converted, err := suite.App.TxFeesKeeper.CalcBaseInCoin(suite.Ctx, tc.inputBaseCoin, tc.targetDenom)
		if !tc.expectedConvertable {
			suite.Require().Error(err, "test: %s", tc.name)
		} else {
			suite.Require().NoError(err, "test: %s", tc.name)
			suite.Require().Equal(tc.expectedOutput, converted, "test: %s", tc.name)

			// Verify bidirectional conversion works
			// Convert back to base denom and check if we get close to the original amount
			// (accounting for truncation in both directions)
			reconverted, err := suite.App.TxFeesKeeper.CalcCoinInBaseDenom(suite.Ctx, converted)
			suite.Require().NoError(err, "test: %s - reconversion", tc.name)

			// The reconverted amount should be less than or equal to the original input
			// due to truncation in both conversions
			suite.Require().LessOrEqual(
				reconverted.Amount.Int64(),
				tc.inputBaseCoin.Amount.Int64(),
				"test: %s - reconversion amount should be <= original due to truncation",
				tc.name,
			)

			// Calculate the difference as a percentage of the original amount
			diff := tc.inputBaseCoin.Amount.Sub(reconverted.Amount)
			diffPercentage := math.LegacyNewDecFromInt(diff).Quo(math.LegacyNewDecFromInt(tc.inputBaseCoin.Amount)).MulInt64(100)

			// The difference should be at most 0.5% of the original amount
			suite.Require().LessOrEqual(
				diffPercentage.MustFloat64(),
				0.5,
				"test: %s - difference too large, got %.2f%%",
				tc.name,
				diffPercentage.MustFloat64(),
			)
		}
	}
}

func (suite *KeeperTestSuite) TestCalcWithMultiRoute() {
	baseDenom := sdk.DefaultBondDenom
	denom := "foo"

	pool1 := sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 10_000_000), sdk.NewInt64Coin("uion", 10_000_000))
	pool2 := sdk.NewCoins(sdk.NewInt64Coin("uion", 10_000_000), sdk.NewInt64Coin(denom, 10_000_000))

	tests := []struct {
		name                string
		inputCoin           sdk.Coin
		expectedConvertable bool
		expectedOutput      sdk.Coin
	}{
		{
			name:                "from coin",
			inputCoin:           sdk.NewInt64Coin(denom, 1000),
			expectedOutput:      sdk.NewInt64Coin(baseDenom, 999),
			expectedConvertable: true,
		},
		{
			name:                "from basedenom",
			inputCoin:           sdk.NewInt64Coin(baseDenom, 1000),
			expectedOutput:      sdk.NewInt64Coin(denom, 999),
			expectedConvertable: true,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()

		_ = suite.PrepareBalancerPoolWithCoins(
			pool1...,
		)

		_ = suite.PrepareBalancerPoolWithCoins(
			pool2...,
		)

		var converted sdk.Coin
		var err error

		if tc.inputCoin.Denom == baseDenom {
			converted, err = suite.App.TxFeesKeeper.CalcBaseInCoin(suite.Ctx, tc.inputCoin, denom)
		} else {
			converted, err = suite.App.TxFeesKeeper.CalcCoinInBaseDenom(suite.Ctx, tc.inputCoin)
		}
		if tc.expectedConvertable {
			suite.Require().NoError(err, "test: %s", tc.name)

			// Calculate the difference as a percentage of the original amount
			diff := tc.expectedOutput.Amount.Sub(converted.Amount)
			diffPercentage := math.LegacyNewDecFromInt(diff).Quo(math.LegacyNewDecFromInt(tc.expectedOutput.Amount)).MulInt64(100)

			// The difference should be at most 0.5% of the original amount
			suite.Require().LessOrEqual(
				diffPercentage.MustFloat64(),
				0.5,
				"test: %s - difference too large, got %.2f%%",
				tc.name,
				diffPercentage.MustFloat64(),
			)
		} else {
			suite.Require().Error(err, "test: %s", tc.name)
		}
	}
}
