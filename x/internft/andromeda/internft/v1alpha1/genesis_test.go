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
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []*internftv1alpha1.Trait{
							{
								Id: traitID,
								Mutability: internftv1alpha1.Trait_MUTABILITY_IMMUTABLE,
							},
						},
						Tokens: []*internftv1alpha1.GenesisToken{
							{
								Id:    tokenIDs[0],
								Properties: []*internftv1alpha1.Property{},
								Owner: addr.String(),
							},
							{
								Id:    tokenIDs[1],
								Properties: []*internftv1alpha1.Property{},
								Owner: addr.String(),
							},
						},
					},
					{
						Id:              classIDs[1],
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{
							{
								Id:    tokenIDs[0],
								Properties: []*internftv1alpha1.Property{},
								Owner: addr.String(),
							},
							{
								Id:    tokenIDs[1],
								Properties: []*internftv1alpha1.Property{},
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
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id: "not-in-bech32",
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidClassID,
		},
		"nil trait id": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []*internftv1alpha1.Trait{
							{
								Mutability: internftv1alpha1.Trait_MUTABILITY_IMMUTABLE,
							},
						},
					},
				},
			},
			err: sdkerrors.ErrNotSupported,
		},
		"duplicate class": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id:              classIDs[0],
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{},
					},
					{
						Id:              classIDs[0],
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		"invalid token id": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id:              classIDs[0],
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{
							{
								Id: "not-in-bech32",
								Properties: []*internftv1alpha1.Property{},
								Owner: addr.String(),
							},
						},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidTokenID,
		},
		"nil trait id in property": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []*internftv1alpha1.Trait{
							{
								Id: traitID,
								Mutability: internftv1alpha1.Trait_MUTABILITY_IMMUTABLE,
							},
						},
						Tokens: []*internftv1alpha1.GenesisToken{
							{
								Id: tokenIDs[0],
								Properties: []*internftv1alpha1.Property{
									{
										Fact: "fact",
									},
								},
								Owner: addr.String(),
							},
						},
					},
				},
			},
			err: sdkerrors.ErrNotSupported,
		},
		"no corresponding trait": {
			s: internftv1alpha1.GenesisState{
				Params: internftv1alpha1.DefaultParams(),
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id: classIDs[0],
						Traits: []*internftv1alpha1.Trait{
							{
								Id: traitID,
								Mutability: internftv1alpha1.Trait_MUTABILITY_IMMUTABLE,
							},
						},
						Tokens: []*internftv1alpha1.GenesisToken{
							{
								Id: tokenIDs[0],
								Properties: []*internftv1alpha1.Property{
									{
										TraitId: "no-such-id",
										Fact: "fact",
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
				Classes: []*internftv1alpha1.GenesisClass{
					{
						Id:              classIDs[0],
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{
							{
								Id:    tokenIDs[0],
								Properties: []*internftv1alpha1.Property{},
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
