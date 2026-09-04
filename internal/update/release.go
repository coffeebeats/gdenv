package update

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/mod/semver"
	"golang.org/x/sync/errgroup"

	"github.com/coffeebeats/gdenv/internal/checksumutil"
	"github.com/coffeebeats/gdenv/internal/client"
	"github.com/coffeebeats/gdenv/internal/extract"
)

const (
	// repositoryURL is the repository from which releases are published.
	repositoryURL = "https://github.com/coffeebeats/gdenv"

	checksumsName = "checksums.txt"

	pathReleaseDownload = "releases/download"
	pathReleaseLatest   = "releases/latest"
	pathReleaseTag      = "releases/tag"
)

var (
	ErrDuplicateBinary = errors.New("duplicate binary in archive")
	ErrInvalidVersion  = errors.New("invalid version")
	ErrMissingBinary   = errors.New("missing binary in archive")
)

/* -------------------------------------------------------------------------- */
/*                           Function: LatestVersion                          */
/* -------------------------------------------------------------------------- */

// LatestVersion resolves the most recently published release version.
//
// NOTE: This reads the 'Location' header of the '/releases/latest' redirect
// rather than querying the GitHub API, which avoids both the unauthenticated
// rate limit and the need for a token.
func LatestVersion(ctx context.Context) (string, error) {
	return latestVersion(ctx, client.NewWithoutRedirects())
}

/* -------------------------- Function: latestVersion ----------------------- */

// latestVersion implements 'LatestVersion' for an explicit client.
func latestVersion(ctx context.Context, c *client.Client) (string, error) {
	u, err := c.Location(ctx, repositoryURL, pathReleaseLatest)
	if err != nil {
		return "", err
	}

	v := path.Base(u.Path)
	if !semver.IsValid(v) {
		return "", fmt.Errorf("%w: %s", ErrInvalidVersion, v)
	}

	return v, nil
}

/* -------------------------------------------------------------------------- */
/*                         Function: ReleaseNotesURL                          */
/* -------------------------------------------------------------------------- */

// ReleaseNotesURL returns the URL of the release notes for a version.
func ReleaseNotesURL(v string) string {
	return repositoryURL + "/" + pathReleaseTag + "/" + v
}

/* -------------------------------------------------------------------------- */
/*                            Function: Normalize                             */
/* -------------------------------------------------------------------------- */

// Normalize returns the canonical, 'v'-prefixed form of a version string, so
// that both '0.7.0' and 'v0.7.0' are accepted from users.
//
// NOTE: Releases are always tagged with a full 'vX.Y.Z', so a partial version
// is completed rather than passed through; 'v0.7' would otherwise validate and
// then fail as a 404 when the release asset is requested. Input which is not
// valid semver is returned 'v'-prefixed but otherwise unchanged, so that the
// caller can reject it and echo back what was typed.
func Normalize(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	v = "v" + strings.TrimPrefix(v, "v")

	if canonical := semver.Canonical(v); canonical != "" {
		return canonical
	}

	return v
}

/* -------------------------------------------------------------------------- */
/*                          Function: IsValidVersion                          */
/* -------------------------------------------------------------------------- */

// IsValidVersion reports whether the provided version string is a valid
// semantic version. Input should be passed through 'Normalize' first.
func IsValidVersion(v string) bool {
	return semver.IsValid(v)
}

/* -------------------------------------------------------------------------- */
/*                            Function: IsUpgrade                             */
/* -------------------------------------------------------------------------- */

// IsUpgrade reports whether 'latest' is a strictly newer version than 'current'.
//
// NOTE: GitHub resolves the "latest" release by publish date rather than by
// version, so a patch published against an older release branch can resolve to
// a *lower* version than the one running. Gating on this comparison ensures
// such a release is never presented as, or applied as, an upgrade.
func IsUpgrade(current, latest string) bool {
	if !semver.IsValid(current) || !semver.IsValid(latest) {
		return false
	}

	return semver.Compare(latest, current) > 0
}

/* --------------------------- Function: releaseHosts ----------------------- */

// releaseHosts returns the domains to which GitHub redirects release asset
// downloads; these must be permitted by the download client.
func releaseHosts() []string {
	return []string{
		"objects.githubusercontent.com",
		"release-assets.githubusercontent.com",
	}
}

/* -------------------------------------------------------------------------- */
/*                           Function: fetchArchive                           */
/* -------------------------------------------------------------------------- */

// fetchArchive downloads the release archive for the specified version into
// 'out' and verifies it against the published checksums file. The path to the
// verified archive is returned.
func fetchArchive(ctx context.Context, t Target, v, out string) (string, error) {
	name := t.ArchiveName(v)

	urlArchive, err := client.ParseURL(repositoryURL, pathReleaseDownload, v, name)
	if err != nil {
		return "", err
	}

	urlChecksums, err := client.ParseURL(repositoryURL, pathReleaseDownload, v, checksumsName)
	if err != nil {
		return "", err
	}

	pathArchive := filepath.Join(out, name)
	pathChecksums := filepath.Join(out, checksumsName)

	log.Infof("downloading 'gdenv' %s: %s", v, name)

	c := client.NewWithRedirectDomains(releaseHosts()...)

	eg, ctxDownload := errgroup.WithContext(ctx)

	eg.Go(func() error { return c.DownloadTo(ctxDownload, urlArchive, pathArchive) })
	eg.Go(func() error { return c.DownloadTo(ctxDownload, urlChecksums, pathChecksums) })

	if err := eg.Wait(); err != nil {
		return "", err
	}

	if err := checksumutil.Compare(ctx, sha256.New(), pathArchive, pathChecksums, name); err != nil {
		return "", err
	}

	return pathArchive, nil
}

/* -------------------------------------------------------------------------- */
/*                          Function: extractArchive                          */
/* -------------------------------------------------------------------------- */

// extractArchive extracts a downloaded release archive into 'out'.
func extractArchive(ctx context.Context, t Target, pathArchive, out string) error {
	if t.OS == osWindows {
		return extract.Zip(ctx, pathArchive, out)
	}

	return extract.TarGz(ctx, pathArchive, out)
}

/* -------------------------------------------------------------------------- */
/*                          Function: locateBinaries                          */
/* -------------------------------------------------------------------------- */

// locateBinaries searches 'dir' for each of the wanted executables and returns
// a mapping from name to the path at which it was found.
//
// NOTE: Entries are matched on their base name so that the directory layout
// within the release archive does not matter.
func locateBinaries(dir string, want []string) (map[string]string, error) {
	wanted := make(map[string]struct{}, len(want))
	for _, name := range want {
		wanted[name] = struct{}{}
	}

	found := make(map[string]string, len(want))

	if err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if _, ok := wanted[name]; !ok {
			return nil
		}

		if existing, ok := found[name]; ok {
			return fmt.Errorf("%w: %s: '%s' and '%s'", ErrDuplicateBinary, name, existing, p)
		}

		found[name] = p

		return nil
	}); err != nil {
		return nil, err
	}

	for _, name := range want {
		if _, ok := found[name]; !ok {
			return nil, fmt.Errorf("%w: %s", ErrMissingBinary, name)
		}
	}

	return found, nil
}
