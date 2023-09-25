// Package entity represents a photo, etc
package entity

import "time"

type Geo struct {
	Lat float64
	Lon float64
	Alt float64
}

type Photo struct {
	Name    string
	Path    string
	Width   int
	Height  int
	TakenAt time.Time
	Geo     Geo
}
