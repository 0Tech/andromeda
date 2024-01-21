package internal_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func (s *KeeperTestSuite) TestMsgCreateAgent() {
	tester := func(subject escrowv1alpha1.MsgCreateAgent) error {
		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		res, err := s.msgServer.CreateAgent(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&escrowv1alpha1.EventCreateAgent{
			Agent:   "", // we can predict it, but it depends on the state.
			Creator: subject.Creator,
		})
		s.Require().NoError(err)

		for _, compare := range []int{1} {
			s.Require().Equal(eventExpected.Attributes[compare], events[0].Attributes[compare])
		}

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.MsgCreateAgent]{
		{
			"nil creator": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid creator": {
				Malleate: func(subject *escrowv1alpha1.MsgCreateAgent) {
					subject.Creator = s.addressBytesToString(s.stranger)
				},
			},
			"invalid creator": {
				Malleate: func(subject *escrowv1alpha1.MsgCreateAgent) {
					subject.Creator = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestMsgSubmitProposal() {
	tester := func(subject escrowv1alpha1.MsgSubmitProposal) error {
		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		res, err := s.msgServer.SubmitProposal(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&escrowv1alpha1.EventSubmitProposal{
			Id:          s.proposalLast + 1,
			Proposer:    subject.Proposer,
			Agent:       subject.Agent,
			PreActions:  subject.PreActions,
			PostActions: subject.PostActions,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.MsgSubmitProposal]{
		{
			"nil proposer": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid proposer": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Proposer = s.addressBytesToString(s.seller)
				},
			},
			"invalid proposer": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Proposer = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil agent": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid agent": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Agent = s.addressBytesToString(s.agentIdle)
				},
			},
			"invalid agent": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Agent = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil pre_actions": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid pre_actions": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.PreActions = []*codectypes.Any{}
				},
			},
		},
		{
			"nil post_actions": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid post_actions": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.PostActions = []*codectypes.Any{}
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestMsgExec() {
	tester := func(subject escrowv1alpha1.MsgExec) error {
		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		res, err := s.msgServer.Exec(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&escrowv1alpha1.EventExec{
			Proposal: subject.Proposal,
			Executor: subject.Executor,
			Actions:  subject.Actions,
		})
		s.Require().NoError(err)

		for _, compare := range []int{0, 1, 2} {
			s.Require().Equal(eventExpected.Attributes[compare], events[0].Attributes[compare])
		}

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.MsgExec]{
		{
			"nil proposal": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid proposal": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Proposal = s.proposalAny
				},
			},
			"proposal not found": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Proposal = s.proposalLast + 1
				},
				Error: func() error {
					return escrowv1alpha1.ErrProposalNotFound
				},
			},
		},
		{
			"nil executor": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid executor": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Executor = s.addressBytesToString(s.stranger)
				},
			},
			"invalid executor": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Executor = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil agent": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid agent": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Agent = s.addressBytesToString(s.agentAny)
				},
			},
			"invalid agent": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Agent = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil actions": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid actions": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Actions = []*codectypes.Any{}
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
