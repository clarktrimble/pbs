// Package entity represents a photo, etc
package entity

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Geo struct {
	Lat float64
	Lon float64
	Alt float64
}

type Image struct {
	//Scale  int // Todo: size name ??
	Width  int
	Height int
	Path   string
}

type Photo struct {
	Id      string
	Name    string
	Path    string // Todo: path is not part of entity?
	TakenAt time.Time
	Geo     Geo
	Images  map[string]Image
}

func DecodePhoto(data []byte) (photo Photo, err error) {

	photo = Photo{}
	err = json.Unmarshal(data, &photo)
	err = errors.Wrapf(err, "failed to decode photo")
	return
}

func (photo Photo) Encode() (data []byte, err error) {

	data, err = json.Marshal(photo)
	err = errors.Wrapf(err, "somehow failed to encode photo")
	return
}

type Photos []Photo

func (photos Photos) String() string {

	data, err := json.MarshalIndent(photos, "", "  ")
	if err != nil {
		return "somehow failed to decode photos"
	}

	return string(data)
}

// Todo: new file plz

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
