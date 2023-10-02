package photosvc

import (
	"context"
	"net/http"
	"xform/entity"
)

type Router interface {
	Set(method, path string, handler http.HandlerFunc) (err error)
}

type Repo interface {
	GetPhotos(ctx context.Context) (photos entity.Photos, err error)
	UpsertPhotos(ctx context.Context, photos entity.Photos) (err error)
}
