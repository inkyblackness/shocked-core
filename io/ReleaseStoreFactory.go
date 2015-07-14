package io

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/res/chunk/store"
	"github.com/inkyblackness/res/serial"
	"github.com/inkyblackness/shocked-core/release"
)

type ReleaseStoreFactory struct {
	source release.Release
	sink   release.Release
}

func NewReleaseStoreFactory(source release.Release, sink release.Release) StoreFactory {
	factory := &ReleaseStoreFactory{
		source: source,
		sink:   sink}

	return factory
}

func (factory *ReleaseStoreFactory) NewChunkStore(name string) (chunkStore chunk.Store, err error) {
	if factory.sink.HasResource(name) {
		chunkStore, err = factory.openChunkStoreFrom(factory.sink, name)
	} else if factory.source.HasResource(name) {
		chunkStore, err = factory.openChunkStoreFrom(factory.source, name)
	} else {
		chunkStore = factory.createSavingStore(chunk.NullProvider(), "", name, func() {})
	}

	return
}

func (factory *ReleaseStoreFactory) openChunkStoreFrom(rel release.Release, name string) (chunkStore chunk.Store, err error) {
	resource, err := rel.GetResource(name)

	if err == nil {
		var reader serial.SeekingReadCloser
		reader, err = resource.AsSource()
		if err == nil {
			var provider chunk.Provider
			provider, err = dos.NewChunkProvider(reader)
			if err == nil {
				chunkStore = factory.createSavingStore(provider, resource.Path(), name, func() { reader.Close() })
			}
		}
	}

	return
}

func (factory *ReleaseStoreFactory) createSavingStore(provider chunk.Provider, path string, name string, closer func()) chunk.Store {
	nullStore := store.NewProviderBacked(chunk.NullProvider(), func() {})
	chunkStore := NewDynamicChunkStore(nullStore)

	closeLastReader := closer

	var resave func()
	resave = func() {
		go chunkStore.Swap(func(oldStore chunk.Store) chunk.Store {
			data := factory.serializeStore(oldStore)
			closeLastReader()

			newProvider, newReader := factory.saveAndReload(data, path, name)
			closeLastReader = func() { newReader.Close() }

			return store.NewProviderBacked(newProvider, resave)
		})
	}

	chunkStore.Swap(func(chunk.Store) chunk.Store {
		return store.NewProviderBacked(provider, resave)
	})

	return chunkStore
}

func (factory *ReleaseStoreFactory) serializeStore(store chunk.Store) []byte {
	buffer := serial.NewByteStore()
	consumer := dos.NewChunkConsumer(buffer)
	ids := store.IDs()

	for _, id := range ids {
		blockStore := store.Get(id)
		consumer.Consume(id, blockStore)
	}
	consumer.Finish()

	return buffer.Data()
}

func (factory *ReleaseStoreFactory) saveAndReload(data []byte, path string, name string) (provider chunk.Provider, reader serial.SeekingReadCloser) {
	var newResource release.Resource
	var err error

	if factory.sink.HasResource(name) {
		newResource, err = factory.sink.GetResource(name)
	} else {
		newResource, err = factory.sink.NewResource(name, path)
	}
	if err == nil {
		var newSink serial.SeekingWriteCloser

		newSink, err = newResource.AsSink()
		if err == nil {
			newSink.Write(data)
			newSink.Close()
		}
	}

	if err == nil {
		reader, err = newResource.AsSource()
	}
	if err == nil {
		provider, err = dos.NewChunkProvider(reader)
		if err != nil {
			reader.Close()
		}
	}
	if err != nil {
		reader = serial.NewByteStoreFromData(data, func([]byte) {})
		provider, _ = dos.NewChunkProvider(reader)
	}

	return
}
