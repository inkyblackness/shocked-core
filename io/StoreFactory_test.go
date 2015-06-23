package io

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/chunk/dos"
	"github.com/inkyblackness/shocked-core/release"
)

import (
	check "gopkg.in/check.v1"
)

type StoreFactorySuite struct {
	source  release.Release
	sink    release.Release
	factory *StoreFactory
}

var _ = check.Suite(&StoreFactorySuite{})

func (suite *StoreFactorySuite) SetUpTest(c *check.C) {
	suite.source = release.NewMemoryRelease()
	suite.sink = release.NewMemoryRelease()
	suite.factory = NewStoreFactory(suite.source, suite.sink)
}

func (suite *StoreFactorySuite) TestNewChunkStoreIsBackedBySinkIfExisting(c *check.C) {
	resource, _ := suite.sink.NewResource("fromSink.res", "")

	{
		writer, _ := resource.AsSink()
		consumer := dos.NewChunkConsumer(writer)
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
		consumer.Finish()
	}
	store, err := suite.factory.NewChunkStore("fromSink.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *StoreFactorySuite) TestNewChunkStoreIsBackedBySinkIfExistingInBoth(c *check.C) {
	resource1, _ := suite.source.NewResource("stillFromSink.res", "")

	{
		writer, _ := resource1.AsSink()
		consumer := dos.NewChunkConsumer(writer)
		consumer.Finish()
	}
	resource2, _ := suite.sink.NewResource("stillFromSink.res", "")

	{
		writer, _ := resource2.AsSink()
		consumer := dos.NewChunkConsumer(writer)
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
		consumer.Finish()
	}
	store, err := suite.factory.NewChunkStore("stillFromSink.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *StoreFactorySuite) TestNewChunkStoreIsBackedBySourceIfExisting(c *check.C) {
	resource, _ := suite.source.NewResource("fromSource.res", "")

	{
		writer, _ := resource.AsSink()
		consumer := dos.NewChunkConsumer(writer)
		consumer.Consume(res.ResourceID(1), chunk.NewBlockHolder(chunk.BasicChunkType, res.Palette, [][]byte{[]byte{}}))
		consumer.Finish()
	}
	store, err := suite.factory.NewChunkStore("fromSource.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	blockStore := store.Get(res.ResourceID(1))
	c.Check(blockStore.BlockCount(), check.Equals, uint16(1))
}

func (suite *StoreFactorySuite) TestNewChunkStoreReturnsEmptyStoreIfNowhereExisting(c *check.C) {
	store, err := suite.factory.NewChunkStore("empty.res")

	c.Assert(err, check.IsNil)
	c.Assert(store, check.NotNil)
	ids := store.IDs()
	c.Check(len(ids), check.Equals, 0)
}
