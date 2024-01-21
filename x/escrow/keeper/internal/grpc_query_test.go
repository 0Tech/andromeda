package internal_test

import (
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

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

		s.Require().NotNil(res.Params)

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
		s.Require().Equal(subject.Address, res.Agent.Address)
		s.Require().NotNil(res.Agent.Creator)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryAgentRequest]{
		{
			"nil address": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid address": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentRequest) {
					subject.Address = s.addressBytesToString(s.agentIdle)
				},
			},
			"invalid address": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentRequest) {
					subject.Address = notInBech32
				},
				Error: func() error {
					return escrowv1alpha1.ErrInvalidAddress
				},
			},
			"address not found": {
				Malleate: func(subject *escrowv1alpha1.QueryAgentRequest) {
					subject.Address = s.addressBytesToString(simtestutil.CreateRandomAccounts(1)[0])
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

func (s *KeeperTestSuite) TestQueryProposal() {
	tester := func(subject escrowv1alpha1.QueryProposalRequest) error {
		res, err := s.queryServer.Proposal(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Proposal)
		s.Require().Equal(subject.Id, res.Proposal.Id)
		s.Require().NotNil(res.Proposal.Proposer)
		s.Require().NotNil(res.Proposal.Agent)
		s.Require().NotNil(res.Proposal.PreActions)
		s.Require().NotNil(res.Proposal.PostActions)

		return nil
	}
	cases := []map[string]testutil.Case[escrowv1alpha1.QueryProposalRequest]{
		{
			"nil id": {
				Error: func() error {
					return escrowv1alpha1.ErrUnimplemented
				},
			},
			"valid id": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalRequest) {
					subject.Id = s.proposalLast
				},
			},
			"id not found": {
				Malleate: func(subject *escrowv1alpha1.QueryProposalRequest) {
					subject.Id = s.proposalLast + 1
				},
				Error: func() error {
					return escrowv1alpha1.ErrProposalNotFound
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

			s.Require().NotZero(proposal.Id, i)
			s.Require().NotNil(proposal.Proposer, i)
			s.Require().NotNil(proposal.Agent, i)
			s.Require().NotNil(proposal.PreActions, i)
			s.Require().NotNil(proposal.PostActions, i)
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