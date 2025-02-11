package types

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountI defines the account contract that must be fulfilled when
// creating a x/gamm keeper.
type AccountI interface {
	NewAccount(context.Context, sdk.AccountI) sdk.AccountI
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
	SetAccount(ctx context.Context, acc sdk.AccountI)
	SetModuleAccount(ctx context.Context, macc sdk.ModuleAccountI)
}

// BankI defines the banking contract that must be fulfilled when
// creating a x/gamm keeper.
type BankI interface {
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// TODO: godoc
type SwapI interface {
	InitializePool(ctx sdk.Context, pool PoolI, creatorAddress sdk.AccAddress) error

	GetPool(ctx sdk.Context, poolId uint64) (PoolI, error)

	SwapExactAmountIn(
		ctx sdk.Context,
		sender sdk.AccAddress,
		pool PoolI,
		tokenIn sdk.Coin,
		tokenOutDenom string,
		tokenOutMinAmount math.Int,
		swapFee math.LegacyDec,
	) (math.Int, error)
	// CalcOutAmtGivenIn calculates the amount of tokenOut given tokenIn and the pool's current state.
	// Returns error if the given pool is not a CFMM pool. Returns error on internal calculations.
	CalcOutAmtGivenIn(
		ctx sdk.Context,
		poolI PoolI,
		tokenIn sdk.Coin,
		tokenOutDenom string,
		swapFee math.LegacyDec,
	) (tokenOut sdk.Coin, err error)

	SwapExactAmountOut(
		ctx sdk.Context,
		sender sdk.AccAddress,
		pool PoolI,
		tokenInDenom string,
		tokenInMaxAmount math.Int,
		tokenOut sdk.Coin,
		swapFee math.LegacyDec,
	) (tokenInAmount math.Int, err error)
	// CalcInAmtGivenOut calculates the amount of tokenIn given tokenOut and the pool's current state.
	// Returns error if the given pool is not a CFMM pool. Returns error on internal calculations.
	CalcInAmtGivenOut(
		ctx sdk.Context,
		poolI PoolI,
		tokenOut sdk.Coin,
		tokenInDenom string,
		swapFee math.LegacyDec,
	) (tokenIn sdk.Coin, err error)
}
