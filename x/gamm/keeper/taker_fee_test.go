package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

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
		tokenOutMinAmount sdk.Int
		expectSwap        bool
		expectError       bool
	}{
		"zero hops": {
			routes:            []poolmanagertypes.SwapAmountInRoute{},
			tokenIn:           sdk.NewCoin("foo", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       true,
		},
		"adym as tokenIn": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "foo",
				},
			},
			tokenIn:           sdk.NewCoin("adym", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       false,
		},
		"adym swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        1,
					TokenOutDenom: "adym",
				},
			},
			tokenIn:           sdk.NewCoin("foo", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       false,
		},
		"usdc swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        2,
					TokenOutDenom: "bar",
				},
			},
			tokenIn:           sdk.NewCoin("foo", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       false,
		},
		"usdc as token in": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        2,
					TokenOutDenom: "foo",
				},
			},
			tokenIn:           sdk.NewCoin("bar", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       false,
		},
		"usdc as token in - no route to dym": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        4,
					TokenOutDenom: "baz",
				},
			},
			tokenIn:           sdk.NewCoin("bar", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       false,
		},
		"usdc swap with dym": {
			routes: []poolmanagertypes.SwapAmountInRoute{
				{
					PoolId:        3,
					TokenOutDenom: "adym",
				},
			},
			tokenIn:           sdk.NewCoin("bar", sdk.NewInt(tokenInAmt)),
			tokenOutMinAmount: sdk.NewInt(1),
			expectError:       false,
		},
	}

	for name, tc := range testcases {
		suite.SetupTest()

		suite.App.TxFeesKeeper.SetBaseDenom(suite.Ctx, "adym")

		suite.FundAcc(suite.TestAccs[0], apptesting.DefaultAcctFunds)
		params := suite.App.GAMMKeeper.GetParams(suite.Ctx)
		params.PoolCreationFee = sdk.NewCoins(
			sdk.NewCoin("adym", sdk.NewInt(100000)),
			sdk.NewCoin("bar", sdk.NewInt(100000)))
		suite.App.GAMMKeeper.SetParams(suite.Ctx, params)

		ctx := suite.Ctx
		msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)

		pool1coins := []sdk.Coin{sdk.NewCoin("adym", sdk.NewInt(100000)), sdk.NewCoin("foo", sdk.NewInt(100000))}
		suite.PrepareBalancerPoolWithCoins(pool1coins...)

		//"bar" is treated as baseDenom (e.g. USDC)
		pool2coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000)), sdk.NewCoin("foo", sdk.NewInt(100000))}
		suite.PrepareBalancerPoolWithCoins(pool2coins...)

		pool3coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000)), sdk.NewCoin("adym", sdk.NewInt(100000))}
		suite.PrepareBalancerPoolWithCoins(pool3coins...)

		pool4coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000)), sdk.NewCoin("baz", sdk.NewInt(100000))}
		suite.PrepareBalancerPoolWithCoins(pool4coins...)

		//get the balance of txfees before swap
		moduleAddrFee := suite.App.AccountKeeper.GetModuleAddress(txfeestypes.ModuleName)
		balancesBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)

		// check taker fee is not 0
		suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(sdk.ZeroDec()))

		// make swap
		_, err := msgServer.SwapExactAmountIn(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountIn{
			Sender:            suite.TestAccs[0].String(),
			Routes:            tc.routes,
			TokenIn:           tc.tokenIn,
			TokenOutMinAmount: tc.tokenOutMinAmount,
		})
		if tc.expectError {
			suite.Require().Error(err, name)
			continue
		}
		suite.Require().NoError(err, name)

		//get the balance of txfees after swap
		balancesAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)

		testDenom := tc.tokenIn.Denom
		if tc.expectSwap {
			testDenom = tc.routes[0].TokenOutDenom
		}
		suite.Require().True(balancesAfter.AmountOf(testDenom).GT(balancesBefore.AmountOf(testDenom)), name)
	}
}

