// Package clientsvc provides a photo client service layer.
package clientsvc

import (
	"context"
	"pbs/entity"
)

const (
	path string = "/photos"
)

// Client specifies an http client interface.
type Client interface {
	SendObject(ctx context.Context, method, path string, snd, rcv any) (err error)
}

// Svc represents a http client service layer.
type Svc struct {
	Client Client
}

// PostPhotos posts photo objects.
func (svc *Svc) PostPhotos(ctx context.Context, photos []entity.Photo) (err error) {

	err = svc.Client.SendObject(ctx, "POST", path, photos, nil)
	return
}
