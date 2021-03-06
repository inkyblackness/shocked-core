package io

import (
	"github.com/inkyblackness/res/objprop"
	"github.com/inkyblackness/res/textprop"
)

// StoreLibrary wraps the methods to contain stores for various data
type StoreLibrary interface {
	// SaveAll requests to persist all pending modifications.
	SaveAll()

	// ChunkStore returns a chunk store for given name.
	ChunkStore(name string) (*DynamicChunkStore, error)

	// ObjpropStore returns an object properties store for given name.
	ObjpropStore(name string) (objprop.Store, error)

	// TextpropStore returns a texture properties store for given name.
	TextpropStore(name string) (textprop.Store, error)
}
