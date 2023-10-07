// Package entity represents a photo, etc
package entity

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Book struct {
	Id       string
	Featured map[string]bool
}

func DecodeBook(data []byte) (book Book, err error) {

	book = Book{}
	err = json.Unmarshal(data, &book)
	err = errors.Wrapf(err, "failed to decode book")
	return
}

func (book Book) Encode() (data []byte, err error) {

	data, err = json.Marshal(book)
	err = errors.Wrapf(err, "somehow failed to encode book")
	return
}
