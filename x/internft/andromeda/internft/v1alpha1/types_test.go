package internftv1alpha1_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
		ids[i] = internftv1alpha1.GetClassID(internftv1alpha1.Address(addr))
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
		var parsed internftv1alpha1.ClassInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.Class]{
		{
			"nil id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid id": {
				malleate: func(subject *internftv1alpha1.Class) {
					subject.Id = createIDs(1, "class")[0]
				},
			},
			"invalid id": {
				malleate: func(subject *internftv1alpha1.Class) {
					subject.Id = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(t, tester, cases)
}

func TestToken(t *testing.T) {
	tester := func(subject internftv1alpha1.Token) error {
		var parsed internftv1alpha1.TokenInternal
		return parsed.Parse(subject)
	}
	cases := []map[string]Case[internftv1alpha1.Token]{
		{
			"nil class id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid class id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.ClassId = createIDs(1, "class")[0]
				},
			},
			"invalid class id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.ClassId = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
		{
			"nil token id": {
				err: func() error {
					return internftv1alpha1.ErrUnimplemented
				},
			},
			"valid token id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.Id = createIDs(1, "token")[0]
				},
			},
			"invalid token id": {
				malleate: func(subject *internftv1alpha1.Token) {
					subject.Id = "not/in/caip19"
				},
				err: func() error {
					return internftv1alpha1.ErrInvalidID
				},
			},
		},
	}

	doTest(t, tester, cases)
}
