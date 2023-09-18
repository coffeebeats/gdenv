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

type Hosted = Remote[Artifact]

/* -------------------------------------------------------------------------- */
/*                               Struct: Remote                               */
/* -------------------------------------------------------------------------- */

// A wrapper around an 'Artifact' which is hosted on the internet.
type Remote[T Artifact] struct {
	Artifact T
	URL      string
}

/* ---------------------------- Method: ParseURL ---------------------------- */

// Returns a parsed URL or fails if it's invalid.
func (r Remote[T]) ParseURL() (*url.URL, error) {
	if r.URL == "" {
		return nil, ErrMissingURL
	}

	// NOTE: Use the stricter 'ParseRequestURI' function instead of 'Parse'.
	urlParsed, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return nil, errors.Join(ErrInvalidURL, err)
	}

	if urlParsed.Host == "" || urlParsed.Scheme == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidURL, r.URL)
	}

	return urlParsed, nil
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (r Remote[T]) Name() string {
	return r.Artifact.Name()
}
