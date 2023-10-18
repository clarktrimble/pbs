// Package photosvc is a service layer between http server and repo.
package photosvc

import (
	"context"
	"encoding/json"
	"net/http"
	"xform/entity"
	"xform/photobook"

	"github.com/clarktrimble/delish"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// Todo: all photos are global, need a scoping concept, i.e.: "Baltic Travels"

// PhotoSvc is a service layer.
type PhotoSvc struct {
	Logger Logger
	Repo   Repo
}

// Register registers routes with the router.
func (svc *PhotoSvc) Register(rtr Router) {

	rtr.Set("GET", "/photos", svc.getPhotos)
	rtr.Set("POST", "/photos", svc.upsertPhotos)
	rtr.Set("POST", "/book", svc.upsertBook)
	rtr.Set("POST", "/featured", svc.setFeatured)
	rtr.Set("GET", "/photobook/{bookId}", svc.getPhotoBook)
}

// unexported

func (svc *PhotoSvc) respond(writer http.ResponseWriter, request *http.Request) (ctx context.Context, rp *delish.Respond) {

	ctx = request.Context()
	rp = &delish.Respond{
		Writer: writer,
		Logger: svc.Logger,
	}

	return
}

func (svc *PhotoSvc) getPhotos(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	photos, err := svc.Repo.GetPhotos()
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.WriteObjects(ctx, map[string]any{"photos": photos})
}

func (svc *PhotoSvc) upsertPhotos(writer http.ResponseWriter, request *http.Request) {

	ctx, rp := svc.respond(writer, request)

	photos, err := entity.ReadPhotos(request.Body)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	err = svc.Repo.UpsertPhotos(photos)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Ok(ctx)
}

func (svc *PhotoSvc) upsertBook(writer http.ResponseWriter, request *http.Request) {

	ctx, rp := svc.respond(writer, request)

	book := entity.Book{}
	err := json.NewDecoder(request.Body).Decode(&book)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}

	// note: all photos are unfeatured initially
	book.Featured = map[string]bool{}

	err = svc.Repo.UpsertBook(book)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Ok(ctx)
}

type featureParam struct {
	BookId   string `json:"book_id"`
	PhotoId  string `json:"photo_id"`
	Featured bool   `json:"featured"`
}

func (svc *PhotoSvc) setFeatured(writer http.ResponseWriter, request *http.Request) {

	// warning!: almost certainly not safe for concurrent use
	// should get away with it for a single user tho

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	param := featureParam{}
	err := json.NewDecoder(request.Body).Decode(&param)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}

	book, err := svc.Repo.GetBook(param.BookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	book.Featured[param.PhotoId] = param.Featured

	err = svc.Repo.UpsertBook(book)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Ok(ctx)
}

func (svc *PhotoSvc) getPhotoBook(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	// rustle up book and photos

	bookId := chi.URLParam(request, "bookId")

	book, err := svc.Repo.GetBook(bookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	photos, err := svc.Repo.GetPhotos()
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	// and mush together

	rp.WriteObjects(ctx, map[string]any{"images": photobook.New(photos, book)})
}
