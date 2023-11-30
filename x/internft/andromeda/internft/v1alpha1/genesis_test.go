package internftv1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
)

func TestTraits(t *testing.T) {
	tester := func(subject []*internftv1alpha1.Trait) error {
		return internftv1alpha1.Traits(subject).ValidateBasic()
	}
	cases := []map[string]Case[[]*internftv1alpha1.Trait]{}
	for _, traitID := range createIDs(2, "trait") {
		traitID := traitID

		added := false
		cases = append(cases, []map[string]Case[[]*internftv1alpha1.Trait]{
			{
				"[nil trait": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						added = false
						*subject = append(*subject, nil)
					},
					err: func() error {
						return sdkerrors.ErrInvalidRequest
					},
				},
				"[non-nil trait": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						added = true
						*subject = append(*subject, &internftv1alpha1.Trait{})
					},
				},
			},
			{
				"nil id": {
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"valid id": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].Id = traitID
					},
				},
				"invalid id": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].Id = "not/in/caip19"
					},
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrInvalidID
					},
				},
			},
			{
				"nil mutability]": {
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"immutable]": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_IMMUTABLE
					},
				},
				"mutable]": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_MUTABLE
					},
				},
			},
		}...)

		addedDuplicate := false
		cases = append(cases, []map[string]Case[[]*internftv1alpha1.Trait]{
			{
				"[no duplicate trait": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						addedDuplicate = false
					},
				},
				"[duplicate trait": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !added {
							return
						}
						addedDuplicate = true
						*subject = append(*subject, &internftv1alpha1.Trait{})
					},
					err: func() error {
						if addedDuplicate {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"nil id": {
					err: func() error {
						if !addedDuplicate {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"valid id": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Id = traitID
					},
				},
				"invalid id": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Id = "not/in/caip19"
					},
					err: func() error {
						if !addedDuplicate {
							return nil
						}
						return internftv1alpha1.ErrInvalidID
					},
				},
			},
			{
				"nil mutability]": {
					err: func() error {
						if !addedDuplicate {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"immutable]": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_IMMUTABLE
					},
				},
				"mutable]": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_MUTABLE
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}

func TestProperties(t *testing.T) {
	tester := func(subject []*internftv1alpha1.Property) error {
		return internftv1alpha1.Properties(subject).ValidateBasic()
	}
	cases := []map[string]Case[[]*internftv1alpha1.Property]{}
	for _, traitID := range createIDs(2, "trait") {
		traitID := traitID

		added := false
		cases = append(cases, []map[string]Case[[]*internftv1alpha1.Property]{
			{
				"[nil property": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						added = false
						*subject = append(*subject, nil)
					},
					err: func() error {
						return sdkerrors.ErrInvalidRequest
					},
				},
				"[non-nil property": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						added = true
						*subject = append(*subject, &internftv1alpha1.Property{})
					},
				},
			},
			{
				"nil trait id": {
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"valid trait id": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].TraitId = traitID
					},
				},
				"invalid trait id": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].TraitId = "not/in/caip19"
					},
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrInvalidID
					},
				},
			},
			{
				"nil fact]": {
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"valid fact]": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].Fact = "valid fact"
					},
				},
			},
		}...)

		addedDuplicate := false
		cases = append(cases, []map[string]Case[[]*internftv1alpha1.Property]{
			{
				"[no duplicate property": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						addedDuplicate = false
					},
				},
				"[duplicate property": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !added {
							return
						}
						addedDuplicate = true
						*subject = append(*subject, &internftv1alpha1.Property{})
					},
					err: func() error {
						if addedDuplicate {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"nil trait id": {
					err: func() error {
						if !addedDuplicate {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"valid trait id": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].TraitId = traitID
					},
				},
				"invalid trait id": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].TraitId = "not/in/caip19"
					},
					err: func() error {
						if !added {
							return nil
						}
						return internftv1alpha1.ErrInvalidID
					},
				},
			},
			{
				"nil fact]": {
					err: func() error {
						if !addedDuplicate {
							return nil
						}
						return internftv1alpha1.ErrUnimplemented
					},
				},
				"valid fact]": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Fact = "valid fact"
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}

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
						Id: "not/in/caip19",
						Traits: []*internftv1alpha1.Trait{},
						Tokens: []*internftv1alpha1.GenesisToken{},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidID,
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
			err: internftv1alpha1.ErrUnimplemented,
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
								Id: "not/in/caip19",
								Properties: []*internftv1alpha1.Property{},
								Owner: addr.String(),
							},
						},
					},
				},
			},
			err: internftv1alpha1.ErrInvalidID,
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
			err: internftv1alpha1.ErrUnimplemented,
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
