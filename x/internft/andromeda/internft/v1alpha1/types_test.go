package internftv1alpha1_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	internftv1alpha1 "github.com/0tech/andromeda/x/internft/andromeda/internft/v1alpha1"
	"github.com/0tech/andromeda/x/internft/testutil"
)

func createAddresses(size int, prefix string) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, size)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(fmt.Sprintf("%s%d", prefix, i))
	}

	return addrs
}

func createIDs(size int, prefix string) []string {
	addrs := createAddresses(size, prefix)
	ids := make([]string, len(addrs))
	for i, addr := range addrs {
		ids[i] = addr.String()
	}

	return ids
}

type Case[T any] struct {
	malleate func(*T)
	err func() error
}

func doTest[T any](
	t *testing.T,
	tester func(T) error,
	cases []map[string]Case[T]) {
	for iter := testutil.NewCaseIterator(cases); iter.Valid(); iter.Next() {
		names := iter.Key()

		var subject T
		var errs []error
		for i, name := range names {
			c := cases[i][name]

			if malleate := c.malleate; malleate != nil {
				malleate(&subject)
			}
			if errGen := c.err; errGen != nil {
				if err := errGen(); err != nil {
					errs = append(errs, err)
				}
			}
		}

		testName := func(names []string) string {
			display := make([]string, 0, len(names))
			for _, name := range names {
				if len(name) != 0 {
					display = append(display, name)
				}
			}
			return strings.Join(display, ",")
		}
		t.Run(testName(names), func(t *testing.T) {
			err := tester(subject)
			if len(errs) != 0 {
				require.Error(t, err, subject)

				for _, candidate := range errs {
					if errors.Is(err, candidate) {
						return
					}
				}

				errStrings := make([]string, len(errs))
				for i, e := range errs {
					errStrings[i] = e.Error()
				}
				require.FailNow(t, fmt.Sprintf(`Received unexpected error:
%s
Expected errors are:
%s`, err.Error(), strings.Join(errStrings, "\n")), subject)
			} else {
				require.NoError(t, err, subject)
			}
		})
	}
}

func TestClass(t *testing.T) {
	tester := func(subject internftv1alpha1.Class) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.Class]{
		{
			"nil id": {
				err: func() error {
					return sdkerrors.ErrNotSupported
				},
			},
			"valid id": {
				malleate: func(subject *internftv1alpha1.Class) {
					subject.Id = createIDs(1, "class")[0]
				},
			},
			"invalid id": {
				malleate: func(subject *internftv1alpha1.Class) {
					subject.Id = "not-in-bech32"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestTraits(t *testing.T) {
	tester := func(subject []*internftv1alpha1.Trait) error {
		return internftv1alpha1.Traits(subject).ValidateBasic()
	}
	cases := []map[string]Case[[]*internftv1alpha1.Trait]{}
	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)

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
						return sdkerrors.ErrNotSupported
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
			},
			{
				"nil mutability]": {
					err: func() error {
						if !added {
							return nil
						}
						return sdkerrors.ErrNotSupported
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

		AddedDuplicate := false
		cases = append(cases, []map[string]Case[[]*internftv1alpha1.Trait]{
			{
				"[no duplicate trait": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						AddedDuplicate = false
					},
				},
				"[duplicate trait": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !added {
							return
						}
						AddedDuplicate = true
						*subject = append(*subject, &internftv1alpha1.Trait{})
					},
					err: func() error {
						if AddedDuplicate {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"nil id": {
					err: func() error {
						if !AddedDuplicate {
							return nil
						}
						return sdkerrors.ErrNotSupported
					},
				},
				"valid id": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !AddedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Id = traitID
					},
				},
			},
			{
				"nil mutability]": {
					err: func() error {
						if !AddedDuplicate {
							return nil
						}
						return sdkerrors.ErrNotSupported
					},
				},
				"immutable]": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !AddedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Mutability = internftv1alpha1.Trait_MUTABILITY_IMMUTABLE
					},
				},
				"mutable]": {
					malleate: func(subject *[]*internftv1alpha1.Trait) {
						if !AddedDuplicate {
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

func TestToken(t *testing.T) {
	tester := func(subject internftv1alpha1.Token) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.Token]{
		{
			"nil class id": {
				err: func() error {
					return sdkerrors.ErrNotSupported
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.ClassId = createIDs(1, "class")[0]
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.ClassId = "not-in-bech32"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"nil token id": {
				err: func() error {
					return sdkerrors.ErrNotSupported
				},
			},
			"valid token id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.Id = createIDs(1, "token")[0]
				},
			},
			"invalid token id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.Id = "not-in-bech32"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidTokenID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestProperties(t *testing.T) {
	tester := func(subject []*internftv1alpha1.Property) error {
		return internftv1alpha1.Properties(subject).ValidateBasic()
	}
	cases := []map[string]Case[[]*internftv1alpha1.Property]{}
	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)
		fact := fmt.Sprintf("fact%02d", i)

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
						return sdkerrors.ErrNotSupported
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
			},
			{
				"nil fact]": {
					err: func() error {
						if !added {
							return nil
						}
						return sdkerrors.ErrNotSupported
					},
				},
				"valid fact]": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !added {
							return
						}
						(*subject)[len(*subject) - 1].Fact = fact
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
						return sdkerrors.ErrNotSupported
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
			},
			{
				"nil fact]": {
					err: func() error {
						if !addedDuplicate {
							return nil
						}
						return sdkerrors.ErrNotSupported
					},
				},
				"valid fact]": {
					malleate: func(subject *[]*internftv1alpha1.Property) {
						if !addedDuplicate {
							return
						}
						(*subject)[len(*subject) - 1].Fact = fact
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}
