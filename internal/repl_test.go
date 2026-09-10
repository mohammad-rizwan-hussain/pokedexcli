package internal

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	cases := map[string][]struct {
		input    string
		expected []string
	}{
		"Split by space": {
			{input: "hello world", expected: []string{"hello", "world"}},
			{input: "foo bar baz", expected: []string{"foo", "bar", "baz"}},
		},
		"Trim whitespace": {
			{input: "   hello   world   ", expected: []string{"hello", "world"}},
			{input: "   foo   bar   baz   ", expected: []string{"foo", "bar", "baz"}},
		},
	}

	for testCaseName, tc := range cases {
		t.Run(testCaseName, func(t *testing.T) {
			for _, c := range tc {
				actual := CleanInput(c.input)
				if diff := cmp.Diff(c.expected, actual); diff != "" {
					t.Errorf("cleanInput(%q) mismatch (-want +got):\n%s", c.input, diff)
				}
			}
		})
	}
}
