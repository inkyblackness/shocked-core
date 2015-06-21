package core

import (
	"fmt"

	"github.com/inkyblackness/shocked-core/release"
)

type TestingReleaseContainer struct {
	releases map[string]*TestingRelease
}

func NewTestingReleaseContainer() *TestingReleaseContainer {
	return &TestingReleaseContainer{releases: make(map[string]*TestingRelease)}
}

// Names returns the list of currently known releases.
func (container *TestingReleaseContainer) Names() []string {
	names := make([]string, 0, len(container.releases))
	for name, _ := range container.releases {
		names = append(names, name)
	}

	return names
}

// Get returns the release with given name, or an error if not possible.
func (container *TestingReleaseContainer) Get(name string) (rel release.Release, err error) {
	rel, existing := container.releases[name]
	if !existing {
		err = fmt.Errorf("Not found")
	}

	return
}

// New creates a new release with given name and returns it, or an error if not possible.
func (container *TestingReleaseContainer) New(name string) (rel release.Release, err error) {
	rel, existing := container.releases[name]
	if !existing {
		testingRel := NewTestingRelease()
		rel = testingRel
		container.releases[name] = testingRel
	} else {
		err = fmt.Errorf("Release exists")
	}

	return
}
