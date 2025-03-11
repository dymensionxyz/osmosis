package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/osmoutils"
	epochstypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"
	pooltypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	poolmanagertypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

// Hooks is the wrapper struct for the txfees keeper.
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}
var _ gammtypes.GammHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

/* -------------------------------------------------------------------------- */
/*                                 epoch hooks                                */
/* -------------------------------------------------------------------------- */
func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}

// at the end of each epoch, swap all non-DYM fees into DYM and burn them
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	if epochIdentifier != k.GetParams(ctx).EpochIdentifier {
		return nil
	}

	moduleAddr := k.accountKeeper.GetModuleAddress(types.ModuleName)
	baseDenom, _ := k.GetBaseDenom(ctx)

	//get all balances of this module
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAddr)

	//swap all to dym
	for _, coinBalance := range balances {
		if coinBalance.Denom == baseDenom {
			continue
		}
		if coinBalance.Amount.IsZero() {
			continue
		}

		feetoken, err := k.GetFeeToken(ctx, coinBalance.Denom)
		if err != nil {
			k.Logger(ctx).Error("unknown fee token", "denom", coinBalance.Denom, "error", err)
			err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coinBalance))
			if err != nil {
				k.Logger(ctx).Error("failed to burn non-native coins", "error", err)
			}
			continue
		}

		// Do the swap of this fee token denom to base denom.
		wrappedRouteExactAmountInFn := func(ctx sdk.Context) error {
			_, err := k.poolManager.RouteExactAmountIn(ctx, moduleAddr, feetoken.Route, coinBalance, math.ZeroInt())
			return err
		}
		err = osmoutils.ApplyFuncIfNoError(ctx, wrappedRouteExactAmountInFn)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed to swap fee token to base token: %v. Trying to burn the tokens", err))
			err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coinBalance))
			if err != nil {
				k.Logger(ctx).Error("failed to burn non-native coins", "error", err)
			}
		}
	}

	// Get all of the txfee payout denom in the module account
	baseDenomCoins := sdk.NewCoins(k.bankKeeper.GetBalance(ctx, moduleAddr, baseDenom))
	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, baseDenomCoins)
	if err != nil {
		return err
	}

	return nil
}

func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

/* -------------------------------------------------------------------------- */
/*                                 pool hooks                                 */
/* -------------------------------------------------------------------------- */
// AfterPoolCreated is called after CreatePool.
// It checks if the base denom is included in the newly created pool.
// If so, it adds the non-native denom as a fee token.
func (h Hooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64) {
	denoms, err := h.k.spotPriceCalculator.GetPoolDenoms(ctx, poolId)
	if err != nil {
		h.k.Logger(ctx).Error("failed to get pool denoms", "error", err)
		return
	}

	newDenom, registeredDenom, err := h.GetNotRegisteredDenom(ctx, denoms)
	if err != nil {
		h.k.Logger(ctx).Error("failed to get non-registered denom", "error", err)
		return
	}

	// get the swapRoute for the 2nd pool asset
	var route []pooltypes.SwapAmountInRoute

	if registeredDenom == h.k.MustGetBaseDenom(ctx) {
		route = append(route, pooltypes.SwapAmountInRoute{
			PoolId:        poolId,
			TokenOutDenom: registeredDenom,
		})
	} else {
		feeToken, err := h.k.GetFeeToken(ctx, newDenom)
		if err != nil {
			h.k.Logger(ctx).Error("failed to get fee token", "error", err)
		return
		}
		route = feeToken.Route
		route = append(route, pooltypes.SwapAmountInRoute{
			PoolId:        poolId,
			TokenOutDenom: registeredDenom,
		})
	}

	feeToken := types.FeeToken{
		Denom: newDenom,
		Route: route,
	}

	// validate the route between feeToken and baseDenom
	// FIXME:

	err = h.k.SetFeeToken(ctx, feeToken)
	if err != nil {
		h.k.Logger(ctx).Error("failed to set fee token", "error", err)
		return
	}

	return
}

// AfterJoinPool hook is a noop.
func (h Hooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, enterCoins sdk.Coins, shareOutAmount math.Int) {
}

// AfterExitPool hook is a noop.
func (h Hooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount math.Int, exitCoins sdk.Coins) {
}

// AfterSwap hook is a noop.
func (h Hooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) {
}

func contains(strarr []string, str string) bool {
	for _, v := range strarr {
		if v == str {
			return true
		}
	}

	return false
}

// GetNotRegisteredDenom returns the non-registered denom in the pool that is neither base denom nor registered fee token
func (h Hooks) GetNotRegisteredDenom(ctx sdk.Context, denoms []string) (string, string, error) {
	if len(denoms) != 2 {
		return "", "", fmt.Errorf("expected exactly 2 pool denoms, got %d: %v", len(denoms), denoms)
	}

	d1Reg := h.k.IsRegisteredDenom(ctx, denoms[0])
	d2Reg := h.k.IsRegisteredDenom(ctx, denoms[1])

	switch {
	case !d1Reg && !d2Reg:
		return "", "", errors.New("both denoms are unregistered")
	case d1Reg && d2Reg:
		return "", "", errors.New("both denoms are already registered")
	case d1Reg:
		return denoms[1], denoms[0], nil
	default: // d2Reg
		return denoms[0], denoms[1], nil
	}
}

func (k Keeper) IsRegisteredDenom(ctx sdk.Context, denom string) bool {
	return k.MustGetBaseDenom(ctx) == denom || k.HasFeeToken(ctx, denom)
}

// getOtherDenom returns the other denom in the pool that is not the base denom
// assumes that the pool has only 2 denoms
func getOtherDenom(denoms []string, idx uint64) string {
	if idx == 0 {
		return denoms[1]
	}
	return denoms[0]
}
