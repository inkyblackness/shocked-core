package core

import (
	"github.com/inkyblackness/shocked-core/release"
)

type Project struct {
	name   string
	source release.Release
	sink   release.Release
}

func NewProject(name string, source release.Release, sink release.Release) *Project {
	project := &Project{
		name:   name,
		source: source,
		sink:   sink}

	return project
}
