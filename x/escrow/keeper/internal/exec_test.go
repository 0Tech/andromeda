package internal_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func (s *KeeperTestSuite) TestExec() {
	type exec struct {
		id       uint64
		executor sdk.AccAddress
		agent    sdk.AccAddress
		actions  []*codectypes.Any
	}

	tester := func(subject exec) error {
		s.NotZero(subject.id)
		s.NotNil(subject.executor)
		s.NotNil(subject.agent)
		s.NotNil(subject.actions)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		err := s.keeper.Exec(ctx, subject.id, subject.executor, subject.agent, subject.actions)
		if err != nil {
			return err
		}

		_, _, err = s.keeper.GetProposal(s.ctx, subject.id)
		s.Assert().NoError(err)

		_, _, err = s.keeper.GetProposal(ctx, subject.id)
		s.Require().Error(err)

		return nil
	}
	cases := []map[string]testutil.Case[exec]{
		{
			"id already exists": {
				Malleate: func(subject *exec) {
					subject.id = s.proposalAny
				},
			},
			"id not found": {
				Malleate: func(subject *exec) {
					subject.id = s.proposalLast + 1
				},
				Error: func() error {
					return escrowv1alpha1.ErrProposalNotFound
				},
			},
		},
		{
			"executor already exists": {
				Malleate: func(subject *exec) {
					subject.executor = s.stranger
				},
			},
		},
		{
			"agent already exists": {
				Malleate: func(subject *exec) {
					subject.agent = s.agentAny
				},
			},
			"agent differs from the proposal's": {
				Malleate: func(subject *exec) {
					subject.agent = s.agentIdle
				},
				Error: func() error {
					return escrowv1alpha1.ErrPermissionDenied
				},
			},
		},
		{
			"valid actions": {
				Malleate: func(subject *exec) {
					subject.actions = []*codectypes.Any{}
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
