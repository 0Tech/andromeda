package testutil

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CaseIterator struct {
	valid  bool
	cursor []int
	keys   [][]string
}

func NewCaseIterator[T any](components []map[string]T) CaseIterator {
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
	return CaseIterator{
		valid:  valid,
		cursor: cursor,
		keys:   keys,
	}
}

func (ci CaseIterator) Valid() bool {
	return ci.valid
}

func (ci *CaseIterator) Next() {
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

func (ci CaseIterator) Key() []string {
	res := make([]string, len(ci.keys))
	for i, j := range ci.cursor {
		res[i] = ci.keys[i][j]
	}
	return res
}

type Case[T any] struct {
	Malleate func(*T)
	Error    func() error
}

func DoTest[T any](
	t *testing.T,
	tester func(T) error,
	cases []map[string]Case[T],
) {
	for iter := NewCaseIterator(cases); iter.Valid(); iter.Next() {
		names := iter.Key()

		var subject T
		var errs []error
		for i, name := range names {
			c := cases[i][name]

			if malleate := c.Malleate; malleate != nil {
				malleate(&subject)
			}
			if errGen := c.Error; errGen != nil {
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
				assert.Error(t, err)

				for _, candidate := range errs {
					if errors.Is(err, candidate) {
						return
					}
				}
				assert.FailNow(t, "unexpected error", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
