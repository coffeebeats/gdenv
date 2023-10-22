package store

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/coffeebeats/gdenv/internal/fstest"
	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/artifacttest"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

const storeName = "store"

/* -------------------------------- Test: Add ------------------------------- */

func TestAdd(t *testing.T) {
	ex := executable.MustParse("Godot_v4.0-stable_linux.x86_64")

	storePathToEx := filepath.Join(storeName, storeDirEx, "v4.0-stable/linux.x86_64")
	storePathToSrc := filepath.Join(storeName, storeDirSrc, "v4.0-stable")

	tests := []struct {
		name  string
		add   []artifact.Local[artifact.Artifact]
		files []fstest.Writer

		want []fstest.Asserter
		err  error
	}{
		{
			name: "unsupported artifact returns error",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: artifacttest.MockArtifact{}, Path: "a/b/c"},
			},
			files: []fstest.Writer{
				fstest.File{Path: "a/b/c"},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: storePathToEx + "/c"},
			},
			err: ErrUnsupportedArtifact,
		},
		{
			name: "missing input file returns error",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: ex, Path: "a/b/c"},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: storePathToEx + "/c"},
			},
			err: fs.ErrNotExist,
		},
		{
			name: "single executable can be added into store",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: ex, Path: "a/b/c"},
			},
			files: []fstest.Writer{
				fstest.File{Path: "a/b/c"},
			},

			want: []fstest.Asserter{
				fstest.File{Path: storePathToEx + "/c"},
			},
		},
		{
			name: "multiple files can be added into store",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: ex, Path: "a/b/c"},
				{Artifact: ex, Path: "a/b/d"},
			},
			files: []fstest.Writer{
				fstest.File{Path: "a/b/c"},
				fstest.File{Path: "a/b/d"},
			},

			want: []fstest.Asserter{
				fstest.File{Path: storePathToEx + "/c"},
				fstest.File{Path: storePathToEx + "/d"},
			},
		},
		{
			name: "source folder can be added into store",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: source.New(ex.Version()), Path: "a/b/c"},
			},
			files: []fstest.Writer{
				fstest.Dir{Path: "a/b/c"},
			},

			want: []fstest.Asserter{
				fstest.Dir{Path: storePathToSrc + "/c"},
			},
		},
		{
			name: "new files overwrite existing files",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: ex, Path: "a"},
			},
			files: []fstest.Writer{
				fstest.File{Path: "a", Contents: "next"},
				fstest.File{Path: storePathToEx + "/a", Contents: "prev"},
			},

			want: []fstest.Asserter{
				fstest.File{Path: storePathToEx + "/a", Contents: "next"},
			},
		},
		{
			name: "a directory can be added into store",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: ex, Path: "a/b"},
			},
			files: []fstest.Writer{
				fstest.Dir{Path: "a/b"},
				fstest.File{Path: "a/b/c"},
			},

			want: []fstest.Asserter{
				fstest.Dir{Path: storePathToEx + "/b"},
				fstest.File{Path: storePathToEx + "/b/c"},
			},
		},
		{
			name: "new directories overwrite existing directories",
			add: []artifact.Local[artifact.Artifact]{
				{Artifact: source.New(ex.Version()), Path: "a"},
			},
			files: []fstest.Writer{
				fstest.File{Path: "a/next"},
				fstest.File{Path: storePathToSrc + "/a/prev"},
			},

			want: []fstest.Asserter{
				fstest.Dir{Path: storePathToSrc + "/a"},
				fstest.File{Path: storePathToSrc + "/a/next"},
				fstest.Absent{Path: storePathToSrc + "/a/prev"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified artifacts with their paths prefixed to the
			// temporary testing directoy.
			for i, a := range tc.add {
				tc.add[i].Path = filepath.Join(tmp, a.Path)
			}

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			// When: Artifacts are added to the store.
			// Then: The expected error value is returned.
			storePath := filepath.Join(tmp, storeName)
			if err := Add(storePath, tc.add...); !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected files exist on the file system.
			for _, f := range tc.want {
				f.Assert(t, tmp)
			}
		})
	}
}

/* ------------------------------- Test: Clear ------------------------------ */

