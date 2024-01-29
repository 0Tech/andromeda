package internal_test

import (
	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func (s *KeeperTestSuite) TestQueryParams() {
	tester := func(subject escrowv1alpha1.QueryParamsRequest) error {
		res, err := s.queryServer.Params(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)
		s.Require().NotZero(res.MaxMetadataLength)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryParamsRequest]{
		{
			"valid request": {},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestQueryAgent() {
	tester := func(subject escrowv1alpha1.QueryAgentRequest) error {
		res, err := s.queryServer.Agent(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Agent)
		s.Require().Equal(subject.Agent, res.Agent.Address)
		s.Require().NotNil(res.Agent.Creator)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryAgentRequest]{
		{
			"nil agent": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid agent": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentRequest) {
					subject.Agent = s.addressBytesToString(s.agentIdle)
				},
			},
			"invalid agent": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentRequest) {
					subject.Agent = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
			"agent not found": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentRequest) {
					subject.Agent = s.addressBytesToString(createRandomAccounts(1)[0])
				},
				Error: func() error {
					return escrowv1alpha1.ErrAgentNotFound
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestQueryAgents() {
	tester := func(subject escrowv1alpha1.QueryAgentsRequest) error {
		res, err := s.queryServer.Agents(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Agents, 1)
		for i, agent := range res.Agents {
			s.Require().NotNil(agent, i)

			s.Require().NotNil(agent.Address, i)
			s.Require().NotNil(agent.Creator, i)
		}

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryAgentsRequest]{
		{
			"valid request": {},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestQueryAgentsByCreator() {
	tester := func(subject escrowv1alpha1.QueryAgentsByCreatorRequest) error {
		res, err := s.queryServer.AgentsByCreator(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Agents, 1)
		for i, agent := range res.Agents {
			s.Require().NotNil(agent, i)

			s.Require().NotNil(agent.Address, i)
			s.Require().Equal(subject.Creator, agent.Creator, i)
		}

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryAgentsByCreatorRequest]{
		{
			"nil creator": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid creator": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentsByCreatorRequest) {
					subject.Creator = s.addressBytesToString(s.seller)
				},
			},
			"invalid creator": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentsByCreatorRequest) {
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

func (s *KeeperTestSuite) TestQueryProposal() {
	tester := func(subject escrowv1alpha1.QueryProposalRequest) error {
		res, err := s.queryServer.Proposal(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Proposal)
		s.Require().Equal(subject.Agent, res.Proposal.Agent)
		s.Require().NotNil(res.Proposal.Proposer)
		s.Require().NotNil(res.Proposal.PreActions)
		s.Require().NotNil(res.Proposal.PostActions)
		s.Require().NotEmpty(res.Proposal.Metadata)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryProposalRequest]{
		{
			"nil agent": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid agent": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalRequest) {
					subject.Agent = s.addressBytesToString(s.agentAny)
				},
			},
			"invalid agent": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalRequest) {
					subject.Agent = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
			"proposal not found": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalRequest) {
					subject.Agent = s.addressBytesToString(s.agentIdle)
				},
				Error: func() error {
					return escrowv1alpha1.ErrProposalNotFound
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestQueryProposalsByProposer() {
	tester := func(subject escrowv1alpha1.QueryProposalsByProposerRequest) error {
		res, err := s.queryServer.ProposalsByProposer(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Proposals, 2)
		for i, proposal := range res.Proposals {
			s.Require().NotNil(proposal, i)

			s.Require().NotNil(proposal.Agent, i)
			s.Require().Equal(subject.Proposer, proposal.Proposer, i)
			s.Require().NotNil(proposal.PreActions, i)
			s.Require().NotNil(proposal.PostActions, i)
			s.Require().NotEmpty(proposal.Metadata, i)
		}
		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryProposalsByProposerRequest]{
		{
			"nil proposer": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid proposer": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalsByProposerRequest) {
					subject.Proposer = s.addressBytesToString(s.seller)
				},
			},
			"invalid proposer": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalsByProposerRequest) {
					subject.Proposer = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestQueryProposals() {
	tester := func(subject escrowv1alpha1.QueryProposalsRequest) error {
		res, err := s.queryServer.Proposals(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Proposals, 2)
		for i, proposal := range res.Proposals {
			s.Require().NotNil(proposal, i)

			s.Require().NotNil(proposal.Agent, i)
			s.Require().NotNil(proposal.Proposer, i)
			s.Require().NotNil(proposal.PreActions, i)
			s.Require().NotNil(proposal.PostActions, i)
			s.Require().NotEmpty(proposal.Metadata, i)
		}
		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryProposalsRequest]{
		{
			"valid request": {},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
