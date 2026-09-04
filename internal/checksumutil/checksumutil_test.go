package checksumutil_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/coffeebeats/gdenv/internal/checksumutil"
)

// contents is the archive payload used throughout these tests.
const contents = "gdenv-binary"

/* -------------------------------------------------------------------------- */
/*                            Function: TestCompute                           */
/* -------------------------------------------------------------------------- */

func TestCompute(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "archive.tar.gz")

	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := checksumutil.Compute(context.Background(), sha256.New(), path)
	if err != nil {
		t.Fatal(err)
	}

	want := sha256Of(t, contents)
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: TestLookup                            */
/* -------------------------------------------------------------------------- */

func TestLookup(t *testing.T) {
	tests := []struct {
		name       string
		checksums  string
		lookupName string
		want       string
		wantErr    error
	}{
		{
			name:       "finds the entry for the requested name",
			checksums:  "aaa  other.tar.gz\nbbb  archive.tar.gz\n",
			lookupName: "archive.tar.gz",
			want:       "bbb",
		},
		{
			name:       "missing entry is reported",
			checksums:  "aaa  other.tar.gz\n",
			lookupName: "archive.tar.gz",
			wantErr:    checksumutil.ErrNotFound,
		},
		{
			name:       "malformed line is rejected",
			checksums:  "aaa\n",
			lookupName: "archive.tar.gz",
			wantErr:    checksumutil.ErrUnrecognizedFormat,
		},
		{
			name:       "conflicting entries are rejected",
			checksums:  "aaa  archive.tar.gz\nbbb  archive.tar.gz\n",
			lookupName: "archive.tar.gz",
			wantErr:    checksumutil.ErrConflicting,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "checksums.txt")
			if err := os.WriteFile(path, []byte(tc.checksums), 0o600); err != nil {
				t.Fatal(err)
			}

			got, err := checksumutil.Lookup(context.Background(), path, tc.lookupName)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err: got %v, want %v", err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("err: got %v, want nil", err)
			}

			if got != tc.want {
				t.Fatalf("got %s, want %s", got, tc.want)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: TestCompare                           */
/* -------------------------------------------------------------------------- */

func TestCompare(t *testing.T) {
	tests := []struct {
		name    string
		record  string
		wantErr error
	}{
		{name: "matching checksum passes"},
		{
			name:    "mismatched checksum is reported",
			record:  "0000000000000000000000000000000000000000000000000000000000000000",
			wantErr: checksumutil.ErrMismatch,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()

			archive := filepath.Join(tmp, "archive.tar.gz")
			if err := os.WriteFile(archive, []byte(contents), 0o600); err != nil {
				t.Fatal(err)
			}

			record := tc.record
			if record == "" {
				record = sha256Of(t, contents)
			}

			checksums := filepath.Join(tmp, "checksums.txt")
			if err := os.WriteFile(checksums, []byte(record+"  archive.tar.gz\n"), 0o600); err != nil {
				t.Fatal(err)
			}

			err := checksumutil.Compare(
				context.Background(),
				sha256.New(),
				archive,
				checksums,
				"archive.tar.gz",
			)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err: got %v, want %v", err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("err: got %v, want nil", err)
			}
		})
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: sha256Of                              */
/* -------------------------------------------------------------------------- */

func sha256Of(t *testing.T, s string) string {
	t.Helper()

	h := sha256.New()
	if _, err := h.Write([]byte(s)); err != nil {
		t.Fatal(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}
