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
	SizeName string
	Width    int
	Height   int
	Url      string // Todo: maybe url is not for here??
}

type Photo struct {
	Id      string
	Name    string
	TakenAt time.Time
	Geo     Geo
	Images  map[string]Image
}

type PhotoFile struct {
	Name string
	Path string
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
