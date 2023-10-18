package photosvc

import (
	"context"
	"net/http"
	"xform/entity"
)

// Router specifies the router interface used here.
type Router interface {
	Set(method, path string, handler http.HandlerFunc)
}

// Repo specifies the backend db interface used here.
type Repo interface {
	GetPhotos() (photos entity.Photos, err error)
	UpsertPhotos(photos entity.Photos) (err error)
	GetBook(bookId string) (book entity.Book, err error)
	UpsertBook(book entity.Book) (err error)
}

// Logger specifies the logging interface used here.
type Logger interface {
	Info(ctx context.Context, msg string, kv ...interface{})
	Error(ctx context.Context, msg string, err error, kv ...interface{})
	WithFields(ctx context.Context, kv ...interface{}) context.Context
}
