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
	GetBook(ctx context.Context, bookId string) (book entity.Book, err error)
	UpsertBook(ctx context.Context, book entity.Book) (err error)
}
