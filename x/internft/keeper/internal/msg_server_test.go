package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		nftID  string
		err error
	}{
		"valid request": {
			nftID: s.nftIDs[s.vendor.String()],
		},
		"insufficient funds": {
			nftID:  s.nftIDs[s.customer.String()],
			err: internftv1alpha1.ErrInsufficientNFT,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgSend{
				Sender:    s.vendor.String(),
				Recipient: s.customer.String(),
				Nft: internftv1alpha1.NFT{
					ClassId: s.vendor.String(),
					Id:      tc.nftID,
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.Send(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgNewClass() {
	testCases := map[string]struct {
		operator sdk.AccAddress
		err   error
	}{
		"valid request": {
			operator: s.customer,
		},
		"class already exists": {
			operator: s.vendor,
			err:   internftv1alpha1.ErrClassAlreadyExists,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgNewClass{
				Operator: tc.operator.String(),
				Class: internftv1alpha1.Class{
					Id: tc.operator.String(),
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.NewClass(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateClass() {
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

			req := &internftv1alpha1.MsgUpdateClass{
				Operator: tc.classID,
				Class: internftv1alpha1.Class{
					Id: tc.classID,
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.UpdateClass(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMintNFT() {
	newNFTID := createIDs(1, "newnft")[0]
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

			req := &internftv1alpha1.MsgMintNFT{
				Operator: tc.classID,
				Recipient: s.customer.String(),
				Nft: internftv1alpha1.NFT{
					ClassId: tc.classID,
					Id: newNFTID,
				},
				Properties: []internftv1alpha1.Property{
					{
						Id: s.mutableTraitID,
					},
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.MintNFT(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnNFT() {
	testCases := map[string]struct {
		nftID  string
		err error
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

			req := &internftv1alpha1.MsgBurnNFT{
				Owner: s.vendor.String(),
				Nft: internftv1alpha1.NFT{
					ClassId: s.vendor.String(),
					Id:      tc.nftID,
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.BurnNFT(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateNFT() {
	testCases := map[string]struct {
		nftID  string
		err error
	}{
		"valid request": {
			nftID: s.nftIDs[s.vendor.String()],
		},
		"nft not found": {
			nftID: s.nftIDs[s.stranger.String()],
			err: internftv1alpha1.ErrNFTNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgUpdateNFT{
				Owner: s.vendor.String(),
				Nft: internftv1alpha1.NFT{
					ClassId: s.vendor.String(),
					Id:      tc.nftID,
				},
				Properties: []internftv1alpha1.Property{
					{
						Id: s.mutableTraitID,
					},
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.UpdateNFT(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}
