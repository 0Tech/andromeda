package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	"github.com/0tech/andromeda/x/test/testutil"
)

func (s *KeeperTestSuite) TestCreate() {
	type create struct {
		creator sdk.AccAddress
		asset   string
	}

	tester := func(subject create) error {
		s.NotNil(subject.creator)
		s.NotEmpty(subject.asset)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		if err := s.keeper.Create(ctx, subject.creator, subject.asset); err != nil {
			return err
		}

		err := s.keeper.HasAsset(s.ctx, subject.creator, subject.asset)
		s.Assert().Error(err)

		err = s.keeper.HasAsset(ctx, subject.creator, subject.asset)
		s.Require().NoError(err)

		return nil
	}
	cases := []map[string]testutil.Case[create]{
		{
			"valid creator": {
				Malleate: func(subject *create) {
					subject.creator = s.catPerson
				},
			},
		},
		{
			"asset not found": {
				Malleate: func(subject *create) {
					subject.asset = s.dog
				},
			},
			"asset already exists": {
				Malleate: func(subject *create) {
					subject.asset = s.cat
				},
				Error: func() error {
					return testv1alpha1.ErrAssetAlreadyExists
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestSend() {
	type send struct {
		sender    sdk.AccAddress
		recipient sdk.AccAddress
		asset     string
	}

	tester := func(subject send) error {
		s.NotNil(subject.sender)
		s.NotNil(subject.recipient)
		s.NotEmpty(subject.asset)

		ctx, _ := sdk.UnwrapSDKContext(s.ctx).CacheContext()
		if err := s.keeper.Send(ctx, subject.sender, subject.recipient, subject.asset); err != nil {
			return err
		}

		err := s.keeper.HasAsset(s.ctx, subject.sender, subject.asset)
		s.Assert().NoError(err)
		err = s.keeper.HasAsset(s.ctx, subject.recipient, subject.asset)
		s.Assert().Error(err)

		err = s.keeper.HasAsset(ctx, subject.sender, subject.asset)
		s.Require().Error(err)
		err = s.keeper.HasAsset(ctx, subject.recipient, subject.asset)
		s.Require().NoError(err)

		return nil
	}
	cases := []map[string]testutil.Case[send]{
		{
			"valid send": {
				Malleate: func(subject *send) {
					subject.sender = s.catPerson
					subject.recipient = s.dogPerson
					subject.asset = s.cat
				},
			},
			"asset already exists on recipient": {
				Malleate: func(subject *send) {
					subject.sender = s.catPerson
					subject.recipient = s.petPerson
					subject.asset = s.cat
				},
				Error: func() error {
					return testv1alpha1.ErrAssetAlreadyExists
				},
			},
			"asset not found on sender": {
				Malleate: func(subject *send) {
					subject.sender = s.catPerson
					subject.recipient = s.notPetPerson
					subject.asset = s.dog
				},
				Error: func() error {
					return testv1alpha1.ErrAssetNotFound
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
