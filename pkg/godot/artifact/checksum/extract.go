package checksum

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/coffeebeats/gdenv/pkg/godot/artifact"
)

const checksumEntryParts = 2

var (
	ErrChecksumNotFound    = errors.New("checksum not found")
	ErrConflictingChecksum = errors.New("conflicting checksum")
	ErrUnrecognizedFormat  = errors.New("unrecognized format")
)

/* -------------------------------------------------------------------------- */
/*                              Function: Extract                             */
/* -------------------------------------------------------------------------- */

// Given a locally-available checksums file, find and return the checksum for
// the specified archive.
func Extract[T artifact.Artifact, U Checksums[T]](
	ctx context.Context,
	local artifact.Local[U],
	a T,
) (string, error) {
	f, err := os.Open(local.Path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	// Build a mapping from filenames to checksums. This enables detection of
	// conflicting entries (i.e. in case the file is malformed).
	scanner, checksums := bufio.NewScanner(f), make(map[string]string)
	for scanner.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		parts := strings.Fields(scanner.Text())
		if len(parts) != checksumEntryParts {
			return "", ErrUnrecognizedFormat
		}

		c, n := parts[0], parts[1]

		if existing, has := checksums[n]; has && existing != c {
			return "", fmt.Errorf("%w: %s", ErrConflictingChecksum, n)
		}

		checksums[n] = c
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	checksum, has := checksums[a.Name()]
	if !has {
		return "", ErrChecksumNotFound
	}

	return checksum, nil
}
