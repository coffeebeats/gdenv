package install

import (
	"errors"
	"fmt"
	"log"

	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/internal/mirror/github"
	"github.com/coffeebeats/gdenv/internal/mirror/tuxfamily"
)

var ErrNoMirrorFound = errors.New("no mirror found")

/* -------------------------------------------------------------------------- */
/*                           Function: ChooseMirror                           */
/* -------------------------------------------------------------------------- */

func ChooseMirror(v version.Version) (mirror.Mirror, error) { //nolint:ireturn
	log.Printf("Checking whether mirror 'GitHub' supports version: %s\n", v.String())

	// NOTE: Use the empty struct to avoid initializing a client before it's
	// necessary.
	if (github.GitHub{}).Supports(v) {
		m := github.New()
		if m.CheckIfSupports(v) {
			log.Printf("Success! Mirror 'GitHub' supports version: %s\n", v.String())
			return m, nil
		}

		log.Printf("Mirror 'GitHub' does not support version: %s\n", v.String())
	}

	log.Printf("Checking whether mirror 'TuxFamily' supports version: %s\n", v.String())

	// NOTE: Use the empty struct to avoid initializing a client before it's
	// necessary.
	if (tuxfamily.TuxFamily{}).Supports(v) {
		m := tuxfamily.New()
		if m.CheckIfSupports(v) {
			log.Printf("Success! Mirror 'TuxFamily' supports version: %s\n", v.String())
			return m, nil
		}

		log.Printf("Mirror 'TuxFamily' does not support version: %s\n", v.String())
	}

	return nil, fmt.Errorf("%w: %s", ErrNoMirrorFound, v)
}
