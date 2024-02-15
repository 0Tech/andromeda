package internal_test

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
)

func (s *KeeperTestSuite) TestExec() {
	type exec struct {
		agents  []sdk.AccAddress
		actions []*codectypes.Any
	}

	tester := func(subject exec) error {
		s.NotNil(subject.agents)
		s.NotNil(subject.actions)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		err := s.keeper.Exec(ctx, subject.agents, subject.actions)
		if err != nil {
			return err
		}

		for i, agent := range subject.agents {
			_, err = s.keeper.GetProposal(s.ctx, agent)
			s.Assert().NoError(err, i)

			_, err = s.keeper.GetProposal(ctx, agent)
			s.Require().Error(err, i)
		}

		return nil
	}
	cases := []map[string]testutil.Case[exec]{
		{
			"proposals already exist": {
				Malleate: func(subject *exec) {
					subject.agents = []sdk.AccAddress{s.agentAny}
				},
			},
			"proposal not found": {
				Malleate: func(subject *exec) {
					subject.agents = []sdk.AccAddress{s.agentIdle}
				},
				Error: func() error {
					return escrowv1alpha1.ErrProposalNotFound
				},
			},
		},
		{
			"valid actions": {
				Malleate: func(subject *exec) {
					subject.actions = s.encodeMsgs([]sdk.Msg{
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
			"actions failing": {
				Malleate: func(subject *exec) {
					subject.actions = s.encodeMsgs([]sdk.Msg{
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.stranger),
							Recipient: s.addressBytesToString(s.agentAny),
							Asset:     "voucher",
						},
						&testv1alpha1.MsgSend{
							Sender:    s.addressBytesToString(s.agentAny),
							Recipient: s.addressBytesToString(s.stranger),
							Asset:     "whale", // agent has "dog"
						},
					})
				},
				Error: func() error {
					return testv1alpha1.ErrAssetNotFound
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
