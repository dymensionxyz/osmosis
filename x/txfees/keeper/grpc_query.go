package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
func (k Keeper) DenomRoute(goCtx context.Context, req *types.QueryDenomRouteRequest) (*types.QueryDenomRouteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeToken, err := k.GetFeeToken(ctx, req.Denom)
	if err != nil {
		return nil, status.Error(codes.NotFound, "fee token not found")
	}

	return &types.QueryDenomRouteResponse{
		Route: feeToken.Route,
	}, nil
}
func (k Keeper) FeeToken(goCtx context.Context, req *types.QueryFeeTokenRequest) (*types.QueryFeeTokenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeToken, err := k.GetFeeToken(ctx, req.Denom)
	if err != nil {
		return nil, status.Error(codes.NotFound, "fee token not found")
	}

	return &types.QueryFeeTokenResponse{
		FeeToken: feeToken,
	}, nil
}
func (q Querier) BaseDenom(ctx context.Context, _ *types.QueryBaseDenomRequest) (*types.QueryBaseDenomResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	baseDenom, err := q.Keeper.GetBaseDenom(sdkCtx)
	if err != nil {
		return nil, err
	}

	return &types.QueryBaseDenomResponse{BaseDenom: baseDenom}, nil
}
