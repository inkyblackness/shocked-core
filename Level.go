package core

import (
	"bytes"
	"encoding/binary"
	"sync"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/data"

	model "github.com/inkyblackness/shocked-model"
)

var tileTypes = map[data.TileType]model.TileType{
	data.Solid: model.Solid,
	data.Open:  model.Open,

	data.DiagonalOpenSouthEast: model.DiagonalOpenSouthEast,
	data.DiagonalOpenSouthWest: model.DiagonalOpenSouthWest,
	data.DiagonalOpenNorthWest: model.DiagonalOpenNorthWest,
	data.DiagonalOpenNorthEast: model.DiagonalOpenNorthEast,

	data.SlopeSouthToNorth: model.SlopeSouthToNorth,
	data.SlopeWestToEast:   model.SlopeWestToEast,
	data.SlopeNorthToSouth: model.SlopeNorthToSouth,
	data.SlopeEastToWest:   model.SlopeEastToWest,

	data.ValleySouthEastToNorthWest: model.ValleySouthEastToNorthWest,
	data.ValleySouthWestToNorthEast: model.ValleySouthWestToNorthEast,
	data.ValleyNorthWestToSouthEast: model.ValleyNorthWestToSouthEast,
	data.ValleyNorthEastToSouthWest: model.ValleyNorthEastToSouthWest,

	data.RidgeNorthWestToSouthEast: model.RidgeNorthWestToSouthEast,
	data.RidgeNorthEastToSouthWest: model.RidgeNorthEastToSouthWest,
	data.RidgeSouthEastToNorthWest: model.RidgeSouthEastToNorthWest,
	data.RidgeSouthWestToNorthEast: model.RidgeSouthWestToNorthEast}

func tileType(modelType model.TileType) (dataType data.TileType) {
	dataType = data.Solid

	for key, value := range tileTypes {
		if value == modelType {
			dataType = key
		}
	}

	return
}

type Level struct {
	id    int
	store chunk.Store

	mutex sync.Mutex

	tileMapStore chunk.BlockStore
	tileMap      []data.TileMapEntry
}

func NewLevel(store chunk.Store, id int) *Level {
	return &Level{id: id, store: store}
}

func (level *Level) bufferTileData() []data.TileMapEntry {
	if level.tileMap == nil {
		level.tileMap = make([]data.TileMapEntry, 64*64)

		level.tileMapStore = level.store.Get(res.ResourceID(4000 + level.id*100 + 5))
		blockData := level.tileMapStore.BlockData(0)
		reader := bytes.NewReader(blockData)
		binary.Read(reader, binary.LittleEndian, &level.tileMap)
	}

	return level.tileMap
}

func (level *Level) onTileDataChanged() {
	buf := bytes.NewBuffer(nil)

	binary.Write(buf, binary.LittleEndian, &level.tileMap)
	level.tileMapStore.SetBlockData(0, buf.Bytes())
}

func (level *Level) ID() int {
	return level.id
}

func (level *Level) Properties() (result model.LevelProperties) {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + 4)).BlockData(0)
	reader := bytes.NewReader(blockData)
	var info data.LevelInformation

	binary.Read(reader, binary.LittleEndian, &info)
	result.CyberspaceFlag = info.IsCyberspace()
	result.HeightShift = int(info.HeightShift)

	return
}

func (level *Level) Textures() (result []int) {
	blockData := level.store.Get(res.ResourceID(4000 + level.id*100 + 7)).BlockData(0)
	reader := bytes.NewReader(blockData)
	var ids [54]uint16

	binary.Read(reader, binary.LittleEndian, &ids)
	for _, id := range ids {
		result = append(result, int(id))
	}

	return
}

func (level *Level) TileProperties(x, y int) (result model.TileProperties) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	entries := level.bufferTileData()

	entry := entries[y*64+x]
	result.Type = new(model.TileType)
	*result.Type = tileTypes[entry.Type]
	result.SlopeHeight = new(model.HeightUnit)
	*result.SlopeHeight = model.HeightUnit(entry.SlopeHeight)
	result.FloorHeight = new(model.HeightUnit)
	*result.FloorHeight = model.HeightUnit(entry.Floor & 0x1F)
	result.CeilingHeight = new(model.HeightUnit)
	*result.CeilingHeight = model.HeightUnit(entry.Ceiling & 0x1F)

	{
		var properties model.RealWorldTileProperties
		var textureIDs = uint16(entry.Textures)

		properties.WallTexture = int(textureIDs & 0x3F)
		properties.CeilingTexture = int((textureIDs >> 6) & 0x1F)
		properties.CeilingTextureRotations = int((entry.Ceiling >> 5) & 0x03)
		properties.FloorTexture = int((textureIDs >> 11) & 0x1F)
		properties.FloorTextureRotations = int((entry.Floor >> 5) & 0x03)
		result.RealWorld = &properties
	}

	return
}

func (level *Level) SetTileProperties(x, y int, properties model.TileProperties) {
	level.mutex.Lock()
	defer level.mutex.Unlock()

	entries := level.bufferTileData()

	entry := &entries[y*64+x]
	if properties.Type != nil {
		entry.Type = tileType(*properties.Type)
	}
	if properties.FloorHeight != nil {
		entry.Floor = (entry.Floor & 0xE0) | (byte(*properties.FloorHeight) & 0x1F)
	}
	if properties.CeilingHeight != nil {
		entry.Ceiling = (entry.Ceiling & 0xE0) | (byte(*properties.CeilingHeight) & 0x1F)
	}
	if properties.SlopeHeight != nil {
		entry.SlopeHeight = byte(*properties.SlopeHeight)
	}

	level.onTileDataChanged()
}
