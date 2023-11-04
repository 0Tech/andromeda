package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestSend() {
	testCases := map[string]struct {
		sender sdk.AccAddress
		nftID     string
		err    error
	}{
		"valid request": {
			nftID: s.nftIDs[s.vendor.String()],
		},
		"insufficient nft": {
			nftID: s.nftIDs[s.customer.String()],
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

			err = s.keeper.Send(ctx, s.vendor, s.customer, nft)
			s.Require().ErrorIs(err, tc.err)
			if err != nil {
				return
			}

			got, err := s.keeper.GetOwner(ctx, nft)
			s.Require().NoError(err)
			s.Require().Equal(s.customer, *got)
		})
	}
}
