package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func (s *KeeperTestSuite) TestSend() {
	type send struct {
		sender sdk.AccAddress
		recipient sdk.AccAddress
		token *internftv1alpha1.Token
	}

	tester := func(subject send) error {
		s.Assert().NotEmpty(subject.sender)
		s.Assert().NotEmpty(subject.recipient)
		s.Assert().NoError(subject.token.ValidateBasic())

		ctx, _ := s.ctx.CacheContext()
		err := s.keeper.Send(ctx, subject.sender, subject.recipient, subject.token)
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
		s.Assert().Equal(subject.sender, *ownerBefore)

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
		s.Require().Equal(subject.recipient, *ownerAfter)

		return nil
	}
	cases := []map[string]Case[send]{
		{
			"sender is owner": {
				malleate: func(subject *send) {
					subject.sender = s.customer
				},
			},
			"sender is not owner": {
				malleate: func(subject *send) {
					subject.sender = s.stranger
				},
				err: func() error {
					return internftv1alpha1.ErrPermissionDenied
				},
			},
		},
		{
			"recipient is sender": {
				malleate: func(subject *send) {
					subject.recipient = s.customer
				},
			},
			"recipient is not sender": {
				malleate: func(subject *send) {
					subject.recipient = s.vendor
				},
			},
		},
		{
			"class exists": {
				malleate: func(subject *send) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.vendor.String(),
					}
				},
			},
			"class not found": {
				malleate: func(subject *send) {
					subject.token = &internftv1alpha1.Token{
						ClassId: s.stranger.String(),
					}
				},
				err: func() error {
					return internftv1alpha1.ErrPermissionDenied
				},
			},
		},
		{
			"token exists": {
				malleate: func(subject *send) {
					subject.token.Id = s.tokenIDs[s.customer.String()]
				},
			},
			"token not found": {
				malleate: func(subject *send) {
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
