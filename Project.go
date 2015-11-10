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
}

func NewProject(name string, source release.Release, sink release.Release) (project *Project, err error) {
	library := io.NewReleaseStoreLibrary(source, sink, 5000)
	var textures *Textures
	var palettes *Palettes

	textures, err = NewTextures(library)

	if err == nil {
		palettes, err = NewPalettes(library)
	}

	if err == nil {
		project = &Project{
			name:     name,
			source:   source,
			sink:     sink,
			library:  library,
			textures: textures,
			palettes: palettes}
	}

	return
}

func (project *Project) Textures() *Textures {
	return project.textures
}

func (project *Project) Palettes() *Palettes {
	return project.palettes
}
