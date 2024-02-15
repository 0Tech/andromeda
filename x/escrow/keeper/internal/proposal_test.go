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
		err := s.keeper.SubmitProposal(ctx, subject.proposer, subject.agent, subject.preActions, subject.postActions, subject.metadata)
		if err != nil {
			return err
		}

		_, err = s.keeper.GetProposal(s.ctx, subject.agent)
		s.Assert().Error(err)

		proposalAfter, err := s.keeper.GetProposal(ctx, subject.agent)
		s.Require().NoError(err)
		s.Require().NotNil(proposalAfter)
		s.Require().Equal(subject.proposer, sdk.AccAddress(proposalAfter.Proposer))
		s.Require().Equal(subject.preActions, proposalAfter.PreActions)
		s.Require().Equal(subject.postActions, proposalAfter.PostActions)
		s.Require().Equal(subject.metadata, proposalAfter.Metadata)

		return nil
	}

	proposer := s.seller
	agent := s.agentIdle
	agentInfo, err := s.keeper.GetAgent(s.ctx, agent)
	s.NoError(err)
	s.Equal(proposer, sdk.AccAddress(agentInfo.Creator))

	validAsset := "snake"
	_, err = s.testQueryServer.Asset(s.ctx, &testv1alpha1.QueryAssetRequest{
		Account: s.addressBytesToString(proposer),
		Asset:   validAsset,
	})
	s.NoError(err)

	invalidAsset := "whale"
	_, err = s.testQueryServer.Asset(s.ctx, &testv1alpha1.QueryAssetRequest{
		Account: s.addressBytesToString(proposer),
		Asset:   invalidAsset,
	})
	s.Error(err)

	param, err := s.keeper.GetParams(s.ctx)
	s.NoError(err)

	cases := []map[string]testutil.Case[submitProposal]{
		{
			"proposer already exists": {
				Malleate: func(subject *submitProposal) {
					subject.proposer = proposer
				},
			},
			"proposer differs from creator": {
				Malleate: func(subject *submitProposal) {
					subject.proposer = s.stranger
				},
				Error: func() error {
					return escrowv1alpha1.ErrPermissionDenied
				},
			},
		},
		{
			"agent already exists": {
				Malleate: func(subject *submitProposal) {
					subject.agent = agent
				},
			},
			"agent not found": {
				Malleate: func(subject *submitProposal) {
					subject.agent = createRandomAddress()
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
							Asset:     validAsset,
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
							Asset:     invalidAsset,
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
					subject.metadata = string(make([]rune, param.MaxMetadataLength+1))
				},
				Error: func() error {
					return escrowv1alpha1.ErrLargeMetadata
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
