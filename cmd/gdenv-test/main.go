//go:debug tarinsecurepath=0
//go:debug zipinsecurepath=0
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/executable"
	"github.com/coffeebeats/gdenv/internal/godot/platform"
	"github.com/coffeebeats/gdenv/internal/godot/version"
	"github.com/coffeebeats/gdenv/internal/mirror/github"
)

func main() {
	m := github.New()

	ex := executable.New(version.MustParse("v4.1.1"), platform.MustParse("macos.universal"))

	remote, err := m.ExecutableArchive(ex)
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	basePath := filepath.Join(wd, "test")

	dstPath := filepath.Join(basePath, executable.Name(ex.Version(), ex.Platform()))

	log.Printf("Downloading asset '%s' to path: %s\n", remote.Name(), dstPath)
	log.Println("	> URL: ", remote.URL)

	// if err := m.DownloadTo(remote, dstPath); err != nil {
	// 	panic(err)
	// }

	log.Printf("Successfully downloaded asset '%s'!\n", remote.Name())

	local := artifact.Local[executable.Archive]{
		Artifact: remote.Artifact,
		Path:     dstPath,
	}

	log.Printf("Extracting asset '%s' to path: %s.\n", local.Name(), basePath)

	folder, err := archive.Extract(local, basePath)
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully extracted executable to path: %s\n", folder.Path)
}
