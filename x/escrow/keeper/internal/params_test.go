package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowv1alpha1 "github.com/0tech/andromeda/x/escrow/andromeda/escrow/v1alpha1"
	"github.com/0tech/andromeda/x/escrow/testutil"
)

func (s *KeeperTestSuite) TestUpdateParams() {
	type updateParams struct {
		newParams *escrowv1alpha1.Params
	}

	tester := func(subject updateParams) error {
		s.NotNil(subject.newParams)
		s.NotZero(subject.newParams.MaxMetadataLength)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		err := s.keeper.UpdateParams(ctx, subject.newParams)
		if err != nil {
			return err
		}

		paramsBefore, err := s.keeper.GetParams(s.ctx)
		s.Assert().NoError(err)
		s.Assert().NotNil(paramsBefore)

		paramsAfter, err := s.keeper.GetParams(ctx)
		s.Require().NoError(err)
		s.Require().NotNil(paramsAfter)
		s.Require().Equal(subject.newParams, paramsAfter)

		return nil
	}
	cases := []map[string]testutil.Case[updateParams]{
		{
			"": {
				Malleate: func(subject *updateParams) {
					subject.newParams = &escrowv1alpha1.Params{}
				},
			},
		},
		{
			"valid max_metadata_length": {
				Malleate: func(subject *updateParams) {
					subject.newParams.MaxMetadataLength = s.keeper.DefaultGenesis().Params.MaxMetadataLength
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
