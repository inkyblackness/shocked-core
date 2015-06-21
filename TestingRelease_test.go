package core

import (
	"fmt"

	"github.com/inkyblackness/shocked-core/release"
)

type TestingRelease struct {
	resources map[string]*TestingResource
}

func NewTestingRelease() *TestingRelease {
	return &TestingRelease{resources: make(map[string]*TestingResource)}
}

// HasResource returns true for a unique resource name if the release contains this resource.
func (rel *TestingRelease) HasResource(name string) bool {
	_, existing := rel.resources[name]

	return existing
}

// GetResource returns the resource identified by given name if existing, or an error otherwise.
func (rel *TestingRelease) GetResource(name string) (res release.Resource, err error) {
	res, existing := rel.resources[name]
	if !existing {
		err = fmt.Errorf("Not found")
	}

	return
}

// NewResource creates a new resource under given path and returns the instance, or an error on failure.
func (rel *TestingRelease) NewResource(name string, path string) (res release.Resource, err error) {
	res, existing := rel.resources[name]
	if !existing {
		testingRes := NewTestingResource(name, path, nil)
		res = testingRes
		rel.resources[name] = testingRes
	} else {
		err = fmt.Errorf("Resource exists")
	}

	return
}
