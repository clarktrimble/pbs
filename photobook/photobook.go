// Package photobook relates to serving the photobook frontend app.
package photobook

import (
	"time"
	"xform/entity"
)

type PbItem struct {
	PhotoId  string    `json:"photo_id"`
	Source   string    `json:"src"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Thumb    string    `json:"thumb"`
	ThumbGs  string    `json:"thumb_gs"`
	Lat      float64   `json:"lat"`
	Lon      float64   `json:"lon"`
	TakenAt  time.Time `json:"taken_at"`
	Featured bool      `json:"featured"`
}

// ahh, yeah more naming plz
type PhotoBook []PbItem

func New(photos entity.Photos, book entity.Book) (pb PhotoBook) {

	//images := []image{}
	pb = PhotoBook{}

	for _, photo := range photos {
		pb = append(pb, PbItem{
			PhotoId:  photo.Id,
			Source:   photo.Images["large"].Path,
			Width:    photo.Images["large"].Width,
			Height:   photo.Images["large"].Height,
			Thumb:    photo.Images["thumb"].Path,
			ThumbGs:  photo.Images["thumb-gs"].Path, // Todo: fix Path ffs!
			Lat:      photo.Geo.Lat,
			Lon:      photo.Geo.Lon,
			TakenAt:  photo.TakenAt,
			Featured: book.Featured[photo.Id],
		})
	}

	return
}
