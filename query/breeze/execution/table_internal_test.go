package execution

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testEntry struct {
	key    interface{}
	val    interface{}
	exists bool
}

func val(k, v interface{}) testEntry {
	return testEntry{
		key:    k,
		val:    v,
		exists: true,
	}
}

func missing() testEntry {
	return testEntry{
		key:    nil,
		val:    nil,
		exists: false,
	}
}

type testcase struct {
	name    string
	entries []testEntry
}

// We only need one test function here despite there being 3 prominent methods
// on the Table struct. This is largely because in order to test Set, we must
// use GetOK and it must work correctly. Get is implemented entirely by GetOK,
// although that is a bit of the implementation detail leaking into these
// tests, which isn't perfect.
func TestTable(t *testing.T) {
	tcs := []testcase{
		{
			"string key get",
			[]testEntry{val("foo", 42)},
		},
		{
			"number key get",
			[]testEntry{val(42.0, 42)},
		},
		{
			"bool key get",
			[]testEntry{val(true, 42)},
		},
		{
			"nil key get",
			[]testEntry{val(nil, 42)},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tbl := newTable()
			for _, entry := range tc.entries {
				if entry.exists {
					tbl.Set(entry.key, entry.val)
				}
			}

			for _, entry := range tc.entries {
				key := entry.key
				actualVal, ok := tbl.GetOK(key)
				if entry.exists {
					if !ok {
						t.Fatalf("expected the value %v for key %q to exist, but it did not", entry.val, key)
					} else {
						require.Equal(t, entry.val, actualVal)
					}
				}

				if !entry.exists && ok {
					t.Fatalf("expected the value for key %q to not exist, but it did: %v", key, actualVal)
				}
			}
		})
	}
}
