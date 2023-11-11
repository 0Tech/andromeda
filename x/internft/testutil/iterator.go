package testutil

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
