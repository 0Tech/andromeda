package internal

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

type queryServer struct {
	keeper Keeper
}

var _ testv1alpha1.QueryServer = (*queryServer)(nil)

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) testv1alpha1.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

func (s queryServer) Asset(ctx context.Context, req *testv1alpha1.QueryAssetRequest) (*testv1alpha1.QueryAssetResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Account == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil account")
	}

	if req.Asset == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil asset")
	}

	account, err := s.keeper.addressStringToBytes(req.Account)
	if err != nil {
		return nil, errorsmod.Wrap(err, "account")
	}

	if err := s.keeper.HasAsset(ctx, account, req.Asset); err != nil {
		return nil, err
	}

	return &testv1alpha1.QueryAssetResponse{
		Asset: &testv1alpha1.QueryAssetResponse_Asset{
			Name: req.Asset,
		},
	}, nil
}

func (s queryServer) Assets(ctx context.Context, req *testv1alpha1.QueryAssetsRequest) (*testv1alpha1.QueryAssetsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if req.Account == "" {
		return nil, testv1alpha1.ErrUnimplemented.Wrap("nil account")
	}

	account, err := s.keeper.addressStringToBytes(req.Account)
	if err != nil {
		return nil, errorsmod.Wrap(err, "account")
	}

	assets, pageRes, err := query.CollectionPaginate(ctx, s.keeper.assets, req.Pagination, func(key collections.Pair[sdk.AccAddress, string], _ testv1alpha1.Asset) (*testv1alpha1.QueryAssetsResponse_Asset, error) {
		asset := key.K2()

		return &testv1alpha1.QueryAssetsResponse_Asset{
			Name: asset,
		}, nil
	}, query.WithCollectionPaginationPairPrefix[sdk.AccAddress, string](account))
	if err != nil {
		return nil, err
	}

	return &testv1alpha1.QueryAssetsResponse{
		Assets:     assets,
		Pagination: pageRes,
	}, nil
}
