package internal_test

import (
	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestQueryParams() {
	tester := func(subject internftv1alpha1.QueryParamsRequest) error {
		res, err := s.queryServer.Params(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Params)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryParamsRequest]{
		{
			"valid request": {
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryClass() {
	tester := func(subject internftv1alpha1.QueryClassRequest) error {
		res, err := s.queryServer.Class(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Class)
		s.Require().Equal(subject.ClassId, res.Class.Id)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryClassRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryClassRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryClassRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryClasses() {
	tester := func(subject internftv1alpha1.QueryClassesRequest) error {
		res, err := s.queryServer.Classes(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Classes, 1)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryClassesRequest]{
		{
			"valid request": {
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryTrait() {
	tester := func(subject internftv1alpha1.QueryTraitRequest) error {
		res, err := s.queryServer.Trait(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Trait)
		s.Require().Equal(subject.TraitId, res.Trait.Id)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryTraitRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryTraitRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryTraitRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil trait id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid trait id": {
				malleate: func(subject *internftv1alpha1.QueryTraitRequest) {
					subject.TraitId = s.immutableTraitID
				},
			},
			"invalid trait id": {
				malleate: func(subject *internftv1alpha1.QueryTraitRequest) {
					subject.TraitId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryTraits() {
	tester := func(subject internftv1alpha1.QueryTraitsRequest) error {
		res, err := s.queryServer.Traits(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Traits, 2)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryTraitsRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryTraitsRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryTraitsRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryToken() {
	tester := func(subject internftv1alpha1.QueryTokenRequest) error {
		res, err := s.queryServer.Token(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Token)
		s.Require().Equal(subject.ClassId, res.Token.ClassId)
		s.Require().Equal(subject.TokenId, res.Token.Id)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryTokenRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryTokenRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryTokenRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id": {
				malleate: func(subject *internftv1alpha1.QueryTokenRequest) {
					subject.TokenId = s.tokenIDs[s.customer.String()]
				},
			},
			"invalid token id": {
				malleate: func(subject *internftv1alpha1.QueryTokenRequest) {
					subject.TokenId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryTokens() {
	tester := func(subject internftv1alpha1.QueryTokensRequest) error {
		res, err := s.queryServer.Tokens(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Tokens, len(s.tokenIDs) - 1)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryTokensRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryTokensRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryTokensRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryProperty() {
	tester := func(subject internftv1alpha1.QueryPropertyRequest) error {
		res, err := s.queryServer.Property(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Property)
		s.Require().Equal(subject.TraitId, res.Property.TraitId)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryPropertyRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryPropertyRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryPropertyRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id": {
				malleate: func(subject *internftv1alpha1.QueryPropertyRequest) {
					subject.TokenId = s.tokenIDs[s.customer.String()]
				},
			},
			"invalid token id": {
				malleate: func(subject *internftv1alpha1.QueryPropertyRequest) {
					subject.TokenId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil trait id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid trait id": {
				malleate: func(subject *internftv1alpha1.QueryPropertyRequest) {
					subject.TraitId = s.immutableTraitID
				},
			},
			"invalid trait id": {
				malleate: func(subject *internftv1alpha1.QueryPropertyRequest) {
					subject.TraitId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryProperties() {
	tester := func(subject internftv1alpha1.QueryPropertiesRequest) error {
		res, err := s.queryServer.Properties(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Properties, 2)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryPropertiesRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryPropertiesRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryPropertiesRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id": {
				malleate: func(subject *internftv1alpha1.QueryPropertiesRequest) {
					subject.TokenId = s.tokenIDs[s.customer.String()]
				},
			},
			"invalid token id": {
				malleate: func(subject *internftv1alpha1.QueryPropertiesRequest) {
					subject.TokenId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestQueryOwner() {
	tester := func(subject internftv1alpha1.QueryOwnerRequest) error {
		res, err := s.queryServer.Owner(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotEmpty(res.Owner)

		return nil
	}
	cases := []map[string]Case[internftv1alpha1.QueryOwnerRequest]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.QueryOwnerRequest) {
					subject.ClassId = s.classID
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.QueryOwnerRequest) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id": {
				malleate: func(subject *internftv1alpha1.QueryOwnerRequest) {
					subject.TokenId = s.tokenIDs[s.customer.String()]
				},
			},
			"invalid token id": {
				malleate: func(subject *internftv1alpha1.QueryOwnerRequest) {
					subject.TokenId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(s, tester, cases)
}
