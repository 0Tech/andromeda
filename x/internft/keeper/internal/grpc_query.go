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

func (s queryServer) Params(ctx context.Context, req *internftv1alpha1.QueryParamsRequest) (*internftv1alpha1.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryParamsInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	params, err := s.keeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryParamsResponse{
		Params: params,
	}, nil
}

func (s queryServer) Class(ctx context.Context, req *internftv1alpha1.QueryClassRequest) (*internftv1alpha1.QueryClassResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryClassInternal
	if err := parsed.Parse(*req); err != nil {
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

func (s queryServer) Classes(ctx context.Context, req *internftv1alpha1.QueryClassesRequest) (*internftv1alpha1.QueryClassesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryClassesInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	classes, pageRes, err := query.CollectionPaginate(ctx, s.keeper.classes, req.Pagination, func(_ string, value internftv1alpha1.Class) (*internftv1alpha1.Class, error) {
		return &value, nil
	})
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryClassesResponse{
		Classes:    classes,
		Pagination: pageRes,
	}, nil
}

func (s queryServer) Trait(ctx context.Context, req *internftv1alpha1.QueryTraitRequest) (*internftv1alpha1.QueryTraitResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryTraitInternal
	if err := parsed.Parse(*req); err != nil {
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

func (s queryServer) Traits(ctx context.Context, req *internftv1alpha1.QueryTraitsRequest) (*internftv1alpha1.QueryTraitsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryTraitsInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	traits, pageRes, err := query.CollectionPaginate(ctx, s.keeper.traits, req.Pagination, func(_ collections.Pair[string, string], value internftv1alpha1.Trait) (*internftv1alpha1.Trait, error) {
		return &value, nil
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

func (s queryServer) Token(ctx context.Context, req *internftv1alpha1.QueryTokenRequest) (*internftv1alpha1.QueryTokenResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryTokenInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	token := &internftv1alpha1.Token{
		ClassId: req.ClassId,
		Id: req.TokenId,
	}

	if err := s.keeper.hasToken(ctx, token); err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryTokenResponse{
		Token: token,
	}, nil
}

func (s queryServer) Tokens(ctx context.Context, req *internftv1alpha1.QueryTokensRequest) (*internftv1alpha1.QueryTokensResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryTokensInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	tokens, pageRes, err := query.CollectionPaginate(ctx, s.keeper.tokens, req.Pagination, func(_ collections.Pair[string, string], value internftv1alpha1.Token) (*internftv1alpha1.Token, error) {
		return &value, nil
	}, func(o *query.CollectionsPaginateOptions[collections.Pair[string, string]]) {
		prefix := collections.PairPrefix[string, string](req.ClassId)
		o.Prefix = &prefix
	})
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryTokensResponse{
		Tokens:       tokens,
		Pagination: pageRes,
	}, nil
}

func (s queryServer) Property(ctx context.Context, req *internftv1alpha1.QueryPropertyRequest) (*internftv1alpha1.QueryPropertyResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryPropertyInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	token := &internftv1alpha1.Token{
		ClassId: req.ClassId,
		Id: req.TokenId,
	}

	property, err := s.keeper.GetProperty(ctx, token, req.TraitId)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryPropertyResponse{
		Property: property,
	}, nil
}

func (s queryServer) Properties(ctx context.Context, req *internftv1alpha1.QueryPropertiesRequest) (*internftv1alpha1.QueryPropertiesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryPropertiesInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	token := internftv1alpha1.Token{
		ClassId: req.ClassId,
		Id: req.TokenId,
	}

	properties, pageRes, err := query.CollectionPaginate(ctx, s.keeper.properties, req.Pagination, func(_ collections.Triple[string, string, string], value internftv1alpha1.Property) (*internftv1alpha1.Property, error) {
		return &value, nil
	}, func(o *query.CollectionsPaginateOptions[collections.Triple[string, string, string]]) {
		prefix := collections.TripleSuperPrefix[string, string, string](token.ClassId, token.Id)
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

func (s queryServer) Owner(ctx context.Context, req *internftv1alpha1.QueryOwnerRequest) (*internftv1alpha1.QueryOwnerResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "empty request")
	}

	var parsed internftv1alpha1.QueryOwnerInternal
	if err := parsed.Parse(*req); err != nil {
		return nil, err
	}

	token := &internftv1alpha1.Token{
		ClassId: req.ClassId,
		Id: req.TokenId,
	}

	owner, err := s.keeper.GetOwner(ctx, token)
	if err != nil {
		return nil, err
	}

	return &internftv1alpha1.QueryOwnerResponse{
		Owner: owner.String(),
	}, nil
}
