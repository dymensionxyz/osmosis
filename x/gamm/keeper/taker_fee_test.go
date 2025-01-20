package keeper_test

import (
	"errors"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"

	"github.com/osmosis-labs/osmosis/v15/testutils/apptesting"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/keeper"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	txfeestypes "github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

/* ----------------------------- Testing ExactIn ---------------------------- */
func (suite *KeeperTestSuite) TestTakerFeeCharged_ExactIn() {
	tokenInAmt := int64(100000)
	testcases := map[string]struct {
		routes            []poolmanagertypes.SwapAmountInRoute
		tokenIn           sdk.Coin
		tokenOutMinAmount math.Int
		expectBeneficiary bool
		expectSwap        bool
		expectError       bool
	}{
		"zero hops": {
			routes:            []poolmanagertypes.SwapAmountInRoute{},
			tokenIn:           sdk.NewCoin("foo", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: false,
			expectError:       true,
		},
		"adym as tokenIn": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "foo",
				},
			},
			tokenIn:           sdk.NewCoin("adym", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: true,
			expectError:       false,
		},
		"adym swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "adym",
				},
			},
			tokenIn:           sdk.NewCoin("foo", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: true,
			expectError:       false,
		},
		"usdc swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        2,
					TokenOutDenom: "bar",
				},
			},
			tokenIn:           sdk.NewCoin("foo", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: false,
			expectError:       false,
		},
		"usdc as token in": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        2,
					TokenOutDenom: "foo",
				},
			},
			tokenIn:           sdk.NewCoin("bar", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: false,
			expectError:       false,
		},
		"usdc as token in - no route to dym": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        4,
					TokenOutDenom: "baz",
				},
			},
			tokenIn:           sdk.NewCoin("bar", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: false,
			expectError:       false,
		},
		"usdc swap with dym": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        3,
					TokenOutDenom: "adym",
				},
			},
			tokenIn:           sdk.NewCoin("bar", math.NewInt(tokenInAmt)),
			tokenOutMinAmount: math.NewInt(1),
			expectBeneficiary: true,
			expectError:       false,
		},
	}

	for name, tc := range testcases {
		suite.Run(name, func() {
			suite.SetupTest()

			// set mock rollapp keeper with a random beneficiary for testing taker fees
			beneficiary := apptesting.CreateRandomAccounts(1)[0]
			rollappKeeper := new(RollappKeeperMock)
			suite.App.GAMMKeeper.SetRollapp(rollappKeeper)
			// we consider adym as a rollapp token for convenience
			// if either IN or OUR token is adym, we expect the rollapp owner to be rewarded
			rollappKeeper.On("GetRollappOwnerByDenom", mock.Anything, "adym").Return(beneficiary, nil)
			rollappKeeper.On("GetRollappOwnerByDenom", mock.Anything, mock.Anything).Return(nil, errors.New("not a rollapp token"))

			suite.App.TxFeesKeeper.SetBaseDenom(suite.Ctx, "adym")

			suite.FundAcc(suite.TestAccs[0], apptesting.DefaultAcctFunds)
			params := suite.App.GAMMKeeper.GetParams(suite.Ctx)
			params.PoolCreationFee = sdk.NewCoins(
				sdk.NewCoin("adym", math.NewInt(100000)),
				sdk.NewCoin("bar", math.NewInt(100000)))
			suite.App.GAMMKeeper.SetParams(suite.Ctx, params)

			ctx := suite.Ctx
			msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)

			pool1coins := []sdk.Coin{sdk.NewCoin("adym", math.NewInt(100000)), sdk.NewCoin("foo", math.NewInt(100000))}
			suite.PrepareBalancerPoolWithCoins(pool1coins...)

			//"bar" is treated as baseDenom (e.g. USDC)
			pool2coins := []sdk.Coin{sdk.NewCoin("bar", math.NewInt(100000)), sdk.NewCoin("foo", math.NewInt(100000))}
			suite.PrepareBalancerPoolWithCoins(pool2coins...)

			pool3coins := []sdk.Coin{sdk.NewCoin("bar", math.NewInt(100000)), sdk.NewCoin("adym", math.NewInt(100000))}
			suite.PrepareBalancerPoolWithCoins(pool3coins...)

			pool4coins := []sdk.Coin{sdk.NewCoin("bar", math.NewInt(100000)), sdk.NewCoin("baz", math.NewInt(100000))}
			suite.PrepareBalancerPoolWithCoins(pool4coins...)

			//get the balance of txfees before swap
			moduleAddrFee := suite.App.AccountKeeper.GetModuleAddress(txfeestypes.ModuleName)
			balancesBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)
			beneficiaryBalancesBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, beneficiary)

			// check taker fee is not 0
			suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(math.LegacyZeroDec()))

			// make swap
			_, err := msgServer.SwapExactAmountIn(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountIn{
				Sender:            suite.TestAccs[0].String(),
				Routes:            tc.routes,
				TokenIn:           tc.tokenIn,
				TokenOutMinAmount: tc.tokenOutMinAmount,
			})
			if tc.expectError {
				suite.Require().Error(err, name)
				return
			}
			suite.Require().NoError(err, name)

			//get the balance of txfees after swap
			balancesAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)

			testDenom := tc.tokenIn.Denom
			if tc.expectSwap {
				testDenom = tc.routes[0].TokenOutDenom
			}
			// x/txfees balance is the same as initially since the fees are distributed immediately
			suite.Require().True(balancesAfter.AmountOf(testDenom).Equal(balancesBefore.AmountOf(testDenom)), name)

			// check if the beneficiary is rewarded
			beneficiaryBalancesAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, beneficiary)
			baseDenom, err := suite.App.TxFeesKeeper.GetBaseDenom(suite.Ctx)
			suite.Require().NoError(err)

			beneficiaryFeeBefore := beneficiaryBalancesBefore.AmountOf(baseDenom)
			beneficiaryFeeAfter := beneficiaryBalancesAfter.AmountOf(baseDenom)
			if tc.expectBeneficiary {
				suite.Require().True(beneficiaryFeeBefore.LT(beneficiaryFeeAfter))
			} else {
				suite.Require().True(beneficiaryFeeBefore.Equal(beneficiaryFeeAfter))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestTakerFeeCharged_ExactOut() {
	tokenInAmt := int64(100000)
	testcases := map[string]struct {
		routes            []poolmanagertypes.SwapAmountOutRoute
		tokenOut          sdk.Coin
		expectBeneficiary bool
		expectSwap        bool
		expectError       bool
	}{
		"zero hops": {
			routes:            []poolmanagertypes.SwapAmountOutRoute{},
			tokenOut:          sdk.NewCoin("foo", math.NewInt(tokenInAmt)),
			expectBeneficiary: false,
			expectError:       true,
		},
		"adym as tokenIn": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: "adym",
				},
			},
			tokenOut:          sdk.NewCoin("foo", math.NewInt(tokenInAmt)),
			expectBeneficiary: true,
			expectError:       false,
		},
		"adym swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: "foo",
				},
			},
			tokenOut:          sdk.NewCoin("adym", math.NewInt(tokenInAmt)),
			expectBeneficiary: true,
			expectError:       false,
		},
		"usdc swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       2,
					TokenInDenom: "foo",
				},
			},
			tokenOut:          sdk.NewCoin("bar", math.NewInt(tokenInAmt)),
			expectBeneficiary: false,
			expectError:       false,
		},
		"usdc as token in": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       2,
					TokenInDenom: "bar",
				},
			},
			tokenOut:          sdk.NewCoin("foo", math.NewInt(tokenInAmt)),
			expectBeneficiary: false,
			expectError:       false,
		},
		"usdc as token in - no route to dym": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       4,
					TokenInDenom: "bar",
				},
			},
			tokenOut:          sdk.NewCoin("baz", math.NewInt(tokenInAmt)),
			expectBeneficiary: false,
			expectError:       false,
		},
		"baz swap with usdc": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       4,
					TokenInDenom: "baz",
				},
			},
			tokenOut:          sdk.NewCoin("bar", math.NewInt(tokenInAmt)),
			expectBeneficiary: false,
			expectSwap:        true,
			expectError:       false,
		},
	}

	for name, tc := range testcases {
		suite.Run(name, func() {
			suite.SetupTest()

			// set mock rollapp keeper with a random beneficiary for testing taker fees
			beneficiary := apptesting.CreateRandomAccounts(1)[0]
			rollappKeeper := new(RollappKeeperMock)
			suite.App.GAMMKeeper.SetRollapp(rollappKeeper)
			// we consider adym as a rollapp token for convenience
			// if either IN or OUR token is adym, we expect the rollapp owner to be rewarded
			rollappKeeper.On("GetRollappOwnerByDenom", mock.Anything, "adym").Return(beneficiary, nil)
			rollappKeeper.On("GetRollappOwnerByDenom", mock.Anything, mock.Anything).Return(nil, errors.New("not a rollapp token"))

			suite.App.TxFeesKeeper.SetBaseDenom(suite.Ctx, "adym")

			suite.FundAcc(suite.TestAccs[0], apptesting.DefaultAcctFunds)
			params := suite.App.GAMMKeeper.GetParams(suite.Ctx)
			params.PoolCreationFee = sdk.NewCoins(
				sdk.NewCoin("adym", math.NewInt(1000)),
				sdk.NewCoin("bar", math.NewInt(1000)))
			suite.App.GAMMKeeper.SetParams(suite.Ctx, params)

			ctx := suite.Ctx
			msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)

			pool1coins := []sdk.Coin{sdk.NewCoin("adym", math.NewInt(100000000)), sdk.NewCoin("foo", math.NewInt(100000000))}
			suite.PrepareBalancerPoolWithCoins(pool1coins...)

			//"bar" is treated as baseDenom (e.g. USDC)
			pool2coins := []sdk.Coin{sdk.NewCoin("bar", math.NewInt(100000000)), sdk.NewCoin("foo", math.NewInt(100000000))}
			suite.PrepareBalancerPoolWithCoins(pool2coins...)

			pool3coins := []sdk.Coin{sdk.NewCoin("bar", math.NewInt(100000000)), sdk.NewCoin("adym", math.NewInt(100000000))}
			suite.PrepareBalancerPoolWithCoins(pool3coins...)

			pool4coins := []sdk.Coin{sdk.NewCoin("bar", math.NewInt(100000000)), sdk.NewCoin("baz", math.NewInt(100000000))}
			suite.PrepareBalancerPoolWithCoins(pool4coins...)

			//get the balance of txfees before swap
			moduleAddrFee := suite.App.AccountKeeper.GetModuleAddress(txfeestypes.ModuleName)
			balancesBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)
			beneficiaryBalancesBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, beneficiary)

			// check taker fee is not 0
			suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(math.LegacyZeroDec()))

			// make swap
			_, err := msgServer.SwapExactAmountOut(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountOut{
				Sender:           suite.TestAccs[0].String(),
				Routes:           tc.routes,
				TokenOut:         tc.tokenOut,
				TokenInMaxAmount: math.NewInt(1000000000000000000),
			})
			if tc.expectError {
				suite.Require().Error(err, name)
				return
			}
			suite.Require().NoError(err, name)

			//get the balance of txfees after swap
			balancesAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)

			testDenom := tc.routes[0].TokenInDenom
			if tc.expectSwap {
				testDenom = tc.tokenOut.Denom
			}
			// x/txfees balance is the same as initially since the fees are distributed immediately
			suite.Require().True(balancesAfter.AmountOf(testDenom).Equal(balancesBefore.AmountOf(testDenom)), testDenom, name)

			// check if the beneficiary is rewarded
			beneficiaryBalancesAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, beneficiary)
			baseDenom, err := suite.App.TxFeesKeeper.GetBaseDenom(suite.Ctx)
			suite.Require().NoError(err)

			beneficiaryFeeBefore := beneficiaryBalancesBefore.AmountOf(baseDenom)
			beneficiaryFeeAfter := beneficiaryBalancesAfter.AmountOf(baseDenom)
			if tc.expectBeneficiary {
				suite.Require().True(beneficiaryFeeBefore.LT(beneficiaryFeeAfter))
			} else {
				suite.Require().True(beneficiaryFeeBefore.Equal(beneficiaryFeeAfter))
			}
		})
	}
}

