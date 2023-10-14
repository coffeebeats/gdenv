package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/coffeebeats/gdenv/internal/godot/artifact"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/archive"
	"github.com/coffeebeats/gdenv/internal/godot/artifact/source"
	"github.com/coffeebeats/gdenv/internal/godot/mirror"
	"github.com/coffeebeats/gdenv/internal/godot/version"
)

func main() {
	m := mirror.NewGitHub()

	remote, err := m.SourceArchive(version.MustParse("v4.1.1"))
	if err != nil {
		panic(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	basePath := filepath.Join(wd, "test")

	dstPath := filepath.Join(basePath, remote.Artifact.Name())

	log.Printf("Downloading asset '%s' to path: %s\n", remote.Artifact.Name(), basePath)
	log.Println("	> URL: ", remote.URL)

	if err := m.Client().DownloadTo(remote.URL, dstPath); err != nil {
		panic(err)
	}

	log.Printf("Successfully downloaded asset '%s'!\n", remote.Artifact.Name())

	local := artifact.Local[source.Archive]{
		Artifact: remote.Artifact,
		Path:     dstPath,
	}

	log.Printf("Extracting asset '%s' to path: %s.\n", local.Artifact.Name(), basePath)

	if err := archive.Extract(local, basePath); err != nil {
		panic(err)
	}

	log.Printf("Successfully extracted executable to path: %s\n", basePath)
}
