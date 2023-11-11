package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestSend() {
	type send struct {
		sender sdk.AccAddress
		recipient sdk.AccAddress
		nft internftv1alpha1.NFT
	}

	tester := func(subject send) error {
		s.Assert().NoError(subject.nft.ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.Send(ctx, subject.sender, subject.recipient, subject.nft)
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
		s.Assert().Equal(subject.sender, *ownerBefore)

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
		s.Require().Equal(subject.recipient, *ownerAfter)

		return nil
	}
	cases := []map[string]Case[send]{
		{
			"sender is owner": {
				malleate: func(subject *send) {
					subject.sender = s.vendor
				},
			},
			"sender is not owner": {
				malleate: func(subject *send) {
					subject.sender = s.stranger
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientNFT
				},
			},
		},
		{
			"recipient is owner": {
				malleate: func(subject *send) {
					subject.recipient = s.vendor
				},
			},
			"recipient is not owner": {
				malleate: func(subject *send) {
					subject.recipient = s.stranger
				},
			},
		},
		{
			"nft exists": {
				malleate: func(subject *send) {
					subject.nft = internftv1alpha1.NFT{
						ClassId: s.vendor.String(),
						Id: s.nftIDs[s.vendor.String()],
					}
				},
			},
			"nft not found": {
				malleate: func(subject *send) {
					subject.nft = internftv1alpha1.NFT{
						ClassId: s.vendor.String(),
						Id: s.nftIDs[s.stranger.String()],
					}
				},
				err: func() error {
					return internftv1alpha1.ErrInsufficientNFT
				},
			},
		},
	}

	doTest(s, tester, cases)
}