func TestClear(t *testing.T) {
	tests := []struct {
		name  string
		files []fstest.Writer

		want []fstest.Asserter
		err  error
	}{
		{
			name: "clearing keeps store directories",

			want: []fstest.Asserter{
				fstest.Dir{Path: filepath.Join(storeName, storeDirBin)},
				fstest.Dir{Path: filepath.Join(storeName, storeDirEx)},
				fstest.Dir{Path: filepath.Join(storeName, storeDirSrc)},
				fstest.File{Path: filepath.Join(storeName, storeFileLayout)},
			},
		},
		{
			name: "clearing removes all cached artifacts",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, storeDirEx, "a")},
				fstest.File{Path: filepath.Join(storeName, storeDirEx, "b/c")},
				fstest.File{Path: filepath.Join(storeName, storeDirSrc, "a")},
				fstest.File{Path: filepath.Join(storeName, storeDirSrc, "b/c")},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: filepath.Join(storeName, storeDirEx, "a")},
				fstest.Absent{Path: filepath.Join(storeName, storeDirEx, "a/b")},
				fstest.Absent{Path: filepath.Join(storeName, storeDirSrc, "a")},
				fstest.Absent{Path: filepath.Join(storeName, storeDirSrc, "a/b")},
			},
		},
		{
			name: "clearing doesn't remove binary files",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, storeDirBin, "a")},
			},

			want: []fstest.Asserter{
				fstest.File{Path: filepath.Join(storeName, storeDirBin, "a")},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			// When: The cached artifacts are cleared from the store.
			// Then: The expected error value is returned.
			storePath := filepath.Join(tmp, storeName)
			if err := Clear(storePath); !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected files exist on the file system.
			for _, f := range tc.want {
				f.Assert(t, tmp)
			}
		})
	}
}

/* ---------------------------- Test: Executables --------------------------- */

func TestExecutables(t *testing.T) {
	tests := []struct {
		name  string
		files []fstest.Writer

		want []LocalEx
		err  error
	}{
		{
			name: "a nil slice is returned when store is empty",

			want: nil,
			err:  nil,
		},

		{
			name: "an empty slice is returned when version/platform directories are empty",
			files: []fstest.Writer{
				fstest.Dir{Path: filepath.Join(storeName, storeDirEx, version.Godot3().String())},
				fstest.Dir{Path: filepath.Join(storeName, storeDirEx, version.Godot4().String(), "linux.x86_64")},
			},

			want: []LocalEx{},
		},

		{
			name: "non-executable files are ignored",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, storeDirBin, "gdenv-cli")},
				fstest.File{Path: filepath.Join(storeName, storeDirSrc, version.Godot3().String(), "godot")},
				fstest.File{Path: filepath.Join(storeName, storeFileLayout)},
				fstest.File{Path: filepath.Join(storeName, ".godot-version")},

				fstest.Dir{Path: filepath.Join(storeName, storeDirEx)},
			},

			want: []LocalEx{},
		},

		{
			name: "all executables are included if they exist",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, storeDirBin, "gdenv-cli")},
				fstest.File{Path: filepath.Join(storeName, storeDirSrc, version.Godot3().String(), "godot")},
				fstest.File{Path: filepath.Join(storeName, storeFileLayout)},
				fstest.File{Path: filepath.Join(storeName, ".godot-version")},

				fstest.File{
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v3.0-stable",
						"x11.64",
						"Godot_v3.0-stable_x11.64",
					),
				},
				fstest.File{
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v3.0-stable",
						"win64",
						"Godot_v3.0-stable_win64.exe",
					),
				},
				fstest.File{
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v4.0-stable",
						"linux.x86_64",
						"Godot_v4.0-stable_linux.x86_64",
					),
				},
				fstest.File{
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v4.0-stable",
						"macos.universal",
						"Godot.app/Contents/MacOS/Godot",
					),
				},
			},

			// NOTE: These results must be alphabetically sorted by path elements.
			want: []LocalEx{
				artifact.Local[executable.Executable]{
					Artifact: executable.MustParse("Godot_v3.0-stable_win64.exe"),
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v3.0-stable",
						"win64",
						"Godot_v3.0-stable_win64.exe",
					),
				},
				artifact.Local[executable.Executable]{
					Artifact: executable.MustParse("Godot_v3.0-stable_x11.64"),
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v3.0-stable",
						"x11.64",
						"Godot_v3.0-stable_x11.64",
					),
				},
				artifact.Local[executable.Executable]{
					Artifact: executable.MustParse("Godot_v4.0-stable_linux.x86_64"),
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v4.0-stable",
						"linux.x86_64",
						"Godot_v4.0-stable_linux.x86_64",
					),
				},
				artifact.Local[executable.Executable]{
					Artifact: executable.MustParse("Godot_v4.0-stable_macos.universal"),
					Path: filepath.Join(
						storeName,
						storeDirEx,
						"v4.0-stable",
						"macos.universal",
						"Godot.app",
					),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified artifacts with their paths prefixed to the
			// temporary testing directoy.
			for i, a := range tc.want {
				tc.want[i].Path = filepath.Join(tmp, a.Path)
			}

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			// When: The cached artifacts in the store are listed.
			// Then: The expected error value is returned.
			got, err := Executables(context.Background(), filepath.Join(tmp, storeName))
			if !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected artifact references are returned.
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("count: got %v, want: %v", got, tc.want)
			}
		})
	}
}

