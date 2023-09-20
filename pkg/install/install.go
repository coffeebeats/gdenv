package install

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/internal/mirror"
	"github.com/coffeebeats/gdenv/pkg/store"
)

const permUserReadWrite = 0700

var (
	ErrInvalidExecutable = errors.New("unsupported platform")
	ErrChecksumMismatch  = errors.New("checksum does not match")
)

/* -------------------------------------------------------------------------- */
/*                            Function: Executable                            */
/* -------------------------------------------------------------------------- */

// Downloads and caches a platform-specific version of Godot.
func Executable(versionRaw, out string) error { //nolint:funlen,cyclop
	// Validate arguments
	versionParsed, err := version.Parse(versionRaw)
	if err != nil {
		return err
	}

	// Define the host 'Platform'.
	platformParsed, err := DetectPlatform()
	if err != nil {
		return err
	}

	// Define the target 'Executable'.
	ex := executable.New(versionParsed, platformParsed)

	log.Println("Selecting mirror for executable:", ex.Name())
	log.Println("	> version:", versionParsed.String())

	m, err := ChooseMirror(versionParsed)
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("Successfully selected mirror: %T", m))

	// Create a temporary directory for the download.
	tmp, err := createTempDir()
	if err != nil {
		return err
	}

	// defer os.RemoveAll(tmp)

	errs := make(chan error)
	exPath, got, want := make(chan string, 1), make(chan string, 1), make(chan string, 1)

	// Download the Godot executable and compute the checksum.
	go func() {
		path, err := fetchExecutable(tmp, m, ex)
		if err != nil {
			errs <- err
			return
		}

		exPath <- path

		log.Println("Computing checksum: " + path)

		checksum, err := ComputeChecksum(path)
		if err != nil {
			errs <- err
			return
		}

		log.Println("Successfully computed checksum: " + checksum)

		got <- checksum
	}()

	// Download the checksums file for the version and extract the checksum.
	go func() {
		checksum, err := fetchChecksum(tmp, m, ex)
		if err != nil {
			errs <- err
			return
		}

		want <- checksum
	}()

	var p, g, w string

	for i := 0; i < 3; i++ {
		select {
		case p = <-exPath:
			// defer os.Remove(p)
		case g = <-got:
		case w = <-want:
		case err := <-errs:
			// NOTE: This approach drops additional errors. Consider fixing this
			// to join multiple errors.
			return err
		}
	}

	if g != w {
		return fmt.Errorf("%w: got '%s', want '%s'", ErrChecksumMismatch, g, w)
	}

	log.Println("Successfully compared checksums!")

	log.Println("Extracting executable from archive: " + p)

	// Unzip the archive now that it's been validated.
	extracted, err := archive.Extract[T](a, out)(p)
	if err != nil {
		return err
	}

	log.Println("Successfully extracted executable!")

	// Finally, add the extracted executable to the specified store.
	// TODO: Fix this.
	if _, err := os.Stat(filepath.Join(filepath.Dir(p), "Godot.app")); !errors.Is(err, fs.ErrNotExist) {
		if err := store.Add(storePath, filepath.Join(filepath.Dir(p), "Godot.app"), ex); err != nil {
			return err
		}

	} else {
		if err := store.Add(storePath, strings.TrimSuffix(p, filepath.Ext(p)), ex); err != nil {
			return err
		}
	}

	for _, path := range extracted {
		log.Println("Adding extracted file to store:", path)

		if err := store.Add(storePath, path, ex); err != nil {
			return err
		}
	}

	t, err := store.ToolPath(storePath, ex)
	if err != nil {
		return err
	}

	log.Println("Successfully added executable to store: " + t)

	return nil
}

/* ------------------------- Function: fetchChecksum ------------------------ */

func fetchChecksum(dir string, m mirror.Mirror, ex godot.Executable) (string, error) {
	asset, err := m.Checksum(ex.Version)
	if err != nil {
		return "", err
	}

	out := filepath.Join(dir, asset.Name())

	log.Println("Downloading asset: " + asset.Name())

	if err := m.DownloadTo(asset, out); err != nil {
		return "", err
	}
	// defer os.Remove(out)

	log.Println("Successfully downloaded asset: " + asset.Name())

	log.Println("Extracting checksum from: " + out)

	checksum, err := ExtractChecksum(out, ex)
	if err != nil {
		return "", err
	}

	log.Println("Successfully extracted checksum: " + checksum)

	return checksum, nil
}

/* ------------------------ Function: fetchExecutable ----------------------- */

func fetchExecutable(dir string, m mirror.Mirror, ex godot.Executable) (string, error) {
	asset, err := m.Executable(ex)
	if err != nil {
		return "", err
	}

	out := filepath.Join(dir, asset.Name())

	log.Println("Downloading asset: " + asset.Name())

	if err := m.DownloadTo(asset, out); err != nil {
		return "", err
	}

	log.Println("Successfully downloaded asset: " + asset.Name())

	return out, nil
}

/* ------------------------- Function: createTempDir ------------------------ */

// Creates a temporary directory useful for working with assets. Permissions
// will be set to '0666' (user R/W).
func createTempDir() (string, error) {
	tmp := filepath.Join(os.TempDir(), "gdenv")

	info, err := os.Stat(tmp)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}

		// Create a directory with user read and write permissions. Files placed
		// in here should *not* be executable and they don't need to be viewable
		// by anyone except the user running this process.
		//
		// NOTE: Don't use 'MkdirAll'; the system temp. directory should exist.
		if err := os.Mkdir(tmp, os.ModeDir|permUserReadWrite); err != nil {
			return "", err
		}
	}

	if info != nil && !info.IsDir() {
		return "", fs.ErrExist
	}

	return tmp, nil
}
