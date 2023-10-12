package photosvc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"xform/entity"

	"github.com/clarktrimble/delish"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// Todo: validation of request body(s)
// Todo: all photos are global, need a scoping concept, i.e.: tartu Jun/Jul 23 or sommat
// Todo: add Ok to responder
// Todo: really want error from Set in router iface??

// PhotoSvc represents a servcie-layer ...
type PhotoSvc struct {
	Server *delish.Server
	Repo   Repo
}

// Register registers routes with the router.
func (svc *PhotoSvc) Register(rtr Router) {

	_ = rtr.Set("GET", "/photos", svc.getPhotos)
	_ = rtr.Set("POST", "/photos", svc.upsertPhotos)
	_ = rtr.Set("POST", "/book", svc.upsertBook)
	_ = rtr.Set("POST", "/featured", svc.setFeatured)
	_ = rtr.Set("GET", "/photobook/{bookId}", svc.getPhotoBook)
}

// unexported

func (svc *PhotoSvc) respond(writer http.ResponseWriter, request *http.Request) (ctx context.Context, rp *delish.Respond) {

	ctx = request.Context()
	rp = &delish.Respond{
		Writer: writer,
		Logger: svc.Server.Logger,
	}

	return
}

func (svc *PhotoSvc) getPhotos(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	photos, err := svc.Repo.GetPhotos(ctx)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.WriteObjects(ctx, map[string]any{"photos": photos})
}

func (svc *PhotoSvc) upsertPhotos(writer http.ResponseWriter, request *http.Request) {

	ctx, rp := svc.respond(writer, request)

	photos := entity.Photos{}
	err := json.NewDecoder(request.Body).Decode(&photos)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}

	err = svc.Repo.UpsertPhotos(ctx, photos)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Write(ctx, []byte(`{"status":"ok"}`))
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

	// Todo: all photos featured intially, or choose?
	book.Featured = map[string]bool{}

	err = svc.Repo.UpsertBook(ctx, book)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Write(ctx, []byte(`{"status":"ok"}`))
}

type featureParam struct {
	BookId   string `json:"book_id"`
	PhotoId  string `json:"photo_id"`
	Featured bool   `json:"featured"`
}

func (svc *PhotoSvc) setFeatured(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	param := featureParam{}
	err := json.NewDecoder(request.Body).Decode(&param)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}

	book, err := svc.Repo.GetBook(ctx, param.BookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	book.Featured[param.PhotoId] = param.Featured

	err = svc.Repo.UpsertBook(ctx, book)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Write(ctx, []byte(`{"status":"ok"}`))
}

type image struct {
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

func (svc *PhotoSvc) getPhotoBook(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	// rustle up book and photos

	bookId := chi.URLParam(request, "bookId")
	book, err := svc.Repo.GetBook(ctx, bookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	photos, err := svc.Repo.GetPhotos(ctx)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	// and mush together

	images := []image{}
	for _, photo := range photos {
		images = append(images, image{
			PhotoId:  photo.Id,
			Source:   photo.Images["large"].Url,
			Width:    photo.Images["large"].Width,
			Height:   photo.Images["large"].Height,
			Thumb:    photo.Images["thumb"].Url,
			ThumbGs:  photo.Images["thumb-gs"].Url, // Todo: choose -_ ffs!
			Lat:      photo.Geo.Lat,
			Lon:      photo.Geo.Lon,
			TakenAt:  photo.TakenAt,
			Featured: book.Featured[photo.Id],
		})
	}

	rp.WriteObjects(ctx, map[string]any{"images": images})
	// Todo: resolve naming issues (above generally) plz
}
