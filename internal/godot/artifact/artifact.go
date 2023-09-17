package artifact

import (
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

/* -------------------------------------------------------------------------- */
/*                             Interface: Artifact                            */
/* -------------------------------------------------------------------------- */

// An interface for different Godot-related files andfolder structures which
// 'gdenv' needs to interact with.
type Artifact interface {
	Name() string
}

/* -------------------------- Interface: Versioned -------------------------- */

// An interface for any artifacts which are tied to a specific version of Godot.
type Versioned interface {
	Version() version.Version
}

/* -------------------------- Interface: Platformed ------------------------- */

// An interface for any artifacts which are tied to a specific operating system
// and CPU architecture.
type Platformed interface {
	Platform() platform.Platform
}
