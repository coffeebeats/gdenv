package mirror

import (
	"github.com/coffeebeats/gdenv/pkg/godot"
	"github.com/go-resty/resty/v2"
)

/* -------------------------------------------------------------------------- */
/*                              Struct: TuxFamily                             */
/* -------------------------------------------------------------------------- */

// A mirror implementation for fetching artifacts via the Godot TuxFamily host.
type TuxFamily struct {
	client *resty.Client
}

/* ------------------------- Function: NewTuxFamily ------------------------- */

func NewTuxFamily() TuxFamily {
	return TuxFamily{}
}

/* ---------------------------- Method: Checksum ---------------------------- */

func (m *TuxFamily) Checksum(v godot.Version) (asset, error) {
	return asset{client: m.client}, nil
}

/* --------------------------- Method: Executable --------------------------- */

func (m *TuxFamily) Executable(ex godot.Executable) (asset, error) {
	return asset{client: m.client}, nil
}

/* ------------------------------- Method: Has ------------------------------ */

func (m *TuxFamily) Has(ex godot.Executable) bool {
	return false
}
