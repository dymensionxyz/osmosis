package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/x/gamm/keeper"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
)

func (suite *KeeperTestSuite) TestDYMIsBurned_ExactIn() {
	suite.SetupTest()
	ctx := suite.Ctx
	msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)
	suite.PrepareBalancerPool()

	//check total supply before swap
	totalSupplyBefore := suite.App.BankKeeper.GetSupply(suite.Ctx, "udym")

	// check taker fee is not 0
	suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(sdk.ZeroDec()))

	routes :=
		[]poolmanagertypes.SwapAmountInRoute{
			{
				PoolId:        1,
				TokenOutDenom: "bar",
			},
		}

	// make swap
	_, err := msgServer.SwapExactAmountIn(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountIn{
		Sender:            suite.TestAccs[0].String(),
		Routes:            routes,
		TokenIn:           sdk.NewCoin("udym", sdk.NewInt(100000)),
		TokenOutMinAmount: sdk.NewInt(100),
	})
	suite.Require().NoError(err)

	// check total supply after swap
	totalSupplyAfter := suite.App.BankKeeper.GetSupply(suite.Ctx, "udym")

	//validate total supply is reduced by taker fee
	takerFeeAmount := suite.App.GAMMKeeper.GetParams(ctx).TakerFee.MulInt(sdk.NewInt(100000)).TruncateInt()
	suite.Require().True(totalSupplyAfter.Amount.LT(totalSupplyBefore.Amount))
	suite.Require().True(totalSupplyBefore.Amount.Sub(totalSupplyAfter.Amount).Equal(takerFeeAmount))
}

func (suite *KeeperTestSuite) TestDYMIsBurned_ExactOut() {
	suite.SetupTest()
	ctx := suite.Ctx
	msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)
	suite.PrepareBalancerPool()

	//check total supply before swap
	totalSupplyBefore := suite.App.BankKeeper.GetSupply(suite.Ctx, "udym")

	// check taker fee is not 0
	suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(sdk.ZeroDec()))

	routes :=
		[]poolmanagertypes.SwapAmountOutRoute{
			{
				PoolId:       1,
				TokenInDenom: "udym",
			},
		}

	// make swap
	resp, err := msgServer.SwapExactAmountOut(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountOut{
		Sender:           suite.TestAccs[0].String(),
		Routes:           routes,
		TokenOut:         sdk.NewCoin("bar", sdk.NewInt(1000)),
		TokenInMaxAmount: sdk.NewInt(100000000),
	})
	suite.Require().NoError(err)
	tokenInCoin := sdk.NewCoin("udym", resp.TokenInAmount)

	// check total supply after swap
	totalSupplyAfter := suite.App.BankKeeper.GetSupply(suite.Ctx, "udym")

	//validate total supply is reduced by taker fee

	_, takerFeeCoin := suite.App.GAMMKeeper.SubTakerFee(tokenInCoin, suite.App.GAMMKeeper.GetParams(ctx).TakerFee)

	suite.Require().True(totalSupplyAfter.Amount.LT(totalSupplyBefore.Amount))
	suite.Require().Equal(takerFeeCoin.Amount, totalSupplyBefore.Amount.Sub(totalSupplyAfter.Amount))
}

func (suite *KeeperTestSuite) TestNonDYMIsSentToCommunity() {
	suite.SetupTest()
	ctx := suite.Ctx
	suite.PrepareBalancerPool()
	msgServer := keeper.NewMsgServerImpl(suite.App.GAMMKeeper)

	//check total supply before swap
	totalSupplyFooBefore := suite.App.BankKeeper.GetSupply(suite.Ctx, "foo")
	totalSupplyDYMBefore := suite.App.BankKeeper.GetSupply(suite.Ctx, "udym")

	// check taker fee is not 0
	suite.Require().True(suite.App.GAMMKeeper.GetParams(ctx).TakerFee.GT(sdk.ZeroDec()))

	routes :=
		[]poolmanagertypes.SwapAmountInRoute{
			{
				PoolId:        1,
				TokenOutDenom: "udym",
			},
		}

	// make swap
	_, err := msgServer.SwapExactAmountIn(sdk.WrapSDKContext(ctx), &types.MsgSwapExactAmountIn{
		Sender:            suite.TestAccs[0].String(),
		Routes:            routes,
		TokenIn:           sdk.NewCoin("foo", sdk.NewInt(100000)),
		TokenOutMinAmount: sdk.NewInt(1),
	})
	suite.Require().NoError(err)

	// check total supply after swap
	totalSupplyFooAfter := suite.App.BankKeeper.GetSupply(suite.Ctx, "foo")
	totalSupplyDYMAfter := suite.App.BankKeeper.GetSupply(suite.Ctx, "udym")

	//validate total supply is NOT reduced by taker fee
	suite.Require().True(totalSupplyFooAfter.Amount.Equal(totalSupplyFooBefore.Amount))
	suite.Require().True(totalSupplyDYMAfter.Amount.Equal(totalSupplyDYMBefore.Amount))

	takerFeeAmount := suite.App.GAMMKeeper.GetParams(ctx).TakerFee.MulInt(sdk.NewInt(100000))

	communityAfter := suite.App.DistrKeeper.GetFeePoolCommunityCoins(ctx)
	suite.Require().True(communityAfter.AmountOf("foo").Equal(takerFeeAmount))
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
			name: "Swap - foo -> udym(pool 1) - udym(pool 2) -> baz ",
			param: param{
				routes: []poolmanagertypes.SwapAmountInRoute{
					{
						PoolId:        1,
						TokenOutDenom: "udym",
					},
					{
						PoolId:        2,
						TokenOutDenom: "baz",
					},
				},
				tokenIn:           sdk.NewCoin("udym", sdk.NewInt(100000)),
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
			name: "Swap - foo -> udym(pool 1) - udym(pool 2) -> baz ",
			param: param{
				routes: []poolmanagertypes.SwapAmountOutRoute{
					{
						PoolId:       1,
						TokenInDenom: "foo",
					},
					{
						PoolId:       2,
						TokenInDenom: "udym",
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
