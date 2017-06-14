package core

import (
	"bytes"
	"fmt"

	"github.com/inkyblackness/res"
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/image"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// Bitmaps is the adapter for general bitmaps.
type Bitmaps struct {
	mfdArt chunk.Store
}

// NewBitmaps returns a new Bitmaps instance, if possible.
func NewBitmaps(library io.StoreLibrary) (bitmaps *Bitmaps, err error) {
	var mfdArt chunk.Store

	if err == nil {
		mfdArt, err = library.ChunkStore("mfdart.res")
	}

	if err == nil {
		bitmaps = &Bitmaps{mfdArt: mfdArt}
	}

	return
}

// Image returns the image data of identified bitmap.
func (bitmaps *Bitmaps) Image(key model.ResourceKey) (bmp image.Bitmap, err error) {
	var blockData []byte

	if key.Type == model.ResourceTypeMfdDataImages {
		holder := bitmaps.mfdArt.Get(res.ResourceID(key.Type))
		if key.Index < holder.BlockCount() {
			blockData = holder.BlockData(key.Index)
		}
	}

	bmp, err = image.Read(bytes.NewReader(blockData))

	return
}

// SetImage requests to set the bitmap data of a resource.
func (bitmaps *Bitmaps) SetImage(key model.ResourceKey, bmp image.Bitmap) (resultKey model.ResourceKey, err error) {
	err = fmt.Errorf("Not Implemented")

	return
}
