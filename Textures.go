package core

import (
	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

type Textures struct {
	cybstrng [model.LanguageCount]chunk.Store
	cp       text.Codepage
}

func NewTextures(factory io.StoreFactory) (textures *Textures, err error) {
	var cybstrng [model.LanguageCount]chunk.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = factory.NewChunkStore(localized[i].cybstrng)
	}

	if err == nil {
		textures = &Textures{cybstrng: cybstrng, cp: text.DefaultCodepage()}
	}

	return
}

func (textures *Textures) TextureCount() int {
	return 273
}

func (textures *Textures) Properties(index int) model.TextureProperties {
	prop := model.TextureProperties{}

	for i := 0; i < model.LanguageCount; i++ {
		names := textures.cybstrng[i].Get(res.ResourceID(0x086A))
		cantBeUseds := textures.cybstrng[i].Get(res.ResourceID(0x086B))

		prop.Name[i] = textures.DecodeString(names.Get(uint16(index)))
		prop.CantBeUsed[i] = textures.DecodeString(cantBeUseds.Get(uint16(index)))
	}

	return prop
}

func (textures *Textures) DecodeString(data []byte) *string {
	value := textures.cp.Decode(data[0 : len(data)-1])

	return &value
}
