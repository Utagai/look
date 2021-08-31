package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testEntry struct {
	key interface{}
	val interface{}
}

func entry(k, v interface{}) testEntry {
	return testEntry{
		key: k,
		val: v,
	}
}

type testcase struct {
	name    string
	entries []testEntry
}

// This may be a bit more generally written than necessary...
func TestTableBasicGetAndSet(t *testing.T) {
	tcs := []testcase{
		{
			"string key get",
			[]testEntry{entry("foo", 42)},
		},
		{
			"number key get",
			[]testEntry{entry(42.0, 42)},
		},
		{
			"bool key get",
			[]testEntry{entry(true, 42)},
		},
		{
			"nil key get",
			[]testEntry{entry(nil, 42)},
		},
		{
			"all types",
			[]testEntry{
				entry("foo", 42),
				entry(42.0, 42),
				entry(true, 42),
				entry(nil, 42),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tbl := newTable()
			for _, entry := range tc.entries {
				tbl.Set(entry.key, entry.val)
			}

			for _, entry := range tc.entries {
				key := entry.key
				actualVal, ok := tbl.GetOK(key)
				if !ok {
					t.Fatalf("expected the value %v for key %q to exist, but it did not", entry.val, key)
				} else {
					require.Equal(t, entry.val, actualVal)
				}
			}
		})
	}
}
