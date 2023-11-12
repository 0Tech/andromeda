package internal_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestNewClass() {
	type newClass struct {
		class *internftv1alpha1.Class
		traits []*internftv1alpha1.Trait
	}

	tester := func(subject newClass) error {
		s.Assert().NoError(subject.class.ValidateBasic())
		s.Assert().NoError(internftv1alpha1.Traits(subject.traits).ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.NewClass(ctx, subject.class, subject.traits)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.class.Id)
		s.Assert().Error(err)
		s.Assert().Nil(classBefore)
		for _, trait := range subject.traits {
			traitBefore, err := s.keeper.GetTrait(s.ctx, subject.class.Id, trait.Id)
			s.Assert().Error(err, trait.Id)
			s.Assert().Nil(traitBefore, trait.Id)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.class.Id)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.class, classAfter)
		for _, trait := range subject.traits {
			traitAfter, err := s.keeper.GetTrait(ctx, subject.class.Id, trait.Id)
			s.Require().NoError(err, trait.Id)
			s.Require().NotNil(traitAfter, trait.Id)
			s.Require().Equal(trait, traitAfter, trait.Id)
		}

		return nil
	}
	cases := []map[string]Case[newClass]{
		{
			"class not found": {
				malleate: func(subject *newClass) {
					subject.class = &internftv1alpha1.Class{
						Id: s.customer.String(),
					}
				},
			},
			"class already exists": {
				malleate: func(subject *newClass) {
					subject.class = &internftv1alpha1.Class{
						Id: s.vendor.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassAlreadyExists
				},
			},
		},
	}
	cases = append(cases, map[string]Case[newClass]{
		"": {
			malleate: func(subject *newClass) {
				subject.traits = []*internftv1alpha1.Trait{}
			},
		},
	})
	for i := 0; i < 4; i++ {
		traitID := fmt.Sprintf("trait%02d", i)

		addedTrait := false
		cases = append(cases, []map[string]Case[newClass]{
			{
				"[nil trait id": {
					malleate: func(subject *newClass) {
						addedTrait = false
					},
				},
				"[valid trait id": {
					malleate: func(subject *newClass) {
						addedTrait = true
						subject.traits = append(subject.traits, &internftv1alpha1.Trait{
							Id: traitID,
						})
					},
				},
			},
			{
				"immutable]": {
					malleate: func(subject *newClass) {
						if !addedTrait {
							return
						}
						subject.traits[len(subject.traits) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_IMMUTABLE
					},
				},
				"mutable]": {
					malleate: func(subject *newClass) {
						if !addedTrait {
							return
						}
						subject.traits[len(subject.traits) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_MUTABLE
					},
				},
			},
		}...)
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestUpdateClass() {
	type updateClass struct {
		class *internftv1alpha1.Class
	}

	tester := func(subject updateClass) error {
		s.Assert().NoError(subject.class.ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.UpdateClass(ctx, subject.class)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.class.Id)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.class, classBefore)

		classAfter, err := s.keeper.GetClass(ctx, subject.class.Id)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.class, classAfter)

		return nil
	}
	cases := []map[string]Case[updateClass]{
		{
			"class exists": {
				malleate: func(subject *updateClass) {
					subject.class = &internftv1alpha1.Class{
						Id: s.vendor.String(),
					}
				},
			},
			"class not found": {
				malleate: func(subject *updateClass) {
					subject.class = &internftv1alpha1.Class{
						Id: s.customer.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassNotFound
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestNewToken() {
	type newToken struct {
		owner sdk.AccAddress
		token *internftv1alpha1.Token
		properties []*internftv1alpha1.Property
	}

	tester := func(subject newToken) error {
		s.Assert().NotEmpty(subject.owner)
		s.Assert().NoError(subject.token.ValidateBasic())
		s.Assert().NoError(internftv1alpha1.Properties(subject.properties).ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.NewToken(ctx, subject.owner, subject.token, subject.properties)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.token.ClassId)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.token.ClassId, classBefore.Id)
		tokenBefore, err := s.keeper.GetToken(s.ctx, subject.token)
		s.Assert().Error(err)
		s.Assert().Nil(tokenBefore)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.token)
		s.Assert().Error(err)
		s.Assert().Nil(ownerBefore)
		for _, property := range subject.properties {
			propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.token, property.TraitId)
			s.Assert().Error(err, property.TraitId)
			s.Assert().Nil(propertyBefore, property.TraitId)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.token.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.token.ClassId, classAfter.Id)
		tokenAfter, err := s.keeper.GetToken(ctx, subject.token)
		s.Require().NoError(err)
		s.Require().NotNil(tokenAfter)
		s.Require().Equal(subject.token, tokenAfter)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.token)
		s.Require().NoError(err)
		s.Require().NotNil(ownerAfter)
		s.Require().Equal(subject.owner, *ownerAfter)
		for _, property := range subject.properties {
			propertyAfter, err := s.keeper.GetProperty(ctx, subject.token, property.TraitId)
			s.Require().NoError(err, property.TraitId)
			s.Require().NotNil(propertyAfter, property.TraitId)
			s.Require().Equal(property, propertyAfter, property.TraitId)
		}

		return nil
	}
	cases := []map[string]Case[newToken]{
		{
			"for operator": {
				malleate: func(subject *newToken) {
					subject.owner = s.vendor
				},
			},
			"not for operator": {
				malleate: func(subject *newToken) {
					subject.owner = s.stranger
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *newToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.vendor.String(),
					}
				},
			},
			"class not found": {
				malleate: func(subject *newToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.stranger.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassNotFound
				},
			},
		},
		{
			"token not found": {
				malleate: func(subject *newToken) {
					subject.token.Id = s.tokenIDs[s.stranger.String()]
				},
			},
			"token already exists": {
				malleate: func(subject *newToken) {
					subject.token.Id = s.tokenIDs[s.customer.String()]
				},
				err: func() error {
					return sdkerrors.ErrInvalidRequest
				},
			},
		},
	}
	cases = append(cases, map[string]Case[newToken]{
		"": {
			malleate: func(subject *newToken) {
				subject.properties = []*internftv1alpha1.Property{}
			},
		},
	})
	for i, traitID := range []string{
		s.immutableTraitID,
		s.mutableTraitID,
	}{
		traitID := traitID
		fact := fmt.Sprintf("fact%02d", i)

		cases = append(cases, map[string]Case[newToken]{
			"not add property": {
			},
			"add property": {
				malleate: func(subject *newToken) {
					subject.properties = append(subject.properties, &internftv1alpha1.Property{
						TraitId: traitID,
						Fact: fact,
					})
				},
			},
		})
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestBurnToken() {
	type burnToken struct {
		owner sdk.AccAddress
		token *internftv1alpha1.Token
	}

	tester := func(subject burnToken) error {
		s.Assert().NotEmpty(subject.owner)
		s.Assert().NoError(subject.token.ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.BurnToken(ctx, subject.owner, subject.token)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.token.ClassId)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.token.ClassId, classBefore.Id)
		tokenBefore, err := s.keeper.GetToken(s.ctx, subject.token)
		s.Assert().NoError(err)
		s.Assert().NotNil(tokenBefore)
		s.Assert().Equal(subject.token, tokenBefore)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.token)
		s.Assert().NoError(err)
		s.Assert().NotNil(ownerBefore)
		s.Assert().Equal(subject.owner, *ownerBefore)
		for _, traitID := range []string{
			s.immutableTraitID,
			s.mutableTraitID,
		}{
			propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.token, traitID)
			s.Assert().NoError(err, traitID)
			s.Assert().NotNil(propertyBefore, traitID)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.token.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.token.ClassId, classAfter.Id)
		tokenAfter, err := s.keeper.GetToken(ctx, subject.token)
		s.Require().Error(err)
		s.Require().Nil(tokenAfter)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.token)
		s.Require().Error(err)
		s.Require().Nil(ownerAfter)
		for _, traitID := range []string{
			s.immutableTraitID,
			s.mutableTraitID,
		}{
			propertyAfter, err := s.keeper.GetProperty(ctx, subject.token, traitID)
			s.Require().Error(err, traitID)
			s.Require().Nil(propertyAfter, traitID)
		}

		return nil
	}
	cases := []map[string]Case[burnToken]{
		{
			"by owner": {
				malleate: func(subject *burnToken) {
					subject.owner = s.customer
				},
			},
			"by operator": {
				malleate: func(subject *burnToken) {
					subject.owner = s.vendor
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientToken
				},
			},
			"by other": {
				malleate: func(subject *burnToken) {
					subject.owner = s.stranger
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientToken
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *burnToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.vendor.String(),
					}
				},
			},
			"class not found": {
				malleate: func(subject *burnToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.stranger.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientToken
				},
			},
		},
		{
			"token exists": {
				malleate: func(subject *burnToken) {
					subject.token.Id = s.tokenIDs[s.customer.String()]
				},
			},
			"token not found": {
				malleate: func(subject *burnToken) {
					subject.token.Id = s.tokenIDs[s.stranger.String()]
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientToken
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestUpdateToken() {
	type updateToken struct {
		token *internftv1alpha1.Token
		properties []*internftv1alpha1.Property
	}

	tester := func(subject updateToken) error {
		s.Assert().NoError(subject.token.ValidateBasic())
		s.Assert().NoError(internftv1alpha1.Properties(subject.properties).ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.UpdateToken(ctx, subject.token, subject.properties)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.token.ClassId)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.token.ClassId, classBefore.Id)
		tokenBefore, err := s.keeper.GetToken(s.ctx, subject.token)
		s.Assert().NoError(err)
		s.Assert().NotNil(tokenBefore)
		s.Assert().Equal(subject.token, tokenBefore)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.token)
		s.Assert().NoError(err)
		s.Assert().NotNil(ownerBefore)
		for _, property := range subject.properties {
			propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.token, property.TraitId)
			s.Assert().NoError(err, property.TraitId)
			s.Assert().NotNil(propertyBefore, property.TraitId)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.token.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.token.ClassId, classAfter.Id)
		tokenAfter, err := s.keeper.GetToken(ctx, subject.token)
		s.Require().NoError(err)
		s.Require().NotNil(tokenAfter)
		s.Require().Equal(subject.token, tokenAfter)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.token)
		s.Require().NoError(err)
		s.Require().NotNil(ownerAfter)
		for _, property := range subject.properties {
			propertyAfter, err := s.keeper.GetProperty(ctx, subject.token, property.TraitId)
			s.Require().NoError(err, property.TraitId)
			s.Require().NotNil(propertyAfter, property.TraitId)
			s.Require().Equal(property, propertyAfter, property.TraitId)
		}

		return nil
	}
	cases := []map[string]Case[updateToken]{
		{
			"class exists": {
				malleate: func(subject *updateToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.vendor.String(),
					}
				},
			},
			"class not found": {
				malleate: func(subject *updateToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.stranger.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrTokenNotFound
				},
			},
		},
		{
			"token exists": {
				malleate: func(subject *updateToken) {
					subject.token.Id = s.tokenIDs[s.customer.String()]
				},
			},
			"token not found": {
				malleate: func(subject *updateToken) {
					subject.token.Id = s.tokenIDs[s.stranger.String()]
				},
				err: func() error {
					return internftv1alpha1.ErrTokenNotFound
				},
			},
		},
	}
	cases = append(cases, map[string]Case[updateToken]{
		"": {
			malleate: func(subject *updateToken) {
				subject.properties = []*internftv1alpha1.Property{}
			},
		},
	})
	for i, traitID := range []string{
		s.immutableTraitID,
		s.mutableTraitID,
	}{
		traitID := traitID
		fact := fmt.Sprintf("newfact%02d", i)

		cases = append(cases, map[string]Case[updateToken]{
			"not add property": {
			},
			"add property": {
				malleate: func(subject *updateToken) {
					subject.properties = append(subject.properties, &internftv1alpha1.Property{
						TraitId: traitID,
						Fact: fact,
					})
				},
				err: func() error {
					if traitID == s.immutableTraitID {
						return internftv1alpha1.ErrTraitImmutable
					}
					return nil
				},
			},
		})
	}

	doTest(s, tester, cases)
}
