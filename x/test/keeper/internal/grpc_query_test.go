package internal_test

import (
	testv1alpha1 "github.com/0tech/andromeda/x/test/andromeda/test/v1alpha1"
	"github.com/0tech/andromeda/x/test/testutil"
)

func (s *KeeperTestSuite) TestQueryAsset() {
	tester := func(subject testv1alpha1.QueryAssetRequest) error {
		res, err := s.queryServer.Asset(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().NotNil(res.Asset)
		s.Require().Equal(subject.Asset, res.Asset.Name)

		return nil
	}
	cases := []map[string]testutil.Case[testv1alpha1.QueryAssetRequest]{
		{
			"nil account": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid account": {
				Malleate: func(subject *testv1alpha1.QueryAssetRequest) {
					subject.Account = s.addressBytesToString(s.catPerson)
				},
			},
			"invalid account": {
				Malleate: func(subject *testv1alpha1.QueryAssetRequest) {
					subject.Account = notInBech32
				},
				Error: func() error {
					return testv1alpha1.ErrInvalidAddress
				},
			},
		},
		{
			"nil asset": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid asset": {
				Malleate: func(subject *testv1alpha1.QueryAssetRequest) {
					subject.Asset = s.cat
				},
			},
			"asset not found": {
				Malleate: func(subject *testv1alpha1.QueryAssetRequest) {
					subject.Asset = s.dog
				},
				Error: func() error {
					return testv1alpha1.ErrAssetNotFound
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}

func (s *KeeperTestSuite) TestQueryAssets() {
	tester := func(subject testv1alpha1.QueryAssetsRequest) error {
		res, err := s.queryServer.Assets(s.ctx, &subject)
		if err != nil {
			return err
		}
		s.Require().NotNil(res)

		s.Require().Len(res.Assets, 2)
		for i, asset := range res.Assets {
			s.Require().NotNil(asset, i)

			s.Require().NotNil(asset.Name, i)
		}

		return nil
	}
	cases := []map[string]testutil.Case[testv1alpha1.QueryAssetsRequest]{
		{
			"nil account": {
				Error: func() error {
					return testv1alpha1.ErrUnimplemented
				},
			},
			"valid account": {
				Malleate: func(subject *testv1alpha1.QueryAssetsRequest) {
					subject.Account = s.addressBytesToString(s.petPerson)
				},
			},
			"invalid account": {
				Malleate: func(subject *testv1alpha1.QueryAssetsRequest) {
					subject.Account = notInBech32
				},
				Error: func() error {
					return testv1alpha1.ErrInvalidAddress
				},
			},
		},
	}

	testutil.DoTest(s.T(), tester, cases)
}
