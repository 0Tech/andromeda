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
				require.Error(t, err)

				for _, candidate := range errs {
					if errors.Is(err, candidate) {
						return
					}
				}
				require.FailNow(t, "unexpected error", err)
			} else {
				require.NoError(t, err)
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
			"valid id": {
				malleate: func(subject *internftv1alpha1.Class) {
					subject.Id = createIDs(1, "class")[0]
				},
			},
			"empty id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestTraits(t *testing.T) {
	tester := func(subject []internftv1alpha1.Trait) error {
		return internftv1alpha1.Traits(subject).ValidateBasic()
	}
	cases := []map[string]Case[[]internftv1alpha1.Trait]{}
	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)

		added := false
		cases = append(cases, []map[string]Case[[]internftv1alpha1.Trait]{
			{
				"no trait": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						added = false
					},
				},
				"add trait": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						added = true
						*subject = append(*subject, internftv1alpha1.Trait{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						if added {
							(*subject)[len(*subject) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if added {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"immutable": {
				},
				"mutable": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						if added {
							(*subject)[len(*subject) - 1].Variable = true
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[[]internftv1alpha1.Trait]{
			{
				"no duplicate trait": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						addedDup = false
					},
				},
				"add duplicate trait": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						addedDup = true
						*subject = append(*subject, internftv1alpha1.Trait{})
					},
					err: func() error {
						if added && addedDup {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						if addedDup {
							(*subject)[len(*subject) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if addedDup {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"immutable": {
				},
				"mutable": {
					malleate: func(subject *[]internftv1alpha1.Trait) {
						if addedDup {
							(*subject)[len(*subject) - 1].Variable = true
						}
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}

func TestNFT(t *testing.T) {
	tester := func(subject internftv1alpha1.NFT) error {
		return subject.ValidateBasic()
	}
	cases := []map[string]Case[internftv1alpha1.NFT]{
		{
			"valid class id": {
				malleate: func(subject *internftv1alpha1.NFT) {
					subject.ClassId = createIDs(1, "class")[0]
				},
			},
			"empty class id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidClassID
				},
			},
		},
		{
			"valid nft id": {
				malleate: func(subject *internftv1alpha1.NFT) {
					subject.Id = createIDs(1, "nft")[0]
				},
			},
			"empty nft id": {
				err: func() error {
					return internftv1alpha1.ErrInvalidNFTID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestProperties(t *testing.T) {
	tester := func(subject []internftv1alpha1.Property) error {
		return internftv1alpha1.Properties(subject).ValidateBasic()
	}
	cases := []map[string]Case[[]internftv1alpha1.Property]{}
	for i := 0; i < 2; i++ {
		traitID := fmt.Sprintf("trait%02d", i)
		fact := fmt.Sprintf("fact%02d", i)

		added := false
		cases = append(cases, []map[string]Case[[]internftv1alpha1.Property]{
			{
				"no property": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						added = false
					},
				},
				"add property": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						added = true
						*subject = append(*subject, internftv1alpha1.Property{})
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						if added {
							(*subject)[len(*subject) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if added {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						if added {
							(*subject)[len(*subject) - 1].Fact = fact
						}
					},
				},
			},
		}...)

		addedDup := false
		cases = append(cases, []map[string]Case[[]internftv1alpha1.Property]{
			{
				"no duplicate property": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						addedDup = false
					},
				},
				"add duplicate property": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						addedDup = true
						*subject = append(*subject, internftv1alpha1.Property{})
					},
					err: func() error {
						if added && addedDup {
							return sdkerrors.ErrInvalidRequest
						}
						return nil
					},
				},
			},
			{
				"of valid id": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						if addedDup {
							(*subject)[len(*subject) - 1].Id = traitID
						}
					},
				},
				"of empty id": {
					err: func() error {
						if addedDup {
							return internftv1alpha1.ErrInvalidTraitID
						}
						return nil
					},
				},
			},
			{
				"with no fact": {
				},
				"with fact": {
					malleate: func(subject *[]internftv1alpha1.Property) {
						if addedDup {
							(*subject)[len(*subject) - 1].Fact = fact
						}
					},
				},
			},
		}...)
	}

	doTest(t, tester, cases)
}
