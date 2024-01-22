package internal_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

func (s *KeeperTestSuite) TestSubmitProposal() {
	type submitProposal struct {
		proposer    sdk.AccAddress
		agent       sdk.AccAddress
		preActions  []*codectypes.Any
		postActions []*codectypes.Any
		metadata    string
	}

	tester := func(subject submitProposal) error {
		s.NotNil(subject.proposer)
		s.NotNil(subject.agent)
		s.NotNil(subject.preActions)
		s.NotNil(subject.postActions)
		s.NotEmpty(subject.metadata)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		id, err := s.keeper.SubmitProposal(ctx, subject.proposer, subject.agent, subject.preActions, subject.postActions, subject.metadata)
		if err != nil {
			return err
		}
		s.NotZero(id)

		proposerBefore, proposalBefore, err := s.keeper.GetProposal(s.ctx, id)
		s.Assert().Error(err)
		s.Assert().Nil(proposerBefore)
		s.Assert().Nil(proposalBefore)

		proposerAfter, proposalAfter, err := s.keeper.GetProposal(ctx, id)
		s.Require().NoError(err)
		s.Require().NotNil(proposerAfter)
		s.Require().NotNil(proposalAfter)
		s.Require().Equal(subject.proposer, proposerAfter)
		s.Require().Equal(subject.agent, sdk.AccAddress(proposalAfter.Agent))
		s.Require().Equal(subject.preActions, proposalAfter.PreActions)
		s.Require().Equal(subject.postActions, proposalAfter.PostActions)
		s.Require().Equal(subject.metadata, proposalAfter.Metadata)

		return nil
	}
	cases := []map[string]testutil.Case[submitProposal]{
		{
			"proposer already exists": {
				Malleate: func(subject *submitProposal) {
					subject.proposer = s.seller
				},
			},
		},
		{
			"agent already exists": {
				Malleate: func(subject *submitProposal) {
					subject.agent = s.agentIdle
				},
			},
			"agent not found": {
				Malleate: func(subject *submitProposal) {
					subject.agent = createRandomAccounts(1)[0]
				},
				Error: func() error {
					return escrowv1alpha1.ErrAgentNotFound
				},
			},
		},
		{
			"valid pre_actions": {
				Malleate: func(subject *submitProposal) {
					subject.preActions = s.encodeMsgs([]sdk.Msg{
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.seller),
							Recipient: s.addressBytesToString(s.agentIdle),
							Asset:     "snake",
						},
					})
				},
			},
			"pre_actions failing": {
				Malleate: func(subject *submitProposal) {
					subject.preActions = s.encodeMsgs([]sdk.Msg{
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.seller),
							Recipient: s.addressBytesToString(s.agentIdle),
							Asset:     "whale", // proposer does not have "whale"
						},
					})
				},
				Error: func() error {
					return testv1alpha1.ErrAssetNotFound
				},
			},
		},
		{
			"valid post_actions": {
				Malleate: func(subject *submitProposal) {
					subject.postActions = s.encodeMsgs([]sdk.Msg{
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
			"valid metadata": {
				Malleate: func(subject *submitProposal) {
					subject.metadata = "sell a snake for a voucher"
				},
			},
			"large metadata": {
				Malleate: func(subject *submitProposal) {
					subject.metadata = string(make([]rune, s.keeper.DefaultGenesis().Params.MaxMetadataLength+1))
				},
				Error: func() error {
					return escrowv1alpha1.ErrLargeMetadata
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
