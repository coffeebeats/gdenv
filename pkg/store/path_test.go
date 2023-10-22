package store

import (
	"errors"
	"testing"
)

/* ------------------------------- Test: Path ------------------------------- */

func TestPath(t *testing.T) {
	tests := []struct {
		env  string
		want string
		err  error
	}{
		// Invalid inputs
		{env: "", err: ErrMissingEnvVar},
		{env: "a", err: ErrInvalidPath},
		{env: "a/b/c", err: ErrInvalidPath},
		{env: "/", err: ErrIllegalPath},
		{env: "/a", err: ErrIllegalPath},

		// Valid inputs
		{env: "/" + storeName, want: "/" + storeName},
		{env: "/." + storeName, want: "/." + storeName},
		{env: "/a/b/" + storeName, want: "/a/b/" + storeName},
		{env: "/a/b/." + storeName, want: "/a/b/." + storeName},
	}

	for _, tc := range tests {
		t.Run(tc.env, func(t *testing.T) {
			t.Setenv(envStore, tc.env)

			got, err := Path()

			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %#v, want %#v", err, tc.err)
			}

			if got != tc.want {
				t.Errorf("output: got %#v, want %#v", got, tc.want)
			}
		})
	}
}
