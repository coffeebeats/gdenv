package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/pkg/godot/artifact/source"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

/* ---------------------------- Test: Executable ---------------------------- */

func TestExecutable(t *testing.T) {
	ex := executable.MustParse("Godot_v4.0-stable_linux.x86_64")

	tests := []struct {
		store string
		ex    executable.Executable

		want string
		err  error
	}{
		{
			store: "",
			ex:    ex,

			err: ErrMissingStore,
		},

		{
			store: storeName,
			ex:    ex,

			want: filepath.Join(
				storeName,
				storeDirEx,
				"v4.0-stable",
				"linux.x86_64",
				ex.Name(),
			),
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s-%s", tc.store, tc.ex.String()), func(t *testing.T) {
			// When: The path to the cached executable is determined.
			got, err := Executable(tc.store, tc.ex)

			// Then: The expected error value is returned.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %s, want: %v", err, tc.err)
			}

			// Then: The expected filepath is returned.
			if got != tc.want {
				t.Errorf("output: got %s, want: %v", got, tc.want)
			}
		})
	}
}

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

/* ------------------------------ Test: Source ------------------------------ */

func TestSource(t *testing.T) {
	srcArchive := source.Archive{Inner: source.New(version.Godot4())}

	tests := []struct {
		store string
		src   source.Archive

		want string
		err  error
	}{
		{
			store: "",
			src:   srcArchive,

			err: ErrMissingStore,
		},

		{
			store: storeName,
			src:   srcArchive,

			want: filepath.Join(
				storeName,
				storeDirSrc,
				"v4.0-stable",
				srcArchive.Name(),
			),
		},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s-%s", tc.store, tc.src.Name()), func(t *testing.T) {
			// When: The path to the cached source directory is determined.
			got, err := Source(tc.store, tc.src.Inner)

			// Then: The expected error value is returned.
			if !errors.Is(err, tc.err) {
				t.Errorf("err: got %s, want: %v", err, tc.err)
			}

			// Then: The expected filepath is returned.
			if got != tc.want {
				t.Errorf("output: got %s, want: %v", got, tc.want)
			}
		})
	}
}