//TODO: test estimation when taker fee is 0

// TestEstimateMultihopSwapExactAmountIn tests that the estimation done via `EstimateSwapExactAmountIn`
func (suite *KeeperTestSuite) TestEstimateMultihopSwapExactAmountIn() {
	type param struct {
		routes            []poolmanagertypes.SwapAmountInRoute
		tokenIn           sdk.Coin
		tokenOutMinAmount math.Int
	}

	tests := []struct {
		name     string
		param    param
		poolType poolmanagertypes.PoolType
	}{
		{
			name: "Proper swap - foo -> bar(pool 1) - bar(pool 2) -> baz",
			param: param{
				routes: []poolmanagertypes.SwapAmountInRoute{
					{
						PoolId:        1,
						TokenOutDenom: "bar",
					},
					{
						PoolId:        2,
						TokenOutDenom: "baz",
					},
				},
				tokenIn:           sdk.NewCoin("foo", math.NewInt(100000)),
				tokenOutMinAmount: math.NewInt(1),
			},
		},
		{
			name: "Swap - foo -> adym(pool 1) - adym(pool 2) -> baz ",
			param: param{
				routes: []poolmanagertypes.SwapAmountInRoute{
					{
						PoolId:        1,
						TokenOutDenom: "adym",
					},
					{
						PoolId:        2,
						TokenOutDenom: "baz",
					},
				},
				tokenIn:           sdk.NewCoin("adym", math.NewInt(100000)),
				tokenOutMinAmount: math.NewInt(1),
			},
		},
	}

	for _, test := range tests {
		// Init suite for each test.
		suite.SetupTest()

		suite.Run(test.name, func() {
			firstEstimatePoolId := suite.PrepareBalancerPool()
			secondEstimatePoolId := suite.PrepareBalancerPool()

			firstEstimatePool, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, firstEstimatePoolId)
			suite.Require().NoError(err)
			secondEstimatePool, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, secondEstimatePoolId)
			suite.Require().NoError(err)

			// calculate token out amount using `EstimateMultihopSwapExactAmountIn`
			queryClient := suite.queryClient
			estimateMultihopTokenOutAmountWithTakerFee, errEstimate := queryClient.EstimateSwapExactAmountIn(
				suite.Ctx,
				&types.QuerySwapExactAmountInRequest{
					TokenIn: test.param.tokenIn.String(),
					Routes:  test.param.routes,
				},
			)
			suite.Require().NoError(errEstimate, "test: %v", test.name)

			// ensure that pool state has not been altered after estimation
			firstEstimatePoolAfterSwap, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, firstEstimatePoolId)
			suite.Require().NoError(err)
			secondEstimatePoolAfterSwap, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, secondEstimatePoolId)
			suite.Require().NoError(err)
			suite.Require().Equal(firstEstimatePool, firstEstimatePoolAfterSwap)
			suite.Require().Equal(secondEstimatePool, secondEstimatePoolAfterSwap)

			// calculate token out amount using `MultihopSwapExactAmountIn`
			poolmanagerKeeper := suite.App.PoolManagerKeeper
			multihopTokenOutAmount, errMultihop := poolmanagerKeeper.MultihopEstimateOutGivenExactAmountIn(
				suite.Ctx,
				test.param.routes,
				test.param.tokenIn,
			)
			suite.Require().NoError(errMultihop, "test: %v", test.name)
			// the pool manager estimation is without taker fee, so it should be higher
			suite.Require().True(multihopTokenOutAmount.GT(estimateMultihopTokenOutAmountWithTakerFee.TokenOutAmount))

			// Now reducing taker fee from the input, we expect the estimation to be the same
			reducedTokenIn := math.LegacyNewDecFromInt(test.param.tokenIn.Amount).MulTruncate(math.LegacyOneDec().Sub(suite.App.GAMMKeeper.GetParams(suite.Ctx).TakerFee))
			reducedTokenInCoin := sdk.NewCoin(test.param.tokenIn.Denom, reducedTokenIn.TruncateInt())

			multihopTokenOutAmountTakerFeeReduced, errMultihop := poolmanagerKeeper.MultihopEstimateOutGivenExactAmountIn(
				suite.Ctx,
				test.param.routes,
				reducedTokenInCoin,
			)
			suite.Require().Equal(estimateMultihopTokenOutAmountWithTakerFee.TokenOutAmount, multihopTokenOutAmountTakerFeeReduced)
		})
	}
}

