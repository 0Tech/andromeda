package internal_test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestQueryParams() {
	testCases := map[string]struct {
		code codes.Code
	}{
		"valid request": {},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryParamsRequest{}

			res, err := s.queryServer.Params(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			// params := res.Params
		})
	}
}

func (s *KeeperTestSuite) TestQueryClass() {
	testCases := map[string]struct {
		classID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"invalid class id": {
			code: codes.InvalidArgument,
		},
		"class not found": {
			classID: s.customer.String(),
			code:    codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryClassRequest{
				ClassId: tc.classID,
			}

			res, err := s.queryServer.Class(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			class := res.Class
			s.Require().NotNil(class)
			s.Require().Equal(tc.classID, class.Id)
		})
	}
}

func (s *KeeperTestSuite) TestQueryClasses() {
	testCases := map[string]struct {
		code codes.Code
	}{
		"valid request": {},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryClassesRequest{}

			res, err := s.queryServer.Classes(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			classes := res.Classes
			s.Require().Len(classes, 1)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTrait() {
	testCases := map[string]struct {
		classID string
		traitID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
			traitID: s.immutableTraitID,
		},
		"invalid class id": {
			classID: "not-in-bech32",
			traitID: s.immutableTraitID,
			code:    codes.InvalidArgument,
		},
		// "invalid trait id": {
		// 	classID: s.vendor.String(),
		// 	code:    codes.InvalidArgument,
		// },
		"trait not found": {
			classID: s.customer.String(),
			traitID: s.immutableTraitID,
			code:    codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryTraitRequest{
				ClassId: tc.classID,
				TraitId: tc.traitID,
			}

			res, err := s.queryServer.Trait(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			trait := res.Trait
			s.Require().NotNil(trait)
			s.Require().Equal(tc.traitID, trait.Id)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTraits() {
	testCases := map[string]struct {
		classID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"invalid class id": {
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryTraitsRequest{
				ClassId: tc.classID,
			}

			res, err := s.queryServer.Traits(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			traits := res.Traits
			s.Require().Len(traits, 2)
		})
	}
}

func (s *KeeperTestSuite) TestQueryToken() {
	testCases := map[string]struct {
		classID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"invalid class id": {
			classID: "not-in-bech32",
			code: codes.InvalidArgument,
		},
		"token not found": {
			classID: s.customer.String(),
			code:    codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryTokenRequest{
				ClassId: tc.classID,
				TokenId:      s.tokenIDs[s.customer.String()],
			}

			res, err := s.queryServer.Token(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			token := res.Token
			s.Require().NotNil(token)
			s.Require().Equal(req.ClassId, token.ClassId)
			s.Require().Equal(req.TokenId, token.Id)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokens() {
	testCases := map[string]struct {
		classID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"invalid class id": {
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryTokensRequest{
				ClassId: tc.classID,
			}

			res, err := s.queryServer.Tokens(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			tokens := res.Tokens
			s.Require().Len(tokens, len(s.tokenIDs)-1)
		})
	}
}

func (s *KeeperTestSuite) TestQueryProperty() {
	testCases := map[string]struct {
		classID    string
		traitID string
		code       codes.Code
	}{
		"valid request": {
			classID:    s.vendor.String(),
			traitID: s.immutableTraitID,
		},
		"invalid class id": {
			classID: "not-in-bech32",
			traitID: s.immutableTraitID,
			code:       codes.InvalidArgument,
		},
		// "invalid trait id": {
		// 	classID: s.vendor.String(),
		// 	code:    codes.InvalidArgument,
		// },
		"trait not found": {
			classID:    s.customer.String(),
			traitID: s.immutableTraitID,
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryPropertyRequest{
				ClassId:    tc.classID,
				TokenId:      s.tokenIDs[s.customer.String()],
				TraitId: tc.traitID,
			}

			res, err := s.queryServer.Property(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			property := res.Property
			s.Require().NotNil(property)
			s.Require().Equal(tc.traitID, property.TraitId)
		})
	}
}

func (s *KeeperTestSuite) TestQueryProperties() {
	testCases := map[string]struct {
		classID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"invalid class id": {
			classID: "not-in-bech32",
			code: codes.InvalidArgument,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryPropertiesRequest{
				ClassId: tc.classID,
				TokenId:   s.tokenIDs[s.customer.String()],
			}

			res, err := s.queryServer.Properties(ctx, req)
			s.Require().Equal(tc.code, status.Code(err))
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			properties := res.Properties
			s.Require().Len(properties, 2)
		})
	}
}

func (s *KeeperTestSuite) TestQueryOwner() {
	testCases := map[string]struct {
		classID string
		code    codes.Code
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"invalid class id": {
			classID: "not-in-bech32",
			code: codes.InvalidArgument,
		},
		"token not found": {
			classID: s.customer.String(),
			code:    codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.QueryOwnerRequest{
				ClassId: tc.classID,
				TokenId:   s.tokenIDs[s.customer.String()],
			}

			res, err := s.queryServer.Owner(ctx, req)
			s.Assert().Equal(tc.code, status.Code(err), err)
			if tc.code != codes.OK {
				return
			}
			s.Require().NotNil(res)

			ownerStr := res.Owner
			_, err = sdk.AccAddressFromBech32(ownerStr)
			s.Require().NoError(err)
		})
	}
}
