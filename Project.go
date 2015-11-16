package core

import (
	"github.com/inkyblackness/shocked-core/io"
	"github.com/inkyblackness/shocked-core/release"
)

type Project struct {
	name   string
	source release.Release
	sink   release.Release

	library io.StoreLibrary

	textures *Textures
	palettes *Palettes
	archive  *Archive
}

// NewProject creates a new project based on given release container.
func NewProject(name string, source release.Release, sink release.Release) (project *Project, err error) {
	library := io.NewReleaseStoreLibrary(source, sink, 5000)
	var textures *Textures
	var palettes *Palettes
	var archive *Archive

	textures, err = NewTextures(library)

	if err == nil {
		palettes, err = NewPalettes(library)
	}
	if err == nil {
		archive, err = NewArchive(library, "archive.dat")
	}

	if err == nil {
		project = &Project{
			name:     name,
			source:   source,
			sink:     sink,
			library:  library,
			textures: textures,
			palettes: palettes,
			archive:  archive}
	}

	return
}

// Name returns the name of the project.
func (project *Project) Name() string {
	return project.name
}

// Textures returns the wrapper for textures.
func (project *Project) Textures() *Textures {
	return project.textures
}

// Palettes returns the wrapper for palettes.
func (project *Project) Palettes() *Palettes {
	return project.palettes
}

// Archive returns the wrapper for the main archive file.
func (project *Project) Archive() *Archive {
	return project.archive
}
