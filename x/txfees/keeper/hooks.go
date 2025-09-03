package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v15/osmoutils"
	epochstypes "github.com/osmosis-labs/osmosis/v15/x/epochs/types"
	gammtypes "github.com/osmosis-labs/osmosis/v15/x/gamm/types"
	pooltypes "github.com/osmosis-labs/osmosis/v15/x/poolmanager/types"
	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

// Hooks is the wrapper struct for the txfees keeper.
type Hooks struct {
	k Keeper
}

var (
	_ epochstypes.EpochHooks = Hooks{}
	_ gammtypes.GammHooks    = Hooks{}
)

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

	// get all balances of this module
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAddr)

	// swap all to dym
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
	var feeToken types.FeeToken

	denoms, err := h.k.gammKeeper.GetPoolDenoms(ctx, poolId)
	if err != nil {
		h.k.Logger(ctx).Error("failed to get pool denoms", "error", err)
		return
	}

	if len(denoms) != 2 {
		h.k.Logger(ctx).Error("expected exactly 2 pool denoms", "denoms", denoms)
		return
	}

	basedenom := h.k.MustGetBaseDenom(ctx)

	// check and handle the case where one of the denoms is basedenom
	// it will override the an existing route if it exists (as it must be a longer path)
	if denoms[0] == basedenom || denoms[1] == basedenom {
		var newDenom string

		if denoms[0] == basedenom {
			newDenom = denoms[1]
		} else {
			newDenom = denoms[0]
		}

		feeToken = types.FeeToken{
			Denom: newDenom,
			Route: []pooltypes.SwapAmountInRoute{
				{
					PoolId:        poolId,
					TokenOutDenom: basedenom,
				},
			},
		}
	} else {
		// no basedenom in the pool, register new token with multi-hop route
		d1Reg := h.k.HasFeeToken(ctx, denoms[0])
		d2Reg := h.k.HasFeeToken(ctx, denoms[1])

		var newDenom, registeredDenom string
		switch {
		case !d1Reg && !d2Reg:
			h.k.Logger(ctx).Error("no route to basedenom exist")
			return
		case d1Reg && d2Reg:
			h.k.Logger(ctx).Debug("both denoms are already registered")
			return
		case d1Reg:
			newDenom, registeredDenom = denoms[1], denoms[0]
		default: // d2Reg
			newDenom, registeredDenom = denoms[0], denoms[1]
		}

		// get the swapRoute for the 2nd pool asset
		feeToken, err = h.k.GetFeeToken(ctx, registeredDenom)
		if err != nil {
			h.k.Logger(ctx).Error("failed to get fee token", "error", err)
			return
		}
		var route []pooltypes.SwapAmountInRoute
		route = append(route, pooltypes.SwapAmountInRoute{
			PoolId:        poolId,
			TokenOutDenom: registeredDenom,
		})
		route = append(route, feeToken.Route...)

		feeToken = types.FeeToken{
			Denom: newDenom,
			Route: route,
		}
	}

	err = h.k.SetFeeToken(ctx, feeToken)
	if err != nil {
		h.k.Logger(ctx).Error("failed to set fee token", "error", err)
		return
	}

	h.k.Logger(ctx).Info("created fee token route for new denom",
		"denom", feeToken.Denom, "poolId", poolId, "routeLength", len(feeToken.Route))
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

// AfterReplacePoolAsset is called after an asset in a pool is replaced.
// It updates the fee token routes by removing the route for the old denom
// and creating a new route for the new denom if applicable.
func (h Hooks) AfterReplacePoolAsset(ctx sdk.Context, poolId uint64, oldDenom, newDenom string) {
	baseDenom := h.k.MustGetBaseDenom(ctx)

	// Remove the old fee token route if it exists
	h.k.DeleteFeeToken(ctx, oldDenom)

	// Get current pool denoms to understand the pool structure
	denoms, err := h.k.gammKeeper.GetPoolDenoms(ctx, poolId)
	if err != nil {
		h.k.Logger(ctx).Error("failed to get pool denoms for new fee token", "poolId", poolId, "error", err)
		return
	}

	// Find the other denom in the pool (not the new token)
	var otherDenom string
	for _, denom := range denoms {
		if denom != newDenom {
			otherDenom = denom
			break
		}
	}

	if otherDenom == "" {
		h.k.Logger(ctx).Error("could not find other denom in pool", "poolId", poolId, "newDenom", newDenom)
		return
	}

	var feeToken types.FeeToken

	// Case 1: Other denom is base denom - direct route
	if otherDenom == baseDenom {
		feeToken = types.FeeToken{
			Denom: newDenom,
			Route: []pooltypes.SwapAmountInRoute{
				{
					PoolId:        poolId,
					TokenOutDenom: baseDenom,
				},
			},
		}
	} else {
		// Case 2: Other denom is a registered fee token - multi-hop route
		registeredFeeToken, err := h.k.GetFeeToken(ctx, otherDenom)
		if err != nil {
			h.k.Logger(ctx).Error("failed to get registered fee token", "denom", otherDenom, "error", err)
			return
		}

		// Create route: newToken -> otherDenom -> ... -> baseDenom
		var route []pooltypes.SwapAmountInRoute
		route = append(route, pooltypes.SwapAmountInRoute{
			PoolId:        poolId,
			TokenOutDenom: otherDenom,
		})
		route = append(route, registeredFeeToken.Route...)

		feeToken = types.FeeToken{
			Denom: newDenom,
			Route: route,
		}
	}

	// Set the new fee token
	err = h.k.SetFeeToken(ctx, feeToken)
	if err != nil {
		h.k.Logger(ctx).Error("failed to set new fee token route", "denom", newDenom, "error", err)
		return
	}

	h.k.Logger(ctx).Info("created fee token route for new denom",
		"denom", newDenom, "poolId", poolId, "routeLength", len(feeToken.Route))
}
