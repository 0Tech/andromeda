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

type caseIterator struct {
	valid bool
	cursor []int
	keys [][]string
}

func NewCaseIterator[T any](components []map[string]T) caseIterator {
	keys := make([][]string, len(components))
	valid := true
	for i, component := range components {
		if len(component) == 0 {
			valid = false
			break
		}

		for key := range component {
			keys[i] = append(keys[i], key)
		}
	}

	cursor := make([]int, len(keys))
	return caseIterator{
		valid: valid,
		cursor: cursor,
		keys: keys,
	}
}

func (ci caseIterator) Valid() bool {
	return ci.valid
}

func (ci *caseIterator) Next() {
	for i := len(ci.keys) - 1; i >= 0; i-- {
		next := ci.cursor[i] + 1
		if next < len(ci.keys[i]) {
			ci.cursor[i] = next
			return
		}

		ci.cursor[i] = 0
	}
	ci.valid = false
}

func (ci caseIterator) Key() []string {
	res := make([]string, len(ci.keys))
	for i, j := range ci.cursor {
		res[i] = ci.keys[i][j]
	}
	return res
}

func doTest[T any](
	t *testing.T,
	modifiers []map[string]func(*T) error,
	tester func(T) error) {
	for iter := NewCaseIterator(modifiers); iter.Valid(); iter.Next() {
		names := iter.Key()

		var subject T
		var errs []error
		for i, name := range names {
			modifier := modifiers[i][name]
			if err := modifier(&subject); err != nil {
				errs = append(errs, err)
			}
		}

		t.Run(strings.Join(names, ","), func(t *testing.T) {
			err := tester(subject)
			if len(errs) != 0 {
				for _, candidate := range errs {
					if errors.Is(err, candidate) {
						return
					}
				}
				require.Fail(t, err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClass(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.Class) error{
		{
			"valid_id": func(subject *internftv1alpha1.Class) error {
				subject.Id = createIDs(1, "class")[0]
				return nil
			},
			"empty_id": func(subject *internftv1alpha1.Class) error {
				return internftv1alpha1.ErrInvalidClassID
			},
		},
	}
	tester := func(subject internftv1alpha1.Class) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestTraits(t *testing.T) {
	modifiers := []map[string]func(*[]internftv1alpha1.Trait) error{
		{
			"no_trait": func(subject *[]internftv1alpha1.Trait) error {
				return nil
			},
			"duplicate_trait": func(subject *[]internftv1alpha1.Trait) error {
				*subject = []internftv1alpha1.Trait{
					{Id: "color"},
					{Id: "color", Variable: true},
				}
				return sdkerrors.ErrInvalidRequest
			},
			"immutable_trait": func(subject *[]internftv1alpha1.Trait) error {
				*subject = []internftv1alpha1.Trait{
					{Id: "color"},
				}
				return nil
			},
			"mutable_trait": func(subject *[]internftv1alpha1.Trait) error {
				*subject = []internftv1alpha1.Trait{
					{Id: "color", Variable: true},
				}
				return nil
			},
			"empty_id_immutable_trait": func(subject *[]internftv1alpha1.Trait) error {
				*subject = []internftv1alpha1.Trait{{}}
				return internftv1alpha1.ErrInvalidTraitID
			},
			"empty_id_mutable_trait": func(subject *[]internftv1alpha1.Trait) error {
				*subject = []internftv1alpha1.Trait{
					{Variable: true},
				}
				return internftv1alpha1.ErrInvalidTraitID
			},
		},
	}
	tester := func(subject []internftv1alpha1.Trait) error {
		return internftv1alpha1.Traits(subject).ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestNFT(t *testing.T) {
	modifiers := []map[string]func(*internftv1alpha1.NFT) error{
		{
			"valid_class_id": func(subject *internftv1alpha1.NFT) error {
				subject.ClassId = createIDs(1, "class")[0]
				return nil
			},
			"empty_class_id": func(subject *internftv1alpha1.NFT) error {
				return internftv1alpha1.ErrInvalidClassID
			},
		},
		{
			"valid_nft_id": func(subject *internftv1alpha1.NFT) error {
				subject.Id = createIDs(1, "nft")[0]
				return nil
			},
			"empty_nft_id": func(subject *internftv1alpha1.NFT) error {
				return internftv1alpha1.ErrInvalidNFTID
			},
		},
	}
	tester := func(subject internftv1alpha1.NFT) error {
		return subject.ValidateBasic()
	}

	doTest(t, modifiers, tester)
}

func TestProperties(t *testing.T) {
	modifiers := []map[string]func(*[]internftv1alpha1.Property) error{
		{
			"no_property": func(subject *[]internftv1alpha1.Property) error {
				return nil
			},
			"duplicate_property": func(subject *[]internftv1alpha1.Property) error {
				*subject = []internftv1alpha1.Property{
					{Id: "color"},
					{Id: "color", Fact: "black"},
				}
				return sdkerrors.ErrInvalidRequest
			},
			"empty_id_empty_fact": func(subject *[]internftv1alpha1.Property) error {
				*subject = []internftv1alpha1.Property{{}}
				return internftv1alpha1.ErrInvalidTraitID
			},
			"empty_id_nonempty_fact": func(subject *[]internftv1alpha1.Property) error {
				*subject = []internftv1alpha1.Property{
					{Fact: "black"},
				}
				return internftv1alpha1.ErrInvalidTraitID
			},
		},
	}
	tester := func(subject []internftv1alpha1.Property) error {
		return internftv1alpha1.Properties(subject).ValidateBasic()
	}

	doTest(t, modifiers, tester)
}
