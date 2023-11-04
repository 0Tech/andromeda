package internal_test

import (
	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestNewClass() {
	testCases := map[string]struct {
		classID string
		err     error
	}{
		"valid request": {
			classID: s.customer.String(),
		},
		"class already exists": {
			classID: s.vendor.String(),
			err:     internftv1alpha1.ErrClassAlreadyExists,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			class := internftv1alpha1.Class{
				Id: tc.classID,
			}
			err := class.ValidateBasic()
			s.Assert().NoError(err)

			traits := []internftv1alpha1.Trait{
				{
					Id: "uri",
				},
			}

			err = s.keeper.NewClass(ctx, class, traits)
			s.Require().ErrorIs(err, tc.err)
			if err != nil {
				return
			}

			got, err := s.keeper.GetClass(ctx, tc.classID)
			s.Require().NoError(err)
			s.Require().NotNil(got)
			s.Require().Equal(class, *got)
		})
	}
}

func (s *KeeperTestSuite) TestUpdateClass() {
	testCases := map[string]struct {
		classID string
		err     error
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		"class not found": {
			classID: s.customer.String(),
			err:     internftv1alpha1.ErrClassNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			class := internftv1alpha1.Class{
				Id: tc.classID,
			}
			err := class.ValidateBasic()
			s.Assert().NoError(err)

			err = s.keeper.UpdateClass(ctx, class)
			s.Require().ErrorIs(err, tc.err)
			if err != nil {
				return
			}

			got, err := s.keeper.GetClass(ctx, tc.classID)
			s.Require().NoError(err)
			s.Require().NotNil(got)
			s.Require().Equal(class, *got)
		})
	}
}

func (s *KeeperTestSuite) TestMintNFT() {
	testCases := map[string]struct {
		classID    string
		propertyID string
		err        error
	}{
		"valid request": {
			classID:    s.vendor.String(),
			propertyID: s.immutableTraitID,
		},
		"class not found": {
			classID:    s.customer.String(),
			propertyID: s.immutableTraitID,
			err:        internftv1alpha1.ErrClassNotFound,
		},
		"trait not found": {
			classID:    s.vendor.String(),
			propertyID: "no-such-a-trait",
			err:        internftv1alpha1.ErrTraitNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := internftv1alpha1.ValidateClassID(tc.classID)
			s.Assert().NoError(err)

			nft := internftv1alpha1.NFT{
				ClassId: tc.classID,
				Id: s.nftIDs[s.stranger.String()],
			}
			properties := []internftv1alpha1.Property{
				{
					Id:   tc.propertyID,
					Fact: randomString(32),
				},
			}
			err = internftv1alpha1.Properties(properties).ValidateBasic()
			s.Assert().NoError(err)

			err = s.keeper.MintNFT(ctx, s.vendor, nft, properties)
			s.Require().ErrorIs(err, tc.err)
			if err != nil {
				return
			}

			_, err = s.keeper.GetNFT(ctx, nft)
			s.Require().NoError(err)

			got, err := s.keeper.GetProperty(ctx, nft, tc.propertyID)
			s.Require().NoError(err)
			s.Require().NotNil(got)
			s.Require().Equal(properties[0], *got)
		})
	}
}

func (s *KeeperTestSuite) TestBurnNFT() {
	testCases := map[string]struct {
		nftID string
		err error
	}{
		"valid request": {
			nftID: s.nftIDs[s.vendor.String()],
		},
		"insufficient nft": {
			nftID:  s.nftIDs[s.stranger.String()],
			err: internftv1alpha1.ErrInsufficientNFT,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			nft := internftv1alpha1.NFT{
				ClassId: s.vendor.String(),
				Id:      tc.nftID,
			}
			err := nft.ValidateBasic()
			s.Assert().NoError(err)

			err = s.keeper.BurnNFT(ctx, s.vendor, nft)
			s.Require().ErrorIs(err, tc.err)
			if err != nil {
				return
			}

			got, err := s.keeper.GetNFT(ctx, nft)
			s.Require().Error(err)
			s.Require().Nil(got)
		})
	}
}

func (s *KeeperTestSuite) TestUpdateNFT() {
	testCases := map[string]struct {
		nftID string
		propertyID string
		err        error
	}{
		"valid request": {
			nftID: s.nftIDs[s.vendor.String()],
			propertyID: s.mutableTraitID,
		},
		"nft not found": {
			nftID: s.nftIDs[s.stranger.String()],
			propertyID: s.mutableTraitID,
			err:        internftv1alpha1.ErrNFTNotFound,
		},
		"trait not found": {
			nftID: s.nftIDs[s.vendor.String()],
			propertyID: "no-such-a-trait",
			err:        internftv1alpha1.ErrTraitNotFound,
		},
		"trait immutable": {
			nftID: s.nftIDs[s.vendor.String()],
			propertyID: s.immutableTraitID,
			err:        internftv1alpha1.ErrTraitImmutable,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			nft := internftv1alpha1.NFT{
				ClassId: s.vendor.String(),
				Id:      tc.nftID,
			}
			err := nft.ValidateBasic()
			s.Assert().NoError(err)

			property := internftv1alpha1.Property{
				Id:   tc.propertyID,
				Fact: randomString(32),
			}
			err = property.ValidateBasic()
			s.Assert().NoError(err)

			err = s.keeper.UpdateNFT(ctx, nft, []internftv1alpha1.Property{property})
			s.Require().ErrorIs(err, tc.err)
			if err != nil {
				return
			}

			got, err := s.keeper.GetProperty(ctx, nft, tc.propertyID)
			s.Require().NoError(err)
			s.Require().NotNil(got)
			s.Require().Equal(property, *got)
		})
	}
}
