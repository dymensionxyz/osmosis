package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/osmosis-labs/osmosis/v15/x/txfees/types"
)

var _ types.QueryServer = Querier{}

// Querier defines a wrapper around the x/txfees keeper providing gRPC method
// handlers.
type Querier struct {
	Keeper
}

func NewQuerier(k Keeper) Querier {
	return Querier{Keeper: k}
}

func (q Querier) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := q.Keeper.GetParams(sdkCtx)

	return &types.QueryParamsResponse{Params: params}, nil
}

func (q Querier) FeeTokens(ctx context.Context, _ *types.QueryFeeTokensRequest) (*types.QueryFeeTokensResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	feeTokens := q.Keeper.GetFeeTokens(sdkCtx)

	return &types.QueryFeeTokensResponse{FeeTokens: feeTokens}, nil
}

func (q Querier) DenomSpotPrice(ctx context.Context, req *types.QueryDenomSpotPriceRequest) (*types.QueryDenomSpotPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if len(req.Denom) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty denom")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	baseDenom, err := q.GetBaseDenom(sdkCtx)
	if err != nil {
		return nil, err
	}

	feeToken, err := q.GetFeeToken(sdkCtx, req.Denom)
	if err != nil {
		return nil, err
	}

	spotPrice, err := q.spotPriceCalculator.CalculateSpotPrice(sdkCtx, feeToken.PoolID, baseDenom, feeToken.Denom)
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomSpotPriceResponse{PoolID: feeToken.PoolID, SpotPrice: spotPrice}, nil
}

func (q Querier) DenomPoolId(ctx context.Context, req *types.QueryDenomPoolIdRequest) (*types.QueryDenomPoolIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if len(req.Denom) == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty denom")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	feeToken, err := q.Keeper.GetFeeToken(sdkCtx, req.GetDenom())
	if err != nil {
		return nil, err
	}

	return &types.QueryDenomPoolIdResponse{PoolID: feeToken.GetPoolID()}, nil
}

func (q Querier) BaseDenom(ctx context.Context, _ *types.QueryBaseDenomRequest) (*types.QueryBaseDenomResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	baseDenom, err := q.Keeper.GetBaseDenom(sdkCtx)
	if err != nil {
		return nil, err
	}

	return &types.QueryBaseDenomResponse{BaseDenom: baseDenom}, nil
}
