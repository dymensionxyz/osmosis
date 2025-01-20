package balancer_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/osmosis-labs/osmosis/v15/x/gamm/pool-models/balancer"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/pool-models/internal/test_helpers"
	"github.com/osmosis-labs/osmosis/v15/x/gamm/types"
)

type BalancerTestSuite struct {
	test_helpers.CfmmCommonTestSuite
}

func TestBalancerTestSuite(t *testing.T) {
	suite.Run(t, new(BalancerTestSuite))
}

func TestBalancerPoolParams(t *testing.T) {
	// Tests that creating a pool with the given pair of swapfee and exit fee
	// errors or succeeds as intended. Furthermore, it checks that
	// NewPool panics in the error case.
	tests := []struct {
		SwapFee   math.LegacyDec
		ExitFee   math.LegacyDec
		shouldErr bool
	}{
		// Should work
		{defaultSwapFee, defaultExitFee, noErr},
		// Can't set the swap fee as negative
		{math.LegacyNewDecWithPrec(-1, 2), defaultExitFee, wantErr},
		// Can't set the swap fee as 1
		{math.LegacyNewDec(1), defaultExitFee, wantErr},
		// Can't set the swap fee above 1
		{math.LegacyNewDecWithPrec(15, 1), defaultExitFee, wantErr},
		// Can't set the exit fee as negative
		{defaultSwapFee, math.LegacyNewDecWithPrec(-1, 2), wantErr},
		// Can't set the exit fee as 1
		{defaultSwapFee, math.LegacyNewDec(1), wantErr},
		// Can't set the exit fee above 1
		{defaultSwapFee, math.LegacyNewDecWithPrec(15, 1), wantErr},
	}

	for i, params := range tests {
		PoolParams := balancer.PoolParams{
			SwapFee: params.SwapFee,
			ExitFee: params.ExitFee,
		}
		err := PoolParams.Validate(dummyPoolAssets)
		if params.shouldErr {
			require.Error(t, err, "unexpected lack of error, tc %v", i)
			// Check that these are also caught if passed to the underlying pool creation func
			_, err = balancer.NewBalancerPool(1, PoolParams, dummyPoolAssets, defaultFutureGovernor, defaultCurBlockTime)
			require.Error(t, err)
		} else {
			require.NoError(t, err, "unexpected error, tc %v", i)
		}
	}
}

func (suite *KeeperTestSuite) TestEnsureDenomInPool() {
	tests := map[string]struct {
		poolAssets  []balancer.PoolAsset
		tokensIn    sdk.Coins
		expectPass  bool
		expectedErr error
	}{
		"all of tokensIn is in pool asset map": {
			poolAssets:  []balancer.PoolAsset{defaultOsmoPoolAsset, defaultAtomPoolAsset},
			tokensIn:    sdk.NewCoins(sdk.NewCoin("uatom", math.OneInt())),
			expectPass:  true,
			expectedErr: nil,
		},
		"one of tokensIn is in pool asset map": {
			poolAssets:  []balancer.PoolAsset{defaultOsmoPoolAsset, defaultAtomPoolAsset},
			tokensIn:    sdk.NewCoins(sdk.NewCoin("uatom", math.OneInt()), sdk.NewCoin("foo", math.OneInt())),
			expectPass:  false,
			expectedErr: types.ErrDenomNotFoundInPool,
		},
		"none of tokensIn is in pool asset map": {
			poolAssets:  []balancer.PoolAsset{defaultOsmoPoolAsset, defaultAtomPoolAsset},
			tokensIn:    sdk.NewCoins(sdk.NewCoin("foo", math.OneInt())),
			expectPass:  false,
			expectedErr: types.ErrDenomNotFoundInPool,
		},
	}

	for name, tc := range tests {
		suite.Run(name, func() {
			suite.SetupTest()

			poolAssetsByDenom, err := balancer.GetPoolAssetsByDenom(tc.poolAssets)
			suite.Require().NoError(err, "test: %s", name)

			err = balancer.EnsureDenomInPool(poolAssetsByDenom, tc.tokensIn)

			if tc.expectPass {
				suite.Require().NoError(err, "test: %s", name)
			} else {
				suite.Require().ErrorIs(err, tc.expectedErr, "test: %s", name)
			}
		})
	}
}
