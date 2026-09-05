// Package checksumutil implements checksum computation and verification
// against a checksums file listing '<checksum> <filename>' entries.
//
// The logic here is deliberately free of any artifact typing; callers which
// need a typed API (see 'pkg/godot/artifact/checksum') wrap these functions
// rather than reimplementing them.
package checksumutil

import (
	"bufio"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"

	"github.com/coffeebeats/gdenv/internal/ioutil"
)

const checksumEntryParts = 2

var (
	ErrMismatch           = errors.New("checksum does not match")
	ErrNotFound           = errors.New("checksum not found")
	ErrConflicting        = errors.New("conflicting checksum")
	ErrUnrecognizedFormat = errors.New("unrecognized format")
)

/* -------------------------------------------------------------------------- */
/*                              Function: Compute                             */
/* -------------------------------------------------------------------------- */

// Compute computes and returns the checksum of the file at the specified path.
func Compute(ctx context.Context, h hash.Hash, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer f.Close()

	if _, err := io.Copy(h, ioutil.NewReaderWithContext(ctx, f.Read)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

/* -------------------------------------------------------------------------- */
/*                               Function: Lookup                             */
/* -------------------------------------------------------------------------- */

// Lookup finds and returns the checksum recorded for 'name' within the
// checksums file at the specified path.
func Lookup(ctx context.Context, checksumsPath, name string) (string, error) {
	f, err := os.Open(checksumsPath)
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
			return "", fmt.Errorf("%w: %s", ErrConflicting, n)
		}

		checksums[n] = c
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	checksum, has := checksums[name]
	if !has {
		return "", ErrNotFound
	}

	return checksum, nil
}

/* -------------------------------------------------------------------------- */
/*                              Function: Compare                             */
/* -------------------------------------------------------------------------- */

// Compare validates that the checksum of the file at 'path' matches the value
// recorded for 'name' within the checksums file at 'checksumsPath'.
func Compare(ctx context.Context, h hash.Hash, path, checksumsPath, name string) error {
	log.Info("verifying checksum of downloaded file")

	eg, ctx := errgroup.WithContext(ctx)

	got, want := make(chan string, 1), make(chan string, 1)
	defer close(got)
	defer close(want)

	eg.Go(func() error {
		value, err := Compute(ctx, h, path)
		if err != nil {
			return err
		}

		log.Debugf("actual checksum: %s", value)

		select {
		case got <- value:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	eg.Go(func() error {
		value, err := Lookup(ctx, checksumsPath, name)
		if err != nil {
			return err
		}

		log.Debugf("expected checksum: %s", value)

		select {
		case want <- value:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	if g, w := <-got, <-want; g != w {
		return fmt.Errorf("%w: %s (got) != %s (want)", ErrMismatch, g, w)
	}

	log.Debug("checksum matched expected value")

	return nil
}
