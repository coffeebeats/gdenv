package artifact

import (
	"errors"
	"net/url"
	"os"

	"github.com/coffeebeats/gdenv/pkg/godot/platform"
	"github.com/coffeebeats/gdenv/pkg/godot/version"
)

var ErrMissingPath = errors.New("missing path")

/* -------------------------------------------------------------------------- */
/*                             Interface: Artifact                            */
/* -------------------------------------------------------------------------- */

// An interface for different Godot-related files andfolder structures which
// 'gdenv' needs to interact with.
type Artifact interface {
	Name() string
}

/* -------------------------- Interface: Platformed ------------------------- */

// An interface for any artifacts which are tied to a specific operating system
// and CPU architecture.
type Platformed interface {
	Artifact

	Platform() platform.Platform
}

/* -------------------------- Interface: Versioned -------------------------- */

// An interface for any artifacts which are tied to a specific version of Godot.
type Versioned interface {
	Artifact

	Version() version.Version
}

/* -------------------------------------------------------------------------- */
/*                                Struct: Local                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is locally-available on the file system.
type Local[T Artifact] struct {
	Artifact T
	Path     string
}

/* ----------------------------- Method: Exists ----------------------------- */

// Returns whether the downloaded file exists on the local file system.
func (l Local[T]) Exists() bool {
	if _, err := os.Stat(l.Path); err != nil {
		return false
	}

	return true
}

/* -------------------------------------------------------------------------- */
/*                               Struct: Remote                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is hosted on the internet.
type Remote[T Artifact] struct {
	Artifact T
	URL      *url.URL
}
