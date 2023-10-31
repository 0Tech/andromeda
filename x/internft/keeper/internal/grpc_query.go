package internal

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"

	"github.com/cosmos/cosmos-sdk/types/query"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

type queryServer struct {
	keeper Keeper
}

var _ internftv1alpha1.QueryServer = (*queryServer)(nil)

// NewQueryServer returns an implementation of the QueryServer interface
// for the provided Keeper.
func NewQueryServer(keeper Keeper) internftv1alpha1.QueryServer {
	return &queryServer{
		keeper: keeper,
	}
}

// Params queries the module params.
func (s queryServer) Params(ctx context.Context, req *internftv1alpha1.QueryParamsRequest) (*internftv1alpha1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	params := s.keeper.GetParams(ctx)

	return &internftv1alpha1.QueryParamsResponse{
		Params: params,
	}, nil
}

// Class queries a class.
func (s queryServer) Class(ctx context.Context, req *internftv1alpha1.QueryClassRequest) (*internftv1alpha1.QueryClassResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if err := internftv1alpha1.ValidateClassID(req.ClassId); err != nil {
		return nil, err
	}

	class, err := s.keeper.GetClass(ctx, req.ClassId)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryClassResponse{
		Class: class,
	}, nil
}

// Classes queries all classes.
func (s queryServer) Classes(ctx context.Context, req *internftv1alpha1.QueryClassesRequest) (*internftv1alpha1.QueryClassesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	classes, pageRes, err := query.CollectionPaginate(ctx, s.keeper.classes, req.Pagination, func(_ string, value internftv1alpha1.Class) (internftv1alpha1.Class, error) {
		return value, nil
	})
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryClassesResponse{
		Classes:    classes,
		Pagination: pageRes,
	}, nil
}

// Trait queries a trait of a class.
func (s queryServer) Trait(ctx context.Context, req *internftv1alpha1.QueryTraitRequest) (*internftv1alpha1.QueryTraitResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if err := internftv1alpha1.ValidateClassID(req.ClassId); err != nil {
		return nil, err
	}

	if err := internftv1alpha1.ValidateTraitID(req.TraitId); err != nil {
		return nil, err
	}

	trait, err := s.keeper.GetTrait(ctx, req.ClassId, req.TraitId)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryTraitResponse{
		Trait: trait,
	}, nil
}

// Traits queries all traits of a class.
func (s queryServer) Traits(ctx context.Context, req *internftv1alpha1.QueryTraitsRequest) (*internftv1alpha1.QueryTraitsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if err := internftv1alpha1.ValidateClassID(req.ClassId); err != nil {
		return nil, err
	}

	traits, pageRes, err := query.CollectionPaginate(ctx, s.keeper.traits, req.Pagination, func(_ collections.Pair[string, string], value internftv1alpha1.Trait) (internftv1alpha1.Trait, error) {
		return value, nil
	}, func(o *query.CollectionsPaginateOptions[collections.Pair[string, string]]) {
		prefix := collections.PairPrefix[string, string](req.ClassId)
		o.Prefix = &prefix
	})
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryTraitsResponse{
		Traits:     traits,
		Pagination: pageRes,
	}, nil
}

// NFT queries an nft.
func (s queryServer) NFT(ctx context.Context, req *internftv1alpha1.QueryNFTRequest) (*internftv1alpha1.QueryNFTResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	nft := internftv1alpha1.NFT{
		ClassId: req.ClassId,
		Id: req.NftId,
	}
	if err := nft.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := s.keeper.hasNFT(ctx, nft); err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryNFTResponse{
		Nft: &nft,
	}, nil
}

// NFTs queries all nfts.
func (s queryServer) NFTs(ctx context.Context, req *internftv1alpha1.QueryNFTsRequest) (*internftv1alpha1.QueryNFTsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	if err := internftv1alpha1.ValidateClassID(req.ClassId); err != nil {
		return nil, err
	}

	nfts, pageRes, err := query.CollectionPaginate(ctx, s.keeper.nfts, req.Pagination, func(_ collections.Pair[string, string], value internftv1alpha1.NFT) (internftv1alpha1.NFT, error) {
		return value, nil
	}, func(o *query.CollectionsPaginateOptions[collections.Pair[string, string]]) {
		prefix := collections.PairPrefix[string, string](req.ClassId)
		o.Prefix = &prefix
	})
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryNFTsResponse{
		Nfts:       nfts,
		Pagination: pageRes,
	}, nil
}

// Property queries a property of a class.
func (s queryServer) Property(ctx context.Context, req *internftv1alpha1.QueryPropertyRequest) (*internftv1alpha1.QueryPropertyResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	nft := internftv1alpha1.NFT{
		ClassId: req.ClassId,
		Id: req.NftId,
	}
	if err := nft.ValidateBasic(); err != nil {
		return nil, err
	}

	if err := internftv1alpha1.ValidateTraitID(req.PropertyId); err != nil {
		return nil, err
	}

	property, err := s.keeper.GetProperty(ctx, nft, req.PropertyId)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryPropertyResponse{
		Property: property,
	}, nil
}

// Properties queries all properties of a class.
func (s queryServer) Properties(ctx context.Context, req *internftv1alpha1.QueryPropertiesRequest) (*internftv1alpha1.QueryPropertiesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	nft := internftv1alpha1.NFT{
		ClassId: req.ClassId,
		Id: req.NftId,
	}
	if err := nft.ValidateBasic(); err != nil {
		return nil, err
	}

	properties, pageRes, err := query.CollectionPaginate(ctx, s.keeper.properties, req.Pagination, func(_ collections.Triple[string, string, string], value internftv1alpha1.Property) (internftv1alpha1.Property, error) {
		return value, nil
	}, func(o *query.CollectionsPaginateOptions[collections.Triple[string, string, string]]) {
		prefix := collections.TripleSuperPrefix[string, string, string](nft.ClassId, nft.Id)
		o.Prefix = &prefix
	})
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryPropertiesResponse{
		Properties: properties,
		Pagination: pageRes,
	}, nil
}

// Owner queries the owner of an nft.
func (s queryServer) Owner(ctx context.Context, req *internftv1alpha1.QueryOwnerRequest) (*internftv1alpha1.QueryOwnerResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	nft := internftv1alpha1.NFT{
		ClassId: req.ClassId,
		Id: req.NftId,
	}
	if err := nft.ValidateBasic(); err != nil {
		return nil, err
	}

	owner, err := s.keeper.GetOwner(ctx, nft)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryOwnerResponse{
		Owner: owner.String(),
	}, nil
}
