package core

import (
	"fmt"

	"github.com/inkyblackness/res/serial"
)

type TestingResource struct {
	name string
	path string

	data       []byte
	readLocks  int
	writeLocks int
}

func NewTestingResource(name string, path string, data []byte) *TestingResource {
	res := &TestingResource{
		name: name,
		path: path,
		data: data}

	return res
}

// Name returns the unique identifier - the file name of the resource.
func (res *TestingResource) Name() string {
	return res.name
}

// Path returns the (relative) path for the resource, based on the release's root.
func (res *TestingResource) Path() string {
	return res.path
}

// AsSource returns an interface for reading the resource.
func (res *TestingResource) AsSource() (buf serial.SeekingReadCloser, err error) {
	if res.writeLocks == 0 {
		res.readLocks++
		buf = serial.NewByteStoreFromData(res.data, func([]byte) { res.readLocks-- })
	} else {
		err = fmt.Errorf("Cannot open for reading")
	}

	return
}

// AsSink returns an interface for writing the resource.
func (res *TestingResource) AsSink() (buf serial.SeekingWriteCloser, err error) {
	if res.readLocks == 0 && res.writeLocks == 0 {
		res.writeLocks++
		buf = serial.NewByteStoreFromData(res.data, func(data []byte) {
			res.data = data
			res.writeLocks--
		})
	} else {
		err = fmt.Errorf("Cannot open for writing")
	}

	return
}
