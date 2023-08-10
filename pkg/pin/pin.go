package pin

import (
	"errors"
	"fmt"
	"os"

	"github.com/coffeebeats/gdenv/internal/godot"
)

var (
	ErrFailedToRead  = errors.New("pin: failed to read")
	ErrFailedToWrite = errors.New("pin: failed to write")
	ErrFileNotFound  = errors.New("pin: file not found")
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
		return godot.Version{}, fmt.Errorf("%w: %w", ErrFailedToRead, err)
	}

	return godot.ParseVersion(string(b))
}

/* ----------------------------- Function: Write ---------------------------- */

// Writes a 'Version' to the specified pin file path.
func Write(v godot.Version, p string) error {
	p, err := Clean(p)
	if err != nil {
		return err
	}

	if err := os.WriteFile(p, []byte(v.String()), 0); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToWrite, err)
	}

	return nil
}
