package core

import (
	"github.com/inkyblackness/res/chunk"
	"github.com/inkyblackness/res/text"
	"github.com/inkyblackness/shocked-core/io"
	model "github.com/inkyblackness/shocked-model"
)

// ElectronicMessages handles all data related to electronic messages.
type ElectronicMessages struct {
	cybstrng [model.LanguageCount]chunk.Store
	cp       text.Codepage
}

// NewElectronicMessages returns a new instance of ElectronicMessages.
func NewElectronicMessages(library io.StoreLibrary) (messages *ElectronicMessages, err error) {
	var cybstrng [model.LanguageCount]chunk.Store

	for i := 0; i < model.LanguageCount && err == nil; i++ {
		cybstrng[i], err = library.ChunkStore(localized[i].cybstrng)
	}
	if err == nil {
		messages = &ElectronicMessages{
			cybstrng: cybstrng,
			cp:       text.DefaultCodepage()}
	}

	return
}

// Message tries to retrieve the message data for given identification.
func (messages *ElectronicMessages) Message(messageType model.ElectronicMessageType, id int) (message model.ElectronicMessage, err error) {

	return
}

// SetMessage updates the properties of a message.
func (messages *ElectronicMessages) SetMessage(messageType model.ElectronicMessageType, id int, message model.ElectronicMessage) (err error) {

	return
}
