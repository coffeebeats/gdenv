package version

import (
	"fmt"
	"testing"
)

/* -------------------------- Test: Version.Normal -------------------------- */

func TestVersionNormal(t *testing.T) {
	type test struct {
		v    Version
		want string
	}

	tests := []test{
		{Version{}, "0.0.0"},

		{Version{major: 1}, "1.0.0"},
		{Version{major: 1, minor: 1}, "1.1.0"},
		{Version{major: 1, minor: 1, patch: 1}, "1.1.1"},

		{Version{minor: 1}, "0.1.0"},
		{Version{minor: 1, patch: 1}, "0.1.1"},

		{Version{patch: 1}, "0.0.1"},
	}

	// Produce an additional test with a specific label applied.
	withLabels := func(tc []test) []test {
		out := make([]test, len(tc))

		for i, t := range tc {
			v := t.v
			v.label = "label"
			out[i] = test{v, t.want}
		}

		return out
	}

	for i, tc := range append(tests, withLabels(tests)...) {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tc.v.Normal()

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}

/* -------------------------- Test: Version.String -------------------------- */

func TestVersionString(t *testing.T) {
	testLabel := "label"

	tests := []struct {
		v    Version
		want string
	}{
		// Default value
		{Version{}, "v0.0-" + labelDefault},

		// Default label
		{Version{major: 1}, "v1.0-" + labelDefault},
		{Version{major: 1, minor: 1}, "v1.1-" + labelDefault},
		{Version{major: 1, minor: 1, patch: 1}, "v1.1.1-" + labelDefault},

		{Version{minor: 1}, "v0.1-" + labelDefault},
		{Version{minor: 1, patch: 1}, "v0.1.1-" + labelDefault},

		{Version{patch: 1}, "v0.0.1-" + labelDefault},

		// Specific label
		{Version{label: testLabel}, "v0.0-" + testLabel},

		{Version{major: 1, label: testLabel}, "v1.0-" + testLabel},
		{Version{major: 1, minor: 1, label: testLabel}, "v1.1-" + testLabel},
		{Version{major: 1, minor: 1, patch: 1, label: testLabel}, "v1.1.1-" + testLabel},

		{Version{minor: 1, label: testLabel}, "v0.1-" + testLabel},
		{Version{minor: 1, patch: 1, label: testLabel}, "v0.1.1-" + testLabel},

		{Version{patch: 1, label: testLabel}, "v0.0.1-" + testLabel},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tc.v.String()

			if got != tc.want {
				t.Fatalf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