/* -------------------------------- Test: Has ------------------------------- */

func TestHas(t *testing.T) {
	ex := executable.MustParse("Godot_v4.0-stable_linux.x86_64")
	src := source.New(ex.Version())

	storePathToEx := filepath.Join(storeName, storeDirEx, "v4.0-stable/linux.x86_64")
	storePathToSrc := filepath.Join(storeName, storeDirSrc, "v4.0-stable")

	tests := []struct {
		name     string
		artifact artifact.Artifact
		files    []fstest.Writer

		want bool
		err  error
	}{
		{
			name:     "unsupported artifact returns false",
			artifact: artifacttest.MockArtifact{},

			want: false,
		},
		{
			name:     "missing executable returns false",
			artifact: ex,

			want: false,
		},
		{
			name:     "present executable returns true",
			artifact: ex,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToEx, ex.Path())},
			},

			want: true,
		},
		{
			name:     "missing source returns false",
			artifact: src,

			want: false,
		},
		{
			name:     "present source returns true",
			artifact: src,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToSrc, src.Name())},
			},

			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			// When: Artifacts are added to the store.
			// Then: The expected error value is returned.
			ok, err := Has(filepath.Join(tmp, storeName), tc.artifact)
			if !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected result is returned.
			if ok != tc.want {
				t.Errorf("got: %v, want: %v", ok, tc.want)
			}
		})
	}
}

/* ------------------------------ Test: Remove ------------------------------ */

