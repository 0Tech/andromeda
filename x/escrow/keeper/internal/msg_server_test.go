package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

func (s *KeeperTestSuite) TestMsgUpdateParams() {
	tester := func(subject escrowv1alpha1.MsgUpdateParams) error {
		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		res, err := s.msgServer.UpdateParams(ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		events := ctx.EventManager().Events()
		s.Require().Len(events, 1)

		eventExpected, err := sdk.TypedEventToEvent(&escrowv1alpha1.EventUpdateParams{
			Authority:         subject.Authority,
			MaxMetadataLength: subject.MaxMetadataLength,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[0])

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.MsgUpdateParams]{
		{
			"nil authority": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid authority": {
				Malleate: func(subject *escrowv1alpha1.MsgUpdateParams) {
					subject.Authority = s.addressBytesToString(s.keeper.GetAuthority())
				},
			},
			"invalid authority": {
				Malleate: func(subject *escrowv1alpha1.MsgUpdateParams) {
					subject.Authority = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil max_metadata_length": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid max_metadata_length": {
				Malleate: func(subject *escrowv1alpha1.MsgUpdateParams) {
					subject.MaxMetadataLength = s.keeper.DefaultGenesis().Params.MaxMetadataLength
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

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
		s.Require().NotEmpty(events)

		eventExpected, err := sdk.TypedEventToEvent(&escrowv1alpha1.EventSubmitProposal{
			Proposal:    s.proposalLast + 1,
			Proposer:    subject.Proposer,
			Agent:       subject.Agent,
			PreActions:  subject.PreActions,
			PostActions: subject.PostActions,
			Metadata:    subject.Metadata,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[len(events)-1])

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
			"proposer not actions signer": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Proposer = s.addressBytesToString(createRandomAccounts(1)[0])
				},
				Error: func() error {
					return escrowv1alpha1.ErrPermissionDenied
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
			"agent not actions signer": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Agent = s.addressBytesToString(createRandomAccounts(1)[0])
				},
				Error: func() error {
					return escrowv1alpha1.ErrPermissionDenied
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
					subject.PreActions = s.encodeMsgs([]sdk.Msg{
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.seller),
							Recipient: s.addressBytesToString(s.agentIdle),
							Asset:     "snake",
						},
					})
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
					subject.PostActions = s.encodeMsgs([]sdk.Msg{
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.agentIdle),
							Recipient: s.addressBytesToString(s.seller),
							Asset:     "voucher",
						},
					})
				},
			},
		},
		{
			"nil metadata": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid metadata": {
				Malleate: func(subject *escrowv1alpha1.MsgSubmitProposal) {
					subject.Metadata = "sell a snake for a voucher"
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
		s.Require().NotEmpty(events)

		eventExpected, err := sdk.TypedEventToEvent(&escrowv1alpha1.EventExec{
			Proposal: subject.Proposal,
			Executor: subject.Executor,
			Actions:  subject.Actions,
		})
		s.Require().NoError(err)
		s.Require().Equal(eventExpected, events[len(events)-1])

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
			"executor not actions signer": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Executor = s.addressBytesToString(createRandomAccounts(1)[0])
				},
				Error: func() error {
					return escrowv1alpha1.ErrPermissionDenied
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
			"agent not actions signer": {
				Malleate: func(subject *escrowv1alpha1.MsgExec) {
					subject.Agent = s.addressBytesToString(createRandomAccounts(1)[0])
				},
				Error: func() error {
					return escrowv1alpha1.ErrPermissionDenied
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
					subject.Actions = s.encodeMsgs([]sdk.Msg{
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.stranger),
							Recipient: s.addressBytesToString(s.agentAny),
							Asset:     "voucher",
						},
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.agentAny),
							Recipient: s.addressBytesToString(s.stranger),
							Asset:     "dog",
						},
					})
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