func (suite *KeeperTestSuite) TestTakerFeeCharged_ExactOut() {
	tokenInAmt := int64(100000)
	testcases := map[string]struct {
		routes      []poolmanagertypes.SwapAmountOutRoute
		tokenOut    sdk.Coin
		expectSwap  bool
		expectError bool
	}{
		"zero hops": {
			routes:      []poolmanagertypes.SwapAmountOutRoute{},
			tokenOut:    sdk.NewCoin("foo", sdk.NewInt(tokenInAmt)),
			expectError: true,
		},
		"adym as tokenIn": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: "adym",
				},
			},
			tokenOut:    sdk.NewCoin("foo", sdk.NewInt(tokenInAmt)),
			expectError: false,
		},
		"adym swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       1,
					TokenInDenom: "foo",
				},
			},
			tokenOut:    sdk.NewCoin("adym", sdk.NewInt(tokenInAmt)),
			expectError: false,
		},
		"usdc swapped in first pool": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       2,
					TokenInDenom: "foo",
				},
			},
			tokenOut:    sdk.NewCoin("bar", sdk.NewInt(tokenInAmt)),
			expectError: false,
		},
		"usdc as token in": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       2,
					TokenInDenom: "bar",
				},
			},
			tokenOut:    sdk.NewCoin("foo", sdk.NewInt(tokenInAmt)),
			expectError: false,
		},
		"usdc as token in - no route to dym": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       4,
					TokenInDenom: "bar",
				},
			},
			tokenOut:    sdk.NewCoin("baz", sdk.NewInt(tokenInAmt)),
			expectError: false,
		},
		"baz swap with usdc": {
			routes: []poolmanagertypes.SwapAmountOutRoute{
				{
					PoolId:       4,
					TokenInDenom: "baz",
				},
			},
			tokenOut:    sdk.NewCoin("bar", sdk.NewInt(tokenInAmt)),
			expectSwap:  true,
			expectError: false,
		},
	}

	for name, tc := range testcases {
		suite.SetupTest()

		suite.App.TxFeesKeeper.SetBaseDenom(suite.Ctx, "adym")

		suite.FundAcc(suite.TestAccs[0], apptesting.DefaultAcctFunds)
		params := suite.App.GAMMKeeper.GetParams(suite.Ctx)
		params.PoolCreationFee = sdk.NewCoins(
			sdk.NewCoin("adym", sdk.NewInt(1000)),
			sdk.NewCoin("bar", sdk.NewInt(1000)))
		suite.App.GAMMKeeper.SetParams(suite.Ctx, params)

		ctx := suite.Ctx
		msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)

		pool1coins := []sdk.Coin{sdk.NewCoin("adym", sdk.NewInt(100000000)), sdk.NewCoin("foo", sdk.NewInt(100000000))}
		suite.PrepareBalancerPoolWithCoins(pool1coins...)

		//"bar" is treated as baseDenom (e.g. USDC)
		pool2coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000000)), sdk.NewCoin("foo", sdk.NewInt(100000000))}
		suite.PrepareBalancerPoolWithCoins(pool2coins...)

		pool3coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000000)), sdk.NewCoin("adym", sdk.NewInt(100000000))}
		suite.PrepareBalancerPoolWithCoins(pool3coins...)

		pool4coins := []sdk.Coin{sdk.NewCoin("bar", sdk.NewInt(100000000)), sdk.NewCoin("baz", sdk.NewInt(100000000))}
		suite.PrepareBalancerPoolWithCoins(pool4coins...)

		//get the balance of txfees before swap
		moduleAddrFee := suite.App.AccountKeeper.GetModuleAddress(txfeestypes.ModuleName)
		balancesBefore := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)

		// check taker fee is not 0
		suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(sdk.ZeroDec()))

		// make swap
		_, err := msgServer.SwapExactAmountOut(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountOut{
			Sender:           suite.TestAccs[0].String(),
			Routes:           tc.routes,
			TokenOut:         tc.tokenOut,
			TokenInMaxAmount: sdk.NewInt(1000000000000000000),
		})
		if tc.expectError {
			suite.Require().Error(err, name)
			continue
		}
		suite.Require().NoError(err, name)

		//get the balance of txfees after swap
		balancesAfter := suite.App.BankKeeper.GetAllBalances(suite.Ctx, moduleAddrFee)

		testDenom := tc.routes[0].TokenInDenom
		if tc.expectSwap {
			testDenom = tc.tokenOut.Denom
		}
		suite.Require().True(balancesAfter.AmountOf(testDenom).GT(balancesBefore.AmountOf(testDenom)), testDenom, name)
	}
}

//TODO: test estimation when taker fee is 0

// TestEstimateMultihopSwapExactAmountIn tests that the estimation done via `EstimateSwapExactAmountIn`
func (suite *KeeperTestSuite) TestEstimateMultihopSwapExactAmountIn() {
	type param struct {
		routes            []poolmanagertypes.SwapAmountInRoute
		tokenIn           sdk.Coin
		tokenOutMinAmount sdk.Int
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
				tokenIn:           sdk.NewCoin("foo", sdk.NewInt(100000)),
				tokenOutMinAmount: sdk.NewInt(1),
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
				tokenIn:           sdk.NewCoin("adym", sdk.NewInt(100000)),
				tokenOutMinAmount: sdk.NewInt(1),
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
			reducedTokenIn := sdk.NewDecFromInt(test.param.tokenIn.Amount).MulTruncate(sdk.OneDec().Sub(suite.App.GAMMKeeper.GetParams(suite.Ctx).TakerFee))
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
		tokenInMinAmount sdk.Int
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
				tokenInMinAmount: sdk.NewInt(1),
				tokenOut:         sdk.NewCoin("baz", sdk.NewInt(100000)),
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
				tokenInMinAmount: sdk.NewInt(1),
				tokenOut:         sdk.NewCoin("baz", sdk.NewInt(100000)),
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
			tokensAfterTakerFeeReduction := sdk.NewDecFromInt(estimateMultihopTokenInAmountWithTakerFee.TokenInAmount).MulTruncate(sdk.OneDec().Sub(takerFee))

			// Now reducing taker fee from the input, we expect the estimation to be the same
			suite.Require().Equal(tokensAfterTakerFeeReduction.TruncateInt(), multihopTokenInAmount)
		})
	}
}
