package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/0tech/andromeda/x/escrow/testutil"
)

func (s *KeeperTestSuite) TestCreateAgent() {
	type createAgent struct {
		creator sdk.AccAddress
	}

	tester := func(subject createAgent) error {
		s.NotNil(subject.creator)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		agent, err := s.keeper.CreateAgent(ctx, subject.creator)
		if err != nil {
			return err
		}
		s.NotNil(agent)

		err = s.keeper.HasAgent(s.ctx, agent, subject.creator)
		s.Assert().Error(err)

		err = s.keeper.HasAgent(ctx, agent, subject.creator)
		s.Require().NoError(err)

		return nil
	}
	cases := []map[string]testutil.Case[createAgent]{
		{
			"valid creator": {
				Malleate: func(subject *createAgent) {
					subject.creator = s.stranger
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
