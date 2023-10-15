package pathutil

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var ErrUnknownMode = errors.New("cannot determine mode")

/* -------------------------------------------------------------------------- */
/*                             Function: Ancestor                             */
/* -------------------------------------------------------------------------- */

// Returns the closest ancestor of the specified 'path' which exists. If 'path'
// itself exists then it will be returned.
func Ancestor(ctx context.Context, path string) (string, error) {
	if path == "" {
		return "", fs.ErrInvalid
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		_, err := os.Stat(path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return "", err
			}

			path = filepath.Dir(path)
			if path == "" {
				return "", fs.ErrNotExist
			}

			continue
		}

		return path, nil
	}
}

/* -------------------------------------------------------------------------- */
/*                            Function: AncestorDir                           */
/* -------------------------------------------------------------------------- */

// Returns the closest ancestor directory of the specified 'path' which exists.
// If 'path' itself exists and is a directory then it will be returned.
func AncestorDir(ctx context.Context, path string) (string, error) {
	if path == "" {
		return "", fs.ErrInvalid
	}

	path, err := Ancestor(ctx, path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return path, nil
	}

	return filepath.Dir(path), nil
}

/* -------------------------------------------------------------------------- */
/*                           Function: AncestorMode                           */
/* -------------------------------------------------------------------------- */

// Returns the 'fs.FileMode' of the closest ancestor directory of the specified
// 'path' which exists. If 'path' itself exists and is a directory then it will
// be returned.
func AncestorMode(ctx context.Context, path string) (fs.FileMode, error) {
	ancestor, err := AncestorDir(ctx, path)
	if err != nil {
		return 0, errors.Join(ErrUnknownMode, err)
	}

	info, err := os.Stat(ancestor)
	if err != nil {
		return 0, errors.Join(ErrUnknownMode, err)
	}

	return info.Mode(), nil
}
