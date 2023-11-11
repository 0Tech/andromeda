package internal_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestNewClass() {
	type newClass struct {
		class internftv1alpha1.Class
		traits []internftv1alpha1.Trait
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
		s.Require().Equal(subject.class, *classAfter)
		for _, trait := range subject.traits {
			traitAfter, err := s.keeper.GetTrait(ctx, subject.class.Id, trait.Id)
			s.Require().NoError(err, trait.Id)
			s.Require().NotNil(traitAfter, trait.Id)
			s.Require().Equal(trait, *traitAfter, trait.Id)
		}

		return nil
	}
	cases := []map[string]Case[newClass]{
		{
			"class not found": {
				malleate: func(subject *newClass) {
					subject.class = internftv1alpha1.Class{
						Id: s.customer.String(),
					}
				},
			},
			"class already exists": {
				malleate: func(subject *newClass) {
					subject.class = internftv1alpha1.Class{
						Id: s.vendor.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrClassAlreadyExists
				},
			},
		},
	}
	for i := 0; i < 4; i++ {
		traitID := fmt.Sprintf("trait%02d", i)

		added := false
		cases = append(cases, []map[string]Case[newClass]{
			{
				"no trait": {
					malleate: func(subject *newClass) {
						added = false
					},
				},
				"add trait": {
					malleate: func(subject *newClass) {
						added = true
						subject.traits = append(subject.traits, internftv1alpha1.Trait{
							Id: traitID,
						})
					},
				},
			},
			{
				"immutable": {
				},
				"mutable": {
					malleate: func(subject *newClass) {
						if added {
							subject.traits[len(subject.traits) - 1].Variable = true
						}
					},
				},
			},
		}...)
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestUpdateClass() {
	type updateClass struct {
		class internftv1alpha1.Class
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
		s.Assert().Equal(subject.class, *classBefore)

		classAfter, err := s.keeper.GetClass(ctx, subject.class.Id)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.class, *classAfter)

		return nil
	}
	cases := []map[string]Case[updateClass]{
		{
			"class exists": {
				malleate: func(subject *updateClass) {
					subject.class = internftv1alpha1.Class{
						Id: s.vendor.String(),
					}
				},
			},
			"class not found": {
				malleate: func(subject *updateClass) {
					subject.class = internftv1alpha1.Class{
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

func (s *KeeperTestSuite) TestMintNFT() {
	type mintNFT struct {
		owner sdk.AccAddress
		nft internftv1alpha1.NFT
		properties []internftv1alpha1.Property
	}

	tester := func(subject mintNFT) error {
		s.Assert().NoError(subject.nft.ValidateBasic())
		s.Assert().NoError(internftv1alpha1.Properties(subject.properties).ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.MintNFT(ctx, subject.owner, subject.nft, subject.properties)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.nft.ClassId)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.nft.ClassId, classBefore.Id)
		nftBefore, err := s.keeper.GetNFT(s.ctx, subject.nft)
		s.Assert().Error(err)
		s.Assert().Nil(nftBefore)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.nft)
		s.Assert().Error(err)
		s.Assert().Nil(ownerBefore)
		for _, property := range subject.properties {
			propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.nft, property.Id)
			s.Assert().Error(err, property.Id)
			s.Assert().Nil(propertyBefore, property.Id)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.nft.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.nft.ClassId, classAfter.Id)
		nftAfter, err := s.keeper.GetNFT(ctx, subject.nft)
		s.Require().NoError(err)
		s.Require().NotNil(nftAfter)
		s.Require().Equal(subject.nft, *nftAfter)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.nft)
		s.Require().NoError(err)
		s.Require().NotNil(ownerAfter)
		s.Require().Equal(subject.owner, *ownerAfter)
		for _, property := range subject.properties {
			propertyAfter, err := s.keeper.GetProperty(ctx, subject.nft, property.Id)
			s.Require().NoError(err, property.Id)
			s.Require().NotNil(propertyAfter, property.Id)
			s.Require().Equal(property, *propertyAfter, property.Id)
		}

		return nil
	}
	cases := []map[string]Case[mintNFT]{
		{
			"for operator": {
				malleate: func(subject *mintNFT) {
					subject.owner = s.vendor
				},
			},
			"not for operator": {
				malleate: func(subject *mintNFT) {
					subject.owner = s.stranger
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *mintNFT) {
					subject.nft.ClassId = s.vendor.String()
				},
			},
			"class not found": {
				malleate: func(subject *mintNFT) {
					subject.nft.ClassId = s.stranger.String()
				},
				err: func() error {
					return internftv1alpha1.ErrClassNotFound
				},
			},
		},
		{
			"nft not found": {
				malleate: func(subject *mintNFT) {
					subject.nft.Id = s.nftIDs[s.stranger.String()]
				},
			},
			"nft already exists": {
				malleate: func(subject *mintNFT) {
					subject.nft.Id = s.nftIDs[s.customer.String()]
				},
				err: func() error {
					return sdkerrors.ErrInvalidRequest
				},
			},
		},
	}
	for i, traitID := range []string{
		s.immutableTraitID,
		s.mutableTraitID,
	}{
		traitID := traitID
		fact := fmt.Sprintf("fact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[mintNFT]{
			{
				"no property": {
					malleate: func(subject *mintNFT) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *mintNFT) {
						added = true
						subject.properties = append(subject.properties, internftv1alpha1.Property{
							Id: traitID,
						})
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *mintNFT) {
						if added {
							subject.properties[len(subject.properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestBurnNFT() {
	type burnNFT struct {
		owner sdk.AccAddress
		nft internftv1alpha1.NFT
	}

	tester := func(subject burnNFT) error {
		s.Assert().NoError(subject.nft.ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.BurnNFT(ctx, subject.owner, subject.nft)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.nft.ClassId)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.nft.ClassId, classBefore.Id)
		nftBefore, err := s.keeper.GetNFT(s.ctx, subject.nft)
		s.Assert().NoError(err)
		s.Assert().NotNil(nftBefore)
		s.Assert().Equal(subject.nft, *nftBefore)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.nft)
		s.Assert().NoError(err)
		s.Assert().NotNil(ownerBefore)
		s.Assert().Equal(subject.owner, *ownerBefore)
		for _, traitID := range []string{
			s.immutableTraitID,
			s.mutableTraitID,
		}{
			propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.nft, traitID)
			s.Assert().NoError(err, traitID)
			s.Assert().NotNil(propertyBefore, traitID)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.nft.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.nft.ClassId, classAfter.Id)
		nftAfter, err := s.keeper.GetNFT(ctx, subject.nft)
		s.Require().Error(err)
		s.Require().Nil(nftAfter)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.nft)
		s.Require().Error(err)
		s.Require().Nil(ownerAfter)
		for _, traitID := range []string{
			s.immutableTraitID,
			s.mutableTraitID,
		}{
			propertyAfter, err := s.keeper.GetProperty(ctx, subject.nft, traitID)
			s.Require().Error(err, traitID)
			s.Require().Nil(propertyAfter, traitID)
		}

		return nil
	}
	cases := []map[string]Case[burnNFT]{
		{
			"by owner": {
				malleate: func(subject *burnNFT) {
					subject.owner = s.customer
				},
			},
			"by operator": {
				malleate: func(subject *burnNFT) {
					subject.owner = s.vendor
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientNFT
				},
			},
			"by other": {
				malleate: func(subject *burnNFT) {
					subject.owner = s.stranger
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientNFT
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *burnNFT) {
					subject.nft.ClassId = s.vendor.String()
				},
			},
			"class not found": {
				malleate: func(subject *burnNFT) {
					subject.nft.ClassId = s.stranger.String()
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientNFT
				},
			},
		},
		{
			"nft exists": {
				malleate: func(subject *burnNFT) {
					subject.nft.Id = s.nftIDs[s.customer.String()]
				},
			},
			"nft not found": {
				malleate: func(subject *burnNFT) {
					subject.nft.Id = s.nftIDs[s.stranger.String()]
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientNFT
				},
			},
		},
	}

	doTest(s, tester, cases)
}

func (s *KeeperTestSuite) TestUpdateNFT() {
	type updateNFT struct {
		nft internftv1alpha1.NFT
		properties []internftv1alpha1.Property
	}

	tester := func(subject updateNFT) error {
		s.Assert().NoError(subject.nft.ValidateBasic())
		s.Assert().NoError(internftv1alpha1.Properties(subject.properties).ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.UpdateNFT(ctx, subject.nft, subject.properties)
		if err != nil {
			return err
		}

		classBefore, err := s.keeper.GetClass(s.ctx, subject.nft.ClassId)
		s.Assert().NoError(err)
		s.Assert().NotNil(classBefore)
		s.Assert().Equal(subject.nft.ClassId, classBefore.Id)
		nftBefore, err := s.keeper.GetNFT(s.ctx, subject.nft)
		s.Assert().NoError(err)
		s.Assert().NotNil(nftBefore)
		s.Assert().Equal(subject.nft, *nftBefore)
		ownerBefore, err := s.keeper.GetOwner(s.ctx, subject.nft)
		s.Assert().NoError(err)
		s.Assert().NotNil(ownerBefore)
		for _, property := range subject.properties {
			propertyBefore, err := s.keeper.GetProperty(s.ctx, subject.nft, property.Id)
			s.Assert().NoError(err, property.Id)
			s.Assert().NotNil(propertyBefore, property.Id)
		}

		classAfter, err := s.keeper.GetClass(ctx, subject.nft.ClassId)
		s.Require().NoError(err)
		s.Require().NotNil(classAfter)
		s.Require().Equal(subject.nft.ClassId, classAfter.Id)
		nftAfter, err := s.keeper.GetNFT(ctx, subject.nft)
		s.Require().NoError(err)
		s.Require().NotNil(nftAfter)
		s.Require().Equal(subject.nft, *nftAfter)
		ownerAfter, err := s.keeper.GetOwner(ctx, subject.nft)
		s.Require().NoError(err)
		s.Require().NotNil(ownerAfter)
		for _, property := range subject.properties {
			propertyAfter, err := s.keeper.GetProperty(ctx, subject.nft, property.Id)
			s.Require().NoError(err, property.Id)
			s.Require().NotNil(propertyAfter, property.Id)
			s.Require().Equal(property, *propertyAfter, property.Id)
		}

		return nil
	}
	cases := []map[string]Case[updateNFT]{
		{
			"class exists": {
				malleate: func(subject *updateNFT) {
					subject.nft.ClassId = s.vendor.String()
				},
			},
			"class not found": {
				malleate: func(subject *updateNFT) {
					subject.nft.ClassId = s.stranger.String()
				},
				err: func() error {
					return internftv1alpha1.ErrNFTNotFound
				},
			},
		},
		{
			"nft exists": {
				malleate: func(subject *updateNFT) {
					subject.nft.Id = s.nftIDs[s.customer.String()]
				},
			},
			"nft not found": {
				malleate: func(subject *updateNFT) {
					subject.nft.Id = s.nftIDs[s.stranger.String()]
				},
				err: func() error {
					return internftv1alpha1.ErrNFTNotFound
				},
			},
		},
	}
	for i, traitID := range []string{
		s.immutableTraitID,
		s.mutableTraitID,
	}{
		traitID := traitID
		fact := fmt.Sprintf("newfact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[updateNFT]{
			{
				"no property": {
					malleate: func(subject *updateNFT) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *updateNFT) {
						added = true
						subject.properties = append(subject.properties, internftv1alpha1.Property{
							Id: traitID,
						})
					},
					err: func() error {
						if traitID == s.immutableTraitID {
							return internftv1alpha1.ErrTraitImmutable
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *updateNFT) {
						if added {
							subject.properties[len(subject.properties) - 1].Fact = fact
						}
					},
				},
			},
		}...)
	}

	doTest(s, tester, cases)
}
