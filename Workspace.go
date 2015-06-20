package core

import (
	"github.com/inkyblackness/shocked-core/release"
)

type Workspace struct {
}

// NewWorkspace takes a Release as a basis for existing resources and returns
// a new workspace instance. With this instance, projects from given projects container
// can be worked with.
func NewWorkspace(source release.Release, projects release.ReleaseContainer) (ws *Workspace, err error) {
	return
}
