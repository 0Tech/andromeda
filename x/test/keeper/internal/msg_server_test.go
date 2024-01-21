package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	"github.com/0tech/andromeda/x/test/testutil"
)

func (s *KeeperTestSuite) TestMsgCreate() {
	tester := func(subject testv1alpha1.MsgCreate) error {
		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		res, err := s.msgServer.Create(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&testv1alpha1.EventCreate{
			Creator: subject.Creator,
			Asset:   subject.Asset,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
	}
	cases := []map[string]testutil.Case[testv1alpha1.MsgCreate]{
		{
			"nil creator": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid creator": {
				Malleate: func(subject *testv1alpha1.MsgCreate) {
					subject.Creator = s.addressBytesToString(s.dogPerson)
				},
			},
			"invalid creator": {
				Malleate: func(subject *testv1alpha1.MsgCreate) {
					subject.Creator = notInBech32
				},
				Error: func() error {
					return testv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil asset": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid asset": {
				Malleate: func(subject *testv1alpha1.MsgCreate) {
					subject.Asset = s.cat
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestMsgSend() {
	tester := func(subject testv1alpha1.MsgSend) error {
		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		res, err := s.msgServer.Send(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&testv1alpha1.EventSend{
			Sender:    subject.Sender,
			Recipient: subject.Recipient,
			Asset:     subject.Asset,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
	}
	cases := []map[string]testutil.Case[testv1alpha1.MsgSend]{
		{
			"nil sender": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid sender": {
				Malleate: func(subject *testv1alpha1.MsgSend) {
					subject.Sender = s.addressBytesToString(s.catPerson)
				},
			},
			"invalid sender": {
				Malleate: func(subject *testv1alpha1.MsgSend) {
					subject.Sender = notInBech32
				},
				Error: func() error {
					return testv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil recipient": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid recipient": {
				Malleate: func(subject *testv1alpha1.MsgSend) {
					subject.Recipient = s.addressBytesToString(s.dogPerson)
				},
			},
			"invalid recipient": {
				Malleate: func(subject *testv1alpha1.MsgSend) {
					subject.Recipient = notInBech32
				},
				Error: func() error {
					return testv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil asset": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid asset": {
				Malleate: func(subject *testv1alpha1.MsgSend) {
					subject.Asset = s.cat
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
