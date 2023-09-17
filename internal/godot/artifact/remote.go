package artifact

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrInvalidURL = errors.New("invalid URL")
	ErrMissingURL = errors.New("missing URL")
)

/* -------------------------------------------------------------------------- */
/*                               Struct: Remote                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is hosted on the internet and available
// for download.
type Remote[A Artifact] struct {
	Artifact A
	URL      *url.URL
}

/* --------------------------- Function: NewRemote -------------------------- */

// Returns a new 'Remote' artifact after validating inputs.
func NewRemote[A Artifact](artifact A, urlRaw string) (Remote[A], error) {
	var remote Remote[A]

	if urlRaw == "" {
		return remote, ErrMissingURL
	}

	// NOTE: Use the stricter 'ParseRequestURI' function instead of 'Parse'.
	urlParsed, err := url.ParseRequestURI(urlRaw)
	if err != nil {
		return remote, fmt.Errorf("%w: %s", ErrInvalidURL, urlRaw)
	}

	remote.Artifact, remote.URL = artifact, urlParsed

	return remote, nil
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (r Remote[A]) Name() string {
	return r.Artifact.Name()
}
