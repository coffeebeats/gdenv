package mirror

import (
	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

/* -------------------------------------------------------------------------- */
/*                               Struct: GitHub                               */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via releases on the Godot
// GitHub repository.
type GitHub struct {
	client *resty.Client
}

/* --------------------------- Function: NewGitHub -------------------------- */

func NewGitHub() GitHub {
	return GitHub{}
}

/* ---------------------------- Method: Checksum ---------------------------- */

func (m *GitHub) Checksum(v godot.Version) (asset, error) {
	return asset{client: m.client}, nil
}

/* --------------------------- Method: Executable --------------------------- */

func (m *GitHub) Executable(ex godot.Executable) (asset, error) {
	return asset{client: m.client}, nil
}

/* ------------------------------- Method: Has ------------------------------ */

func (m *GitHub) Has(ex godot.Executable) bool {
	return false
}
