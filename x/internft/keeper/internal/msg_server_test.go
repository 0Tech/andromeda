package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestMsgSendToken() {
	testCases := map[string]struct {
		tokenID  string
		err error
	}{
		"valid request": {
			tokenID: s.tokenIDs[s.vendor.String()],
		},
		"not the owner": {
			tokenID:  s.tokenIDs[s.customer.String()],
			err: internftv1alpha1.ErrPermissionDenied,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgSendToken{
				Sender:    s.vendor.String(),
				Recipient: s.customer.String(),
				Token: &internftv1alpha1.Token{
					ClassId: s.classID,
					Id:      tc.tokenID,
				},
			}
			res, err := s.msgServer.SendToken(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgCreateClass() {
	testCases := map[string]struct {
		operator sdk.AccAddress
		err   error
	}{
		"valid request": {
			operator: s.customer,
		},
		// "not the operator": {
		// 	operator: s.customer,
		// 	err:   internftv1alpha1.ErrPermissionDenied,
		// },
		"class already exists": {
			operator: s.vendor,
			err:   internftv1alpha1.ErrClassAlreadyExists,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgCreateClass{
				Operator: tc.operator.String(),
				Class: &internftv1alpha1.Class{
					Id: internftv1alpha1.GetClassID(internftv1alpha1.Address(tc.operator)),
				},
			}
			res, err := s.msgServer.CreateClass(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateTrait() {
	return
	testCases := map[string]struct {
		classID string
		err     error
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		// "not the operator": {
		// 	classID: s.vendor.String(),
		// 	err:     internftv1alpha1.ErrPermissionDenied,
		// },
		"class not found": {
			classID: s.customer.String(),
			err:     internftv1alpha1.ErrClassNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgUpdateTrait{
				Operator: tc.classID,
				Class: &internftv1alpha1.Class{
					Id: tc.classID,
				},
			}
			res, err := s.msgServer.UpdateTrait(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMintToken() {
	return
	newTokenID := createIDs(1, "newtoken")[0]
	testCases := map[string]struct {
		classID string
		err     error
	}{
		"valid request": {
			classID: s.vendor.String(),
		},
		// "not the operator": {
		// 	classID: s.customer.String(),
		// 	err:     internftv1alpha1.ErrPermissionDenied,
		// },
		"class not found": {
			classID: s.customer.String(),
			err:     internftv1alpha1.ErrClassNotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgMintToken{
				Operator: s.vendor.String(),
				Token: &internftv1alpha1.Token{
					ClassId: tc.classID,
					Id: newTokenID,
				},
			}
			res, err := s.msgServer.MintToken(ctx, req)
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
		"not the owner": {
			tokenID: s.tokenIDs[s.customer.String()],
			err: internftv1alpha1.ErrPermissionDenied,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgBurnToken{
				Owner: s.vendor.String(),
				Token: &internftv1alpha1.Token{
					ClassId: s.classID,
					Id:      tc.tokenID,
				},
			}
			res, err := s.msgServer.BurnToken(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgUpdateProperty() {
	return
	testCases := map[string]struct {
		tokenID  string
		err error
	}{
		"valid request": {
			tokenID: s.tokenIDs[s.vendor.String()],
		},
		"not the owner": {
			tokenID: s.tokenIDs[s.stranger.String()],
			err: internftv1alpha1.ErrPermissionDenied,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &internftv1alpha1.MsgUpdateProperty{
				Operator: s.vendor.String(),
				Token: &internftv1alpha1.Token{
					ClassId: s.classID,
					Id:      tc.tokenID,
				},
				Property: &internftv1alpha1.Property{
					TraitId: s.mutableTraitID,
					Fact: "new-fact",
				},
			}
			res, err := s.msgServer.UpdateProperty(ctx, req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}
			s.Require().NotNil(res)
		})
	}
}
