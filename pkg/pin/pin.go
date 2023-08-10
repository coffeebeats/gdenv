package pin

import (
	"errors"
	"io/fs"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot"
)

var (
	ErrIOFailed     = errors.New("pin: IO failed")
	ErrFileNotFound = errors.New("pin: file not found")
)

/* ----------------------------- Function: Read ----------------------------- */

// Parses a 'Version' from the specified pin file.
func Read(p string) (godot.Version, error) {
	p, err := Clean(p)
	if err != nil {
		return godot.Version{}, err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return godot.Version{}, errors.Join(ErrIOFailed, err)
	}

	return godot.ParseVersion(string(b))
}

/* ----------------------------- Function: Write ---------------------------- */

// Deletes the specified pin file.
func Remove(p string) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if err := os.Remove(p); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return errors.Join(ErrIOFailed, err)
		}
	}

	return nil
}

/* ----------------------------- Function: Write ---------------------------- */

// Writes a 'Version' to the specified pin file path.
func Write(v godot.Version, p string) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if err := os.WriteFile(p, []byte(v.String()), 0); err != nil {
		return errors.Join(ErrIOFailed, err)
	}

	return nil
}
