package mirror

import (
	"errors"
	"net/url"
)

var (
	ErrMissingName = errors.New("missing name")
	ErrMissingURL  = errors.New("missing URL")
)

/* -------------------------------------------------------------------------- */
/*                                Struct: Asset                               */
/* -------------------------------------------------------------------------- */

// A struct representing a file/directory that can be downloaded from a Godot
// project mirror.
type Asset struct {
	name string
	url  *url.URL
}

/* --------------------------- Function: NewAsset --------------------------- */

// Returns a new 'Asset' after validating inputs.
func NewAsset(name string, u *url.URL) (Asset, error) {
	var asset Asset

	if name == "" {
		return asset, ErrMissingName
	}

	if u == nil {
		return asset, ErrMissingURL
	}

	asset.name, asset.url = name, u

	return asset, nil
}

/* ------------------------------ Method: Name ------------------------------ */

// Returns the filename of the asset to download.
func (a *Asset) Name() string {
	return a.name
}

/* ------------------------------- Method: URL ------------------------------ */

// Returns the URL of the asset to download.
func (a *Asset) URL() *url.URL {
	return a.url
}
