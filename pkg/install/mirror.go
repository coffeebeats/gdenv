package install

import (
	"errors"
	"fmt"
	"log"

	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

var ErrNoMirrorFound = errors.New("no mirror found")

/* -------------------------------------------------------------------------- */
/*                           Function: ChooseMirror                           */
/* -------------------------------------------------------------------------- */

func ChooseMirror(v version.Version) (mirror.Mirror, error) { //nolint:ireturn
	log.Printf("Checking whether mirror 'GitHub' supports version: %s\n", v.String())

	// NOTE: Use a zero value to avoid initializing a client before necessary.
	if (mirror.GitHub{}).Supports(v) {
		m := mirror.NewGitHub()
		if m.CheckIfSupports(v) {
			log.Printf("Success! Mirror 'GitHub' supports version: %s\n", v.String())
			return m, nil
		}

		log.Printf("Mirror 'GitHub' does not support version: %s\n", v.String())
	}

	log.Printf("Checking whether mirror 'TuxFamily' supports version: %s\n", v.String())

	// NOTE: Use a zero value to avoid initializing a client before necessary.
	if (mirror.TuxFamily{}).Supports(v) {
		m := mirror.NewTuxFamily()
		if m.CheckIfSupports(v) {
			log.Printf("Success! Mirror 'TuxFamily' supports version: %s\n", v.String())
			return m, nil
		}

		log.Printf("Mirror 'TuxFamily' does not support version: %s\n", v.String())
	}

	return nil, fmt.Errorf("%w: %s", ErrNoMirrorFound, v)
}
