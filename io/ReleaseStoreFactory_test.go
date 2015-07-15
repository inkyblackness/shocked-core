package io

import (
	"time"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/shocked-core/release"
)

import (
	check "gopkg.in/check.v1"
)

type ReleaseStoreFactorySuite struct {
	source  release.Release
	sink    release.Release
	factory StoreFactory
}

var _ = check.Suite(&ReleaseStoreFactorySuite{})

func (suite *ReleaseStoreFactorySuite) SetUpTest(c *check.C) {
	suite.source = release.NewMemoryRelease()
	suite.sink = release.NewMemoryRelease()
	suite.factory = NewReleaseStoreFactory(suite.source, suite.sink, 0)
}

func (suite *ReleaseStoreFactorySuite) createChunkResource(rel release.Release, name string, filler func(consumer chunk.Consumer)) {
	resource, _ := rel.NewResource(name, "")
	writer, _ := resource.AsSink()
	consumer := dos.NewChunkConsumer(writer)
	filler(consumer)
	consumer.Finish()
}

func (suite *ReleaseStoreFactorySuite) TestNewChunkStoreIsBackedBySinkIfExisting(c *check.C) {
	suite.createChunkResource(suite.sink, "fromSink.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.factory.NewChunkStore("fromSink.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *ReleaseStoreFactorySuite) TestNewChunkStoreIsBackedBySinkIfExistingInBoth(c *check.C) {
	suite.createChunkResource(suite.source, "stillFromSink.res", func(consumer chunk.Consumer) {})
	suite.createChunkResource(suite.sink, "stillFromSink.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.factory.NewChunkStore("stillFromSink.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *ReleaseStoreFactorySuite) TestNewChunkStoreIsBackedBySourceIfExisting(c *check.C) {
	suite.createChunkResource(suite.source, "fromSource.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.factory.NewChunkStore("fromSource.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *ReleaseStoreFactorySuite) TestNewChunkStoreReturnsEmptyStoreIfNowhereExisting(c *check.C) {
	store, err := suite.factory.NewChunkStore("empty.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	ids := store.IDs()
	c.Check(len(ids), check.Equals, 0)
}

func (suite *ReleaseStoreFactorySuite) TestModifyingSourceSavesNewSink(c *check.C) {
	suite.createChunkResource(suite.source, "source.res", func(consumer chunk.Consumer) {
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
	})
	store, err := suite.factory.NewChunkStore("source.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)

	store.Del(res.ResourceID(1))

	time.Sleep(100 * time.Millisecond)

	c.Check(suite.sink.HasResource("source.res"), check.Equals, true)
}
