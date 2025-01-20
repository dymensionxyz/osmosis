package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountI defines the account contract that must be fulfilled when
// creating a x/gamm keeper.
type AccountI interface {
	NewAccount(sdk.Context, authtypes.AccountI) authtypes.AccountI
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
}

// BankI defines the banking contract that must be fulfilled when
// creating a x/gamm keeper.
type BankI interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
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
