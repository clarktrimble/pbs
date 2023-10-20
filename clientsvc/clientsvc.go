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

// Svc is an http client service layer.
type ClientSvc struct {
	Client Client
}

// PostPhotos posts photo objects to an api.
func (svc *ClientSvc) PostPhotos(ctx context.Context, photos []entity.Photo) (err error) {

	err = svc.Client.SendObject(ctx, "POST", path, photos, nil)
	return
}
