package internal_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func (s *KeeperTestSuite) TestSubmitProposal() {
	type submitProposal struct {
		proposer    sdk.AccAddress
		agent       sdk.AccAddress
		preActions  []*codectypes.Any
		postActions []*codectypes.Any
	}

	tester := func(subject submitProposal) error {
		s.NotNil(subject.proposer)
		s.NotNil(subject.agent)
		s.NotNil(subject.preActions)
		s.NotNil(subject.postActions)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		id, err := s.keeper.SubmitProposal(ctx, subject.proposer, subject.agent, subject.preActions, subject.postActions)
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
					subject.agent = simtestutil.CreateRandomAccounts(1)[0]
				},
				Error: func() error {
					return escrowv1alpha1.ErrAgentNotFound
				},
			},
		},
		{
			"valid pre_actions": {
				Malleate: func(subject *submitProposal) {
					subject.preActions = []*codectypes.Any{}
				},
			},
		},
		{
			"valid post_actions": {
				Malleate: func(subject *submitProposal) {
					subject.postActions = []*codectypes.Any{}
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
