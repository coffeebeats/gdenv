package godot

import (
	"bufio"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	checksumEntryParts = 2
)

var (
	ErrConflictingChecksum = errors.New("conflicting checksum")
	ErrMissingChecksum     = errors.New("missing checksum")
	ErrMissingPath         = errors.New("missing path")
	ErrUnrecognizedFormat  = errors.New("unrecognized format")
)

/* -------------------------------------------------------------------------- */
/*                          Function: ComputeChecksum                         */
/* -------------------------------------------------------------------------- */

// Computes and returns the sha-512 checksum of the specified file.
func ComputeChecksum(path string) (string, error) {
	if path == "" {
		return "", ErrMissingPath
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

/* -------------------------------------------------------------------------- */
/*                          Function: ExtractChecksum                         */
/* -------------------------------------------------------------------------- */

// Returns the checksum matching the provided 'Executable' from the specified
// checksum file 'path'.
func ExtractChecksum(path string, ex Executable) (string, error) {
	if path == "" {
		return "", ErrMissingPath
	}

	name, err := ex.Name()
	if err != nil {
		return "", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Build a mapping from filenames to checksums. This enables detection of
	// conflicting entries (i.e. in case the file is malformed).
	scanner, checksums := bufio.NewScanner(f), make(map[string]string)
	for scanner.Scan() {
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

	// NOTE: This is brittle; find a better way of handling the extensions.
	checksum, has := checksums[name+".zip"]
	if !has {
		return "", ErrMissingChecksum
	}

	return checksum, nil
}
