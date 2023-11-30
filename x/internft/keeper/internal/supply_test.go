package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestCreateClass() {
	type createClass struct {
		class *internftv1alpha1.Class
	}

	tester := func(subject createClass) error {
		s.Assert().NoError((&internftv1alpha1.ClassInternal{}).Parse(*subject.class))

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.CreateClass(ctx, subject.class)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.class.Id)
		s.Assert().Error(err)
		s.Assert().Nil(classBefore)

		classAfter, err := s.keeper.GetClass(ctx, subject.class.Id)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.class, classAfter)

		return nil
	}
	cases := []map[string]Case[createClass]{
		{
			"class not found": {
				malleate: func(subject *createClass) {
					subject.class = &internftv1alpha1.Class{
						Id: "no-class",
					}
				},
			},
			"class already exists": {
				malleate: func(subject *createClass) {
					subject.class = &internftv1alpha1.Class{
						Id: s.classID,
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassAlreadyExists
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestUpdateTrait() {
	type updateTrait struct {
		class *internftv1alpha1.Class
		trait *internftv1alpha1.Trait
	}

	tester := func(subject updateTrait) error {
		s.Assert().NoError((&internftv1alpha1.ClassInternal{}).Parse(*subject.class))
		s.Assert().NoError((&internftv1alpha1.TraitInternal{}).Parse(*subject.trait))

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.UpdateTrait(ctx, subject.class, subject.trait)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.class.Id)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.class, classBefore)
		traitBefore, err := s.keeper.GetTrait(s.ctx, subject.class.Id, subject.trait.Id)
		if err != nil {
			s.Assert().Nil(traitBefore)
		} else {
			s.Assert().NotNil(traitBefore)
			s.Assert().NotEqual(internftv1alpha1.Trait_MUTABILITY_IMMUTABLE, traitBefore.Mutability)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.class.Id)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.class, classAfter)
		traitAfter, err := s.keeper.GetTrait(ctx, subject.class.Id, subject.trait.Id)
		s.Require().NoError(err)
		s.Require().NotNil(traitAfter)
		s.Require().Equal(subject.trait, traitAfter)

		return nil
	}
	cases := []map[string]Case[updateTrait]{
		{
			"class exists": {
				malleate: func(subject *updateTrait) {
					subject.class = &internftv1alpha1.Class{
						Id: s.classID,
					}
				},
			},
			"class not found": {
				malleate: func(subject *updateTrait) {
					subject.class = &internftv1alpha1.Class{
						Id: "no-class",
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassNotFound
				},
			},
		},
		{
			"trait exists and mutable": {
				malleate: func(subject *updateTrait) {
					subject.trait = &internftv1alpha1.Trait{
						Id: s.mutableTraitID,
					}
				},
			},
			"trait not found": {
				malleate: func(subject *updateTrait) {
					subject.trait = &internftv1alpha1.Trait{
						Id: "new-trait",
					}
				},
			},
			"trait exists and immutable": {
				malleate: func(subject *updateTrait) {
					subject.trait = &internftv1alpha1.Trait{
						Id: s.immutableTraitID,
					}
				},
				err: func() error {
					return internftv1alpha1.ErrTraitImmutable
				},
			},
		},
		{
			"set to immutable": {
				malleate: func(subject *updateTrait) {
					subject.trait.Mutability = internftv1alpha1.Trait_MUTABILITY_IMMUTABLE
				},
			},
			"set to mutable": {
				malleate: func(subject *updateTrait) {
					subject.trait.Mutability = internftv1alpha1.Trait_MUTABILITY_MUTABLE
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestMintToken() {
	type mintToken struct {
		owner sdk.AccAddress
		token *internftv1alpha1.Token
	}

	tester := func(subject mintToken) error {
		s.Assert().NotEmpty(subject.owner)
		s.Assert().NoError((&internftv1alpha1.TokenInternal{}).Parse(*subject.token))

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.MintToken(ctx, subject.owner, subject.token)
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

		return nil
	}
	cases := []map[string]Case[mintToken]{
		{
			"for operator": {
				malleate: func(subject *mintToken) {
					subject.owner = s.vendor
				},
			},
			"not for operator": {
				malleate: func(subject *mintToken) {
					subject.owner = s.stranger
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *mintToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.classID,
					}
				},
			},
			"class not found": {
				malleate: func(subject *mintToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: "no-class",
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassNotFound
				},
			},
		},
		{
			"token not found": {
				malleate: func(subject *mintToken) {
					subject.token.Id = s.tokenIDs[s.stranger.String()]
				},
			},
			"token already exists": {
				malleate: func(subject *mintToken) {
					subject.token.Id = s.tokenIDs[s.customer.String()]
				},
				err: func() error {
					return internftv1alpha1.ErrTokenAlreadyExists
				},
			},
		},
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
		s.Assert().NoError((&internftv1alpha1.TokenInternal{}).Parse(*subject.token))

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
					return internftv1alpha1.ErrPermissionDenied
				},
			},
			"by other": {
				malleate: func(subject *burnToken) {
					subject.owner = s.stranger
				},
				err: func() error {
					return internftv1alpha1.ErrPermissionDenied
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *burnToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.classID,
					}
				},
			},
			"class not found": {
				malleate: func(subject *burnToken) {
					subject.token = &internftv1alpha1.Token{
						ClassId: "no-class",
					}
				},
				err: func() error {
					return internftv1alpha1.ErrPermissionDenied
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
					return internftv1alpha1.ErrPermissionDenied
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestUpdateProperty() {
	type updateProperty struct {
		token *internftv1alpha1.Token
		property *internftv1alpha1.Property
	}

	tester := func(subject updateProperty) error {
		s.Assert().NoError((&internftv1alpha1.TokenInternal{}).Parse(*subject.token))
		s.Assert().NoError((&internftv1alpha1.PropertyInternal{}).Parse(*subject.property))

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.UpdateProperty(ctx, subject.token, subject.property)
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
		propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.token, subject.property.TraitId)
		if err != nil {
			s.Assert().Nil(propertyBefore)
		} else {
			s.Assert().NotNil(propertyBefore)
		}
		traitBefore, err := s.keeper.GetTrait(s.ctx, subject.token.ClassId, subject.property.TraitId)
		s.Assert().NoError(err)
		s.Assert().NotNil(traitBefore)
		s.Assert().NotEqual(internftv1alpha1.Trait_MUTABILITY_IMMUTABLE, traitBefore.Mutability)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.token)
		s.Assert().NoError(err)
		s.Assert().NotNil(ownerBefore)

		classAfter, err := s.keeper.GetClass(ctx, subject.token.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.token.ClassId, classAfter.Id)
		tokenAfter, err := s.keeper.GetToken(ctx, subject.token)
		s.Require().NoError(err)
		s.Require().NotNil(tokenAfter)
		s.Require().Equal(subject.token, tokenAfter)
		propertyAfter, err := s.keeper.GetProperty(ctx, subject.token, subject.property.TraitId)
		s.Require().NoError(err)
		s.Require().NotNil(propertyAfter)
		s.Require().Equal(subject.property, propertyAfter)
		traitAfter, err := s.keeper.GetTrait(s.ctx, subject.token.ClassId, subject.property.TraitId)
		s.Require().NoError(err)
		s.Require().NotNil(traitAfter)
		s.Require().NotEqual(internftv1alpha1.Trait_MUTABILITY_IMMUTABLE, traitBefore.Mutability)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.token)
		s.Require().NoError(err)
		s.Require().NotNil(ownerAfter)

		return nil
	}
	cases := []map[string]Case[updateProperty]{
		{
			"class exists": {
				malleate: func(subject *updateProperty) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.classID,
					}
				},
			},
			"class not found": {
				malleate: func(subject *updateProperty) {
					subject.token = &internftv1alpha1.Token{
						ClassId: "no-class",
					}
				},
				err: func() error {
					return internftv1alpha1.ErrTokenNotFound
				},
			},
		},
		{
			"token exists": {
				malleate: func(subject *updateProperty) {
					subject.token.Id = s.tokenIDs[s.customer.String()]
				},
			},
			"token not found": {
				malleate: func(subject *updateProperty) {
					subject.token.Id = s.tokenIDs[s.stranger.String()]
				},
				err: func() error {
					return internftv1alpha1.ErrTokenNotFound
				},
			},
		},
		{
			"trait exists and mutable": {
				malleate: func(subject *updateProperty) {
					subject.property = &internftv1alpha1.Property{
						TraitId: s.mutableTraitID,
						Fact: "new-fact",
					}
				},
			},
			"trait not found": {
				malleate: func(subject *updateProperty) {
					subject.property = &internftv1alpha1.Property{
						TraitId: "no-trait",
						Fact: "new-fact",
					}
				},
				err: func() error {
					return internftv1alpha1.ErrTraitNotFound
				},
			},
			"trait exists and immutable": {
				malleate: func(subject *updateProperty) {
					subject.property = &internftv1alpha1.Property{
						TraitId: s.immutableTraitID,
						Fact: "new-fact",
					}
				},
				err: func() error {
					return internftv1alpha1.ErrTraitImmutable
				},
			},
		},
	}

	doTest(s, tester, cases)
}
