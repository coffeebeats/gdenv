package download

import (
	"fmt"
	"io/fs"
	"os"
)

/* -------------------------------------------------------------------------- */
/*                         Function: checkIsDirectory                         */
/* -------------------------------------------------------------------------- */

// checkIsDirectory is a convenience function which returns whether the provided
// path is a directory.
func checkIsDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: expected a directory: '%s'", fs.ErrInvalid, path)
	}

	return nil
}
