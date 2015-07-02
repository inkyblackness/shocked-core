package core

import (
	"github.com/inkyblackness/shocked-core/io"
	"github.com/inkyblackness/shocked-core/release"
)

type Project struct {
	name   string
	source release.Release
	sink   release.Release

	factory io.StoreFactory

	textures *Textures
}

func NewProject(name string, source release.Release, sink release.Release) (project *Project, err error) {
	factory := io.NewReleaseStoreFactory(source, sink)
	textures, err := NewTextures(factory)

	if err == nil {
		project = &Project{
			name:     name,
			source:   source,
			sink:     sink,
			factory:  factory,
			textures: textures}
	}

	return
}

func (project *Project) Textures() *Textures {
	return project.textures
}
