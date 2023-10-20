// Package entity defines the entities of the project.
package entity

import (
	"time"
)

// PhotoFile tracks a photo in the fs.
type PhotoFile struct {
	Name string
	Path string
}

// Book tracks which photos are featured.
type Book struct {
	Id       string
	Featured map[string]bool
}

// PbItem is a photobook item tailored for use in frontend.
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

// PhotoBook is slice of photo book items used with photobook app.
type PhotoBook []PbItem

// New creates a photobook from photos and book.
func New(photos Photos, book Book) (pb PhotoBook) {

	pb = PhotoBook{}
	for _, photo := range photos {
		pb = append(pb, PbItem{
			PhotoId:  photo.Id,
			Source:   photo.Images["large"].Path,
			Width:    photo.Images["large"].Width,
			Height:   photo.Images["large"].Height,
			Thumb:    photo.Images["thumb"].Path,
			ThumbGs:  photo.Images["thumb-gs"].Path,
			Lat:      photo.Geo.Lat,
			Lon:      photo.Geo.Lon,
			TakenAt:  photo.TakenAt,
			Featured: book.Featured[photo.Id],
		})
	}

	return
}
