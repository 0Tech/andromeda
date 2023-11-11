package internftv1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func TestGenesisState(t *testing.T) {
	classIDs := createIDs(2, "class")
	tokenIDs := createIDs(2, "token")
	const traitID = "uri"
	addr := createAddresses(1, "addr")[0]

	testCases := map[string]struct {
		s   internftv1alpha1.GenesisState
		err error
	}{
		"default genesis": {
			s: *internftv1alpha1.DefaultGenesisState(),
		},
		"all features": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []internftv1alpha1.Trait{
							{
								Id: traitID,
							},
						},
						Tokens: []internftv1alpha1.GenesisToken{
							{
								Id:    tokenIDs[0],
								Owner: addr.String(),
							},
							{
								Id:    tokenIDs[1],
								Owner: addr.String(),
							},
						},
					},
					{
						Id:              classIDs[1],
						Tokens: []internftv1alpha1.GenesisToken{
							{
								Id:    tokenIDs[0],
								Owner: addr.String(),
							},
							{
								Id:    tokenIDs[1],
								Owner: addr.String(),
							},
						},
					},
				},
			},
		},
		"invalid class id": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
					},
				},
			},
			err: internftv1alpha1.ErrInvalidClassID,
		},
		"invalid trait id": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []internftv1alpha1.Trait{
							{},
						},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidTraitID,
		},
		"duplicate class": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id:              classIDs[0],
					},
					{
						Id:              classIDs[0],
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		"invalid token id": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id:              classIDs[0],
						Tokens: []internftv1alpha1.GenesisToken{
							{
								Owner: addr.String(),
							},
						},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidTokenID,
		},
		"invalid property id": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []internftv1alpha1.Trait{
							{
								Id: traitID,
							},
						},
						Tokens: []internftv1alpha1.GenesisToken{
							{
								Id: tokenIDs[0],
								Properties: []internftv1alpha1.Property{
									{},
								},
								Owner: addr.String(),
							},
						},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidTraitID,
		},
		"no corresponding trait": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []internftv1alpha1.Trait{
							{
								Id: traitID,
							},
						},
						Tokens: []internftv1alpha1.GenesisToken{
							{
								Id: tokenIDs[0],
								Properties: []internftv1alpha1.Property{
									{
										TraitId: "nosuchid",
									},
								},
								Owner: addr.String(),
							},
						},
					},
				},
			},
			err: internftv1alpha1.ErrTraitNotFound,
		},
		"invalid owner": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []internftv1alpha1.GenesisClass{
					{
						Id:              classIDs[0],
						Tokens: []internftv1alpha1.GenesisToken{
							{
								Id:    tokenIDs[0],
								Owner: "invalid",
							},
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.s.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
		})
	}
}
