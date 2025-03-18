package keeper

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetFeeToken returns the fee token record for a specific denom,
// In our case the baseDenom is adym.
func (k Keeper) GetBaseDenom(ctx sdk.Context) (denom string, err error) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(types.BaseDenomKey) {
		return "", types.ErrNoBaseDenom
	}

	bz := store.Get(types.BaseDenomKey)

	return string(bz), nil
}

// MustGetBaseDenom returns the baseDenom or panics
func (k Keeper) MustGetBaseDenom(ctx sdk.Context) string {
	denom, err := k.GetBaseDenom(ctx)
	if err != nil {
		panic(err)
	}
	return denom
}

// SetBaseDenom sets the base fee denom for the chain. Should only be used once.
func (k Keeper) SetBaseDenom(ctx sdk.Context, denom string) error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.BaseDenomKey, []byte(denom))
	return nil
}

// HasFeeToken checks if a fee token record exists for a specific denom.
func (k Keeper) HasFeeToken(ctx sdk.Context, denom string) bool {
	prefixStore := k.getFeeTokensStore(ctx)
	return prefixStore.Has([]byte(denom))
}

// GetFeeToken returns a unique fee token record for a specific denom.
// If the denom doesn't exist, returns an error.
func (k Keeper) GetFeeToken(ctx sdk.Context, denom string) (types.FeeToken, error) {
	prefixStore := k.getFeeTokensStore(ctx)
	if !prefixStore.Has([]byte(denom)) {
		return types.FeeToken{}, fmt.Errorf("denom not found (%s)", denom)
	}
	bz := prefixStore.Get([]byte(denom))

	feeToken := types.FeeToken{}
	err := proto.Unmarshal(bz, &feeToken)
	if err != nil {
		return types.FeeToken{}, err
	}

	return feeToken, nil
}

func (k Keeper) GetFeeTokens(ctx sdk.Context) (feetokens []types.FeeToken) {
	prefixStore := k.getFeeTokensStore(ctx)

	// this entire store just contains FeeTokens, so iterate over all entries.
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	feeTokens := []types.FeeToken{}

	for ; iterator.Valid(); iterator.Next() {
		feeToken := types.FeeToken{}

		err := proto.Unmarshal(iterator.Value(), &feeToken)
		if err != nil {
			panic(err)
		}

		feeTokens = append(feeTokens, feeToken)
	}
	return feeTokens
}

func (k Keeper) SetFeeTokens(ctx sdk.Context, feetokens []types.FeeToken) error {
	for _, feeToken := range feetokens {
		err := k.SetFeeToken(ctx, feeToken)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetFeeToken sets a new fee token record for a specific denom.
// PoolID is just the pool to swap rate between alt fee token and native fee token.
// If the len of the feeToken route is 0, deletes the fee Token entry.
func (k Keeper) SetFeeToken(ctx sdk.Context, feeToken types.FeeToken) error {
	prefixStore := k.getFeeTokensStore(ctx)

	if len(feeToken.Route) == 0 {
		if prefixStore.Has([]byte(feeToken.Denom)) {
			prefixStore.Delete([]byte(feeToken.Denom))
		}
		return nil
	}

	bz, err := proto.Marshal(&feeToken)
	if err != nil {
		return err
	}

	prefixStore.Set([]byte(feeToken.Denom), bz)
	return nil
}
