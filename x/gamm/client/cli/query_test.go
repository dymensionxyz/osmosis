package cli_test

import (
	gocontext "context"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v15/testutils/apptesting"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
)

type QueryTestSuite struct {
	apptesting.KeeperTestHelper
	queryClient types.QueryClient
}

func (s *QueryTestSuite) SetupSuite() {
	s.Setup()
	s.queryClient = types.NewQueryClient(s.QueryHelper)
	// create a new pool
	s.PrepareBalancerPool()
}

func (s *QueryTestSuite) TestQueriesNeverAlterState() {
	var (
		fooDenom  = apptesting.DefaultPoolAssets[0].Token.Denom
		barDenom  = apptesting.DefaultPoolAssets[1].Token.Denom
		bazDenom  = apptesting.DefaultPoolAssets[2].Token.Denom
		adymDenom = apptesting.DefaultPoolAssets[3].Token.Denom

		basicValidTokensIn = sdk.NewCoins(
			sdk.NewCoin(fooDenom, math.OneInt()),
			sdk.NewCoin(barDenom, math.OneInt()),
			sdk.NewCoin(bazDenom, math.OneInt()),
			sdk.NewCoin(adymDenom, math.OneInt()))
	)

	testCases := []struct {
		name   string
		query  string
		input  interface{}
		output interface{}
	}{
		{
			"Query pools",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/Pools",
			&types.QueryPoolsRequest{},
			&types.QueryPoolsResponse{},
		},
		{
			"Query single pool",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/Pool",
			&types.QueryPoolRequest{PoolId: 1},
			&types.QueryPoolsResponse{},
		},
		{
			"Query num pools",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/NumPools",
			&types.QueryNumPoolsRequest{},
			&types.QueryNumPoolsResponse{},
		},
		{
			"Query pool params",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/PoolParams",
			&types.QueryPoolParamsRequest{PoolId: 1},
			&types.QueryPoolParamsResponse{},
		},
		{
			"Query spot price",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/SpotPrice",
			&types.QuerySpotPriceRequest{PoolId: 1, BaseAssetDenom: fooDenom, QuoteAssetDenom: barDenom},
			&types.QuerySpotPriceResponse{},
		},
		{
			"Query total liquidity",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/TotalLiquidity",
			&types.QueryTotalLiquidityRequest{},
			&types.QueryTotalLiquidityResponse{},
		},
		{
			"Query pool total liquidity",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/TotalPoolLiquidity",
			&types.QueryTotalPoolLiquidityRequest{PoolId: 1},
			&types.QueryTotalPoolLiquidityResponse{},
		},
		{
			"Query total shares",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/TotalShares",
			&types.QueryTotalSharesRequest{PoolId: 1},
			&types.QueryTotalSharesResponse{},
		},
		{
			"Query estimate for join pool shares with no swap",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/CalcJoinPoolNoSwapShares",
			&types.QueryCalcJoinPoolNoSwapSharesRequest{PoolId: 1, TokensIn: basicValidTokensIn},
			&types.QueryCalcJoinPoolNoSwapSharesResponse{},
		},
		{
			"Query estimate for join pool shares with no swap",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/CalcJoinPoolShares",
			&types.QueryCalcJoinPoolSharesRequest{PoolId: 1, TokensIn: basicValidTokensIn},
			&types.QueryCalcJoinPoolSharesResponse{},
		},
		{
			"Query exit pool coins from shares",
			"/dymensionxyz.dymension.gamm.v1beta1.Query/CalcExitPoolCoinsFromShares",
			&types.QueryCalcExitPoolCoinsFromSharesRequest{PoolId: 1, ShareInAmount: math.OneInt()},
			&types.QueryCalcExitPoolCoinsFromSharesResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupSuite()
			err := s.QueryHelper.Invoke(gocontext.Background(), tc.query, tc.input, tc.output)
			s.Require().NoError(err)
			s.StateNotAltered()
		})
	}
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}
