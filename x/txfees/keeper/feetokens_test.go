package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestBaseDenom() {
	suite.SetupTest()

	// Test getting basedenom (should be default from genesis)
	baseDenom, err := suite.App.TxFeesKeeper.GetBaseDenom(suite.Ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(sdk.DefaultBondDenom, baseDenom)

	converted, err := suite.App.TxFeesKeeper.ConvertToBaseToken(suite.Ctx, math.NewInt64Coin(sdk.DefaultBondDenom, 10))
	suite.Require().True(converted.IsEqual(math.NewInt64Coin(sdk.DefaultBondDenom, 10)))
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestFeeTokenConversions() {
	baseDenom := sdk.DefaultBondDenom

	tests := []struct {
		name                string
		baseDenomPoolInput  sdk.Coin
		feeTokenPoolInput   sdk.Coin
		inputFee            sdk.Coin
		expectedConvertable bool
		expectedOutput      sdk.Coin
	}{
		{
			name:                "equal value",
			baseDenomPoolInput:  math.NewInt64Coin(baseDenom, 100),
			feeTokenPoolInput:   math.NewInt64Coin("uion", 100),
			inputFee:            math.NewInt64Coin("uion", 10),
			expectedOutput:      math.NewInt64Coin(baseDenom, 10),
			expectedConvertable: true,
		},
		{
			name:               "unequal value",
			baseDenomPoolInput: math.NewInt64Coin(baseDenom, 100),
			feeTokenPoolInput:  math.NewInt64Coin("foo", 200),
			inputFee:           math.NewInt64Coin("foo", 10),
			// expected to get approximately 5 base denom
			// foo supply / stake supply =  200 / 100 = 2 foo for 1 stake
			// 10 foo in / 2 foo for 1 stake = 5 base denom
			expectedOutput:      math.NewInt64Coin(baseDenom, 5),
			expectedConvertable: true,
		},
		{
			name:                "basedenom value",
			baseDenomPoolInput:  math.NewInt64Coin(baseDenom, 100),
			feeTokenPoolInput:   math.NewInt64Coin("foo", 200),
			inputFee:            math.NewInt64Coin(baseDenom, 10),
			expectedOutput:      math.NewInt64Coin(baseDenom, 10),
			expectedConvertable: true,
		},
		{
			name:                "convert non-existent",
			baseDenomPoolInput:  math.NewInt64Coin(baseDenom, 100),
			feeTokenPoolInput:   math.NewInt64Coin("uion", 200),
			inputFee:            math.NewInt64Coin("foo", 10),
			expectedOutput:      sdk.Coin{},
			expectedConvertable: false,
		},
	}

	for _, tc := range tests {
		suite.SetupTest()

		_ = suite.PrepareBalancerPoolWithCoins(
			tc.baseDenomPoolInput,
			tc.feeTokenPoolInput,
		)

		converted, err := suite.App.TxFeesKeeper.ConvertToBaseToken(suite.Ctx, tc.inputFee)
		if tc.expectedConvertable {
			suite.Require().NoError(err, "test: %s", tc.name)
			suite.Require().Equal(tc.expectedOutput, converted)
		} else {
			suite.Require().Error(err, "test: %s", tc.name)
		}
	}
}