func TestRemove(t *testing.T) {
	ex := executable.MustParse("Godot_v4.0-stable_linux.x86_64")
	src := source.New(ex.Version())

	storePathToEx := filepath.Join(storeName, storeDirEx, "v4.0-stable/linux.x86_64")
	storePathToSrc := filepath.Join(storeName, storeDirSrc, "v4.0-stable")

	tests := []struct {
		name   string
		remove artifact.Artifact
		files  []fstest.Writer

		want []fstest.Asserter
		err  error
	}{
		{
			name:   "unsupported artifact is no-op",
			remove: artifacttest.MockArtifact{},

			err: nil,
		},
		{
			name:   "removing missing artifact is a no-op",
			remove: ex,
			files:  []fstest.Writer{},

			err: nil,
		},
		{
			name:   "remove executable deletes artifact",
			remove: ex,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToEx, ex.Path())},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: filepath.Join(storePathToEx, ex.Path())},
			},
		},
		{
			name:   "remove source deletes artifact",
			remove: src,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToSrc, src.Name())},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: filepath.Join(storePathToSrc, src.Name())},
			},
		},
		{
			name:   "remove executable deletes sibling, not cousin artifacts",
			remove: ex,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToEx, ex.Path())},
				fstest.File{Path: filepath.Join(storePathToEx, "sibling")},
				fstest.File{Path: filepath.Join(filepath.Dir(storePathToEx), "parent-sibling", "cousin")},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: filepath.Join(storePathToEx, ex.Path())},
				fstest.Absent{Path: filepath.Join(storePathToEx, "sibling")},
				fstest.File{Path: filepath.Join(filepath.Dir(storePathToEx), "parent-sibling", "cousin")},
			},
		},
		{
			name:   "remove source doesn't delete sibling artifact",
			remove: src,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToSrc, src.Name())},
				fstest.File{Path: filepath.Join(storePathToSrc, "sibling")},
			},

			want: []fstest.Asserter{
				fstest.Absent{Path: filepath.Join(storePathToSrc, src.Name())},
				fstest.File{Path: filepath.Join(storePathToSrc, "sibling")},
			},
		},
		{
			name:   "remove executable cleans up empty directory",
			remove: ex,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToEx, ex.Path())},
			},

			want: []fstest.Asserter{
				fstest.Dir{Path: filepath.Join(storeName, storeDirEx)},
				fstest.Absent{Path: filepath.Join(storeName, storeDirEx, ex.Version().String())},
			},
		},
		{
			name:   "remove source cleans up empty directory",
			remove: src,
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storePathToSrc, src.Name())},
			},

			want: []fstest.Asserter{
				fstest.Dir{Path: filepath.Join(storeName, storeDirSrc)},
				fstest.Absent{Path: filepath.Join(storeName, storeDirSrc, src.Version().String())},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			// When: The specified artifact is removed from the store.
			// Then: The expected error value is returned.
			storePath := filepath.Join(tmp, storeName)
			if err := Remove(storePath, tc.remove); !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected files exist on the file system.
			for _, f := range tc.want {
				f.Assert(t, tmp)
			}
		})
	}
}

/* ------------------------------- Test: Touch ------------------------------ */

func TestTouch(t *testing.T) {
	tests := []struct {
		name  string
		files []fstest.Writer

		want []fstest.Asserter
		err  error
	}{
		{
			name: "fully creates the store layout",

			want: []fstest.Asserter{
				fstest.Dir{Path: filepath.Join(storeName, storeDirBin)},
				fstest.Dir{Path: filepath.Join(storeName, storeDirEx)},
				fstest.Dir{Path: filepath.Join(storeName, storeDirSrc)},
				fstest.File{Path: filepath.Join(storeName, storeFileLayout)},
			},
		},
		{
			name: "adds missing folders without overriding existing",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, storeDirBin, "a")},
				fstest.File{Path: filepath.Join(storeName, storeDirEx, "a")},
			},

			want: []fstest.Asserter{
				fstest.File{Path: filepath.Join(storeName, storeDirBin, "a")},
				fstest.File{Path: filepath.Join(storeName, storeDirEx, "a")},
				fstest.Dir{Path: filepath.Join(storeName, storeDirSrc)},
				fstest.File{Path: filepath.Join(storeName, storeFileLayout)},
			},
		},
		{
			name: "overwrites a malformed layout file",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, storeFileLayout), Contents: "invalid"},
			},

			want: []fstest.Asserter{
				fstest.File{Path: filepath.Join(storeName, storeFileLayout), Contents: ""},
			},
		},
		{
			name: "doesn't overwrite a pin file in the store",
			files: []fstest.Writer{
				fstest.File{Path: filepath.Join(storeName, ".godot-version"), Contents: "v4.0-stable"},
			},

			want: []fstest.Asserter{
				fstest.File{Path: filepath.Join(storeName, ".godot-version"), Contents: "v4.0-stable"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			// Given: The specified files exist on the file system.
			for _, f := range tc.files {
				f.Write(t, tmp)
			}

			// When: A store is initialized at the specified path.
			// Then: The expected error value is returned.
			storePath := filepath.Join(tmp, storeName)
			if err := Touch(storePath); !errors.Is(err, tc.err) {
				t.Errorf("got: %v, want: %v", err, tc.err)
			}

			// Then: The expected files exist on the file system.
			for _, f := range tc.want {
				f.Assert(t, tmp)
			}
		})
	}
}