func (suite *KeeperTestSuite) TestEstimateMultihopSwapExactAmountOut() {
	type param struct {
		routes           []poolmanagertypes.SwapAmountOutRoute
		tokenOut         sdk.Coin
		tokenInMinAmount math.Int
	}

	tests := []struct {
		name     string
		param    param
		poolType poolmanagertypes.PoolType
	}{
		{
			name: "Proper swap - foo -> bar(pool 1) - bar(pool 2) -> baz",
			param: param{
				routes: []poolmanagertypes.SwapAmountOutRoute{
					{
						PoolId:       1,
						TokenInDenom: "foo",
					},
					{
						PoolId:       2,
						TokenInDenom: "bar",
					},
				},
				tokenInMinAmount: math.NewInt(1),
				tokenOut:         sdk.NewCoin("baz", math.NewInt(100000)),
			},
		},
		{
			name: "Swap - foo -> adym(pool 1) - adym(pool 2) -> baz ",
			param: param{
				routes: []poolmanagertypes.SwapAmountOutRoute{
					{
						PoolId:       1,
						TokenInDenom: "foo",
					},
					{
						PoolId:       2,
						TokenInDenom: "adym",
					},
				},
				tokenInMinAmount: math.NewInt(1),
				tokenOut:         sdk.NewCoin("baz", math.NewInt(100000)),
			},
		},
	}

	for _, test := range tests {
		// Init suite for each test.
		suite.SetupTest()

		suite.Run(test.name, func() {
			firstEstimatePoolId := suite.PrepareBalancerPool()
			secondEstimatePoolId := suite.PrepareBalancerPool()

			firstEstimatePool, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, firstEstimatePoolId)
			suite.Require().NoError(err)
			secondEstimatePool, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, secondEstimatePoolId)
			suite.Require().NoError(err)

			// calculate token out amount using `EstimateMultihopSwapExactAmountIn`
			queryClient := suite.queryClient
			estimateMultihopTokenInAmountWithTakerFee, errEstimate := queryClient.EstimateSwapExactAmountOut(
				suite.Ctx,
				&types.QuerySwapExactAmountOutRequest{
					TokenOut: test.param.tokenOut.String(),
					Routes:   test.param.routes,
				},
			)
			suite.Require().NoError(errEstimate, "test: %v", test.name)

			// ensure that pool state has not been altered after estimation
			firstEstimatePoolAfterSwap, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, firstEstimatePoolId)
			suite.Require().NoError(err)
			secondEstimatePoolAfterSwap, err := suite.App.GAMMKeeper.GetPoolAndPoke(suite.Ctx, secondEstimatePoolId)
			suite.Require().NoError(err)
			suite.Require().Equal(firstEstimatePool, firstEstimatePoolAfterSwap)
			suite.Require().Equal(secondEstimatePool, secondEstimatePoolAfterSwap)

			// calculate token out amount using `MultihopSwapExactAmountIn`
			poolmanagerKeeper := suite.App.PoolManagerKeeper
			multihopTokenInAmount, errMultihop := poolmanagerKeeper.MultihopEstimateInGivenExactAmountOut(
				suite.Ctx,
				test.param.routes,
				test.param.tokenOut,
			)
			suite.Require().NoError(errMultihop, "test: %v", test.name)
			// the pool manager estimation is without taker fee, so it should be lower (less tokens in for same amount out)
			suite.Require().True(multihopTokenInAmount.LT(estimateMultihopTokenInAmountWithTakerFee.TokenInAmount))

			takerFee := suite.App.GAMMKeeper.GetParams(suite.Ctx).TakerFee
			tokensAfterTakerFeeReduction := math.LegacyNewDecFromInt(estimateMultihopTokenInAmountWithTakerFee.TokenInAmount).MulTruncate(math.LegacyOneDec().Sub(takerFee))

			// Now reducing taker fee from the input, we expect the estimation to be the same
			suite.Require().Equal(tokensAfterTakerFeeReduction.TruncateInt(), multihopTokenInAmount)
		})
	}
}

type RollappKeeperMock struct {
	mock.Mock
}

func (m *RollappKeeperMock) GetRollappOwnerByDenom(ctx sdk.Context, denom string) (sdk.AccAddress, error) {
	args := m.Called(ctx, denom)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(sdk.AccAddress), args.Error(1)
}
