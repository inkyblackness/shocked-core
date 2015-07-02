package io

import (
	"github.com/inkyblackness/res/chunk"
)

// StoreFactory wraps the methods to create stores for various data
type StoreFactory interface {
	// NewChunkStore returns a chunk store for given name.
	NewChunkStore(name string) (chunk.Store, error)
}
