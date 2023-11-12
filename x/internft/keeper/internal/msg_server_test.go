package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		tokenID  string
		err error
	}{
		"valid request": {
			tokenID: s.tokenIDs[s.vendor.String()],
		},
		"insufficient funds": {
			tokenID:  s.tokenIDs[s.customer.String()],
			err: internftv1alpha1.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgSend{
				Sender:    s.vendor.String(),
				Recipient: s.customer.String(),
				Token: &internftv1alpha1.Token{
					ClassId: s.vendor.String(),
					Id:      tc.tokenID,
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
				Class: &internftv1alpha1.Class{
					Id: tc.operator.String(),
				},
				Traits: []*internftv1alpha1.Trait{},
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
				Class: &internftv1alpha1.Class{
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

func (s *KeeperTestSuite) TestMsgNewToken() {
	newTokenID := createIDs(1, "newtoken")[0]
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

			req := &internftv1alpha1.MsgNewToken{
				Operator: tc.classID,
				Recipient: s.customer.String(),
				Token: &internftv1alpha1.Token{
					ClassId: tc.classID,
					Id: newTokenID,
				},
				Properties: []*internftv1alpha1.Property{
					{
						TraitId: s.mutableTraitID,
						Fact: "fact",
					},
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.NewToken(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurnToken() {
	testCases := map[string]struct {
		tokenID  string
		err error
	}{
		"valid request": {
			tokenID: s.tokenIDs[s.vendor.String()],
		},
		"insufficient token": {
			tokenID: s.tokenIDs[s.customer.String()],
			err: internftv1alpha1.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgBurnToken{
				Owner: s.vendor.String(),
				Token: &internftv1alpha1.Token{
					ClassId: s.vendor.String(),
					Id:      tc.tokenID,
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.BurnToken(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateToken() {
	testCases := map[string]struct {
		tokenID  string
		err error
	}{
		"valid request": {
			tokenID: s.tokenIDs[s.vendor.String()],
		},
		"not the owner": {
			tokenID: s.tokenIDs[s.stranger.String()],
			err: internftv1alpha1.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgUpdateToken{
				Owner: s.vendor.String(),
				Token: &internftv1alpha1.Token{
					ClassId: s.vendor.String(),
					Id:      tc.tokenID,
				},
				Properties: []*internftv1alpha1.Property{
					{
						TraitId: s.mutableTraitID,
						Fact: "newfact",
					},
				},
			}
			err := req.ValidateBasic()
			s.Assert().NoError(err)

			res, err := s.msgServer.UpdateToken(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}
