package balancer

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ErrMsgFormatRepeatingPoolAssetsNotAllowed = formatRepeatingPoolAssetsNotAllowedErrFormat
	ErrMsgFormatNoPoolAssetFound              = formatNoPoolAssetFoundErrFormat
)

var (
	ErrMsgFormatFailedInterimLiquidityUpdate = failedInterimLiquidityUpdateErrFormat

	CalcPoolSharesOutGivenSingleAssetIn   = calcPoolSharesOutGivenSingleAssetIn
	CalcSingleAssetInGivenPoolSharesOut   = calcSingleAssetInGivenPoolSharesOut
	UpdateIntermediaryPoolAssetsLiquidity = updateIntermediaryPoolAssetsLiquidity

	GetPoolAssetsByDenom = getPoolAssetsByDenom
	EnsureDenomInPool    = ensureDenomInPool
)

func (p *Pool) CalcSingleAssetJoin(tokenIn sdk.Coin, swapFee math.LegacyDec, tokenInPoolAsset PoolAsset, totalShares math.Int) (numShares math.Int, err error) {
	return p.calcSingleAssetJoin(tokenIn, swapFee, tokenInPoolAsset, totalShares)
}

func (p *Pool) CalcJoinSingleAssetTokensIn(tokensIn sdk.Coins, totalSharesSoFar math.Int, poolAssetsByDenom map[string]PoolAsset, swapFee math.LegacyDec) (math.Int, sdk.Coins, error) {
	return p.calcJoinSingleAssetTokensIn(tokensIn, totalSharesSoFar, poolAssetsByDenom, swapFee)
}
