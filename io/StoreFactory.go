package io

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/res/chunk/store"
	"github.com/inkyblackness/shocked-core/release"
)

type StoreFactory struct {
	source release.Release
	sink   release.Release
}

func NewStoreFactory(source release.Release, sink release.Release) *StoreFactory {
	factory := &StoreFactory{
		source: source,
		sink:   sink}

	return factory
}

func (factory *StoreFactory) NewChunkStore(name string) (chunkStore chunk.Store, err error) {
	if factory.sink.HasResource(name) {
		chunkStore, err = factory.openChunkStoreFrom(factory.sink, name)
	} else if factory.source.HasResource(name) {
		chunkStore, err = factory.openChunkStoreFrom(factory.source, name)
	} else {
		chunkStore = store.NewProviderBacked(chunk.NullProvider(), func() {})
	}

	return
}

func (factory *StoreFactory) openChunkStoreFrom(rel release.Release, name string) (chunkStore chunk.Store, err error) {
	resource, err := rel.GetResource(name)

	if err == nil {
		reader, err := resource.AsSource()
		if err == nil {
			provider, err := dos.NewChunkProvider(reader)
			if err == nil {
				chunkStore = store.NewProviderBacked(provider, func() {})
			}
		}
	}

	return
}
