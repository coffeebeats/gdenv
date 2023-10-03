package archive

const extensionZip = ".zip"

/* -------------------------------------------------------------------------- */
/*                                 Struct: Zip                                */
/* -------------------------------------------------------------------------- */

// A struct representing a 'zip'-compressed archive.
type Zip[T Archivable] struct {
	Artifact T
}

/* ----------------------------- Impl: Artifact ----------------------------- */

func (a Zip[T]) Name() string {
	name := a.Artifact.Name()
	if name != "" {
		name += extensionZip
	}

	return name
}

/* ------------------------------ Impl: Archive ----------------------------- */

// Extracts the archived contents to the specified directory.
func (a Zip[T]) extract(path, out string) error { //nolint:revive
	return nil // TODO: Implement the archive extraction.
}

// // Extracts the archive to the specified file path.
// func (a Zip[T]) extract(archiveFilepath, out string) error {
// 	log.Println("Opening archive for decompression:", archiveFilepath)

// 	archive, err := zip.OpenReader(archiveFilepath)
// 	if err != nil {
// 		return err
// 	}

// 	defer archive.Close()

// 	log.Println("Successfully opened archive for decompression!")

// 	// Extract all files within the archive next to the archive.
// 	for _, f := range archive.File {
// 		log.Println("  >", f.Name)

// 		name := filepath.Clean(f.Name)
// 		// if filepath.IsAbs(name) {
// 		// 	return fmt.Errorf("%w: filepath '%s' is absolute", zip.ErrInsecurePath, name)
// 		// }

// 		// // Prevent an archived file trying to escape the archive.
// 		// if p, err := filepath.Rel(archiveFilepath, name); err == nil && strings.HasPrefix(p, "..") {
// 		// 	return fmt.Errorf("%w: path '%s' outside of archive", zip.ErrInsecurePath, name)
// 		// }

// 		mode := f.FileInfo().Mode()
// 		baseFilepath := filepath.Dir(archiveFilepath)
// 		dstFilepath := filepath.Join(baseFilepath, name)

// 		if f.FileInfo().IsDir() {
// 			if err := os.MkdirAll(filepath.Join(baseFilepath, name), mode); err != nil {
// 				return err
// 			}

// 			continue
// 		}

// 		if err := extractFile(archive, name, dstFilepath, mode); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // Extracts the specified file at a path relative to the directory containing
// // the archive.
// //
// // NOTE: This method does minimal validation of the filepaths within the Zip
// // archive. This is because 'gdenv' enables the 'zipinsecurepath' debug flag
// // which allows 'zip.Open' to do these extra security checks instead.
// func extractFile(a *zip.ReadCloser, name, dstFilepath string, mode fs.FileMode) error {
// 	src, err := a.Open(name)
// 	if err != nil {
// 		return err
// 	}

// 	defer src.Close()

// 	dst, err := os.OpenFile(
// 		// Extract the executable next to the archive.
// 		dstFilepath,
// 		// Only write to 'dst'; create a new file/overwrite an existing.
// 		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
// 		// Use the same permission (e.g. executable) as the source file.
// 		mode,
// 	)

// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()

// 	if _, err := io.Copy(dst, src); err != nil {
// 		return err
// 	}

// 	return nil
// }
