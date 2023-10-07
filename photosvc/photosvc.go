package photosvc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"xform/entity"

	"github.com/clarktrimble/delish"
	"github.com/pkg/errors"
)

// PhotoSvc represents a servcie-layer ...
type PhotoSvc struct {
	Server *delish.Server
	Repo   Repo
}

// Register registers routes with the router.
func (svc *PhotoSvc) Register(rtr Router) {

	rtr.Set("GET", "/photos", svc.getPhotos)
	rtr.Set("POST", "/photos", svc.upsertPhotos)
	rtr.Set("POST", "/book", svc.upsertBook)
	rtr.Set("POST", "/featured", svc.setFeatured)
	rtr.Set("GET", "/photobook", svc.getPhotoBook)
	rtr.Set("POST", "/photobook", svc.getPhotoBook) // Todo: grrrrrr
	//rtr.Set("POST", "/imgbook", svc.getImgBook)     // Todo: grrrrrr
	return
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

	// note: photos are not validated beyond simple unmarshal

	err = svc.Repo.UpsertPhotos(ctx, photos)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	// Todo: add Ok to responder
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

	// Todo: all photos featured intially
	book.Featured = map[string]bool{}

	// note: book is not validated beyond simple unmarshal

	err = svc.Repo.UpsertBook(ctx, book)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	rp.Write(ctx, []byte(`{"status":"ok"}`))
}

type photoBookRequest struct {
	BookId   string `json:"book_id"`
	PhotoId  string `json:"photo_id"`
	Featured bool   `json:"featured"`
}

func (svc *PhotoSvc) setFeatured(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	pbr := photoBookRequest{}
	err := json.NewDecoder(request.Body).Decode(&pbr)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}
	book, err := svc.Repo.GetBook(ctx, pbr.BookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	book.Featured[pbr.PhotoId] = pbr.Featured

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

	pbr := photoBookRequest{}
	err := json.NewDecoder(request.Body).Decode(&pbr)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}
	// Todo: validate plz prolly with decode over in nty
	book, err := svc.Repo.GetBook(ctx, pbr.BookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	// Todo: hmmm need a way to scope photos and book, something rotten here :/
	// Todo: params for gets ??

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
			Source:   photo.Images["large"].Path,
			Width:    photo.Images["large"].Width,
			Height:   photo.Images["large"].Height,
			Thumb:    photo.Images["thumb"].Path,
			ThumbGs:  photo.Images["thumb-gs"].Path, // Todo: choose -_ ffs!
			Lat:      photo.Geo.Lat,
			Lon:      photo.Geo.Lon,
			TakenAt:  photo.TakenAt,
			Featured: book.Featured[photo.Id],
		})
		//"TakenAt": "2023-07-07T19:03:16Z",
		//"Geo": {
		//"Lat": 54.897641699999994,
		//"Lon": 23.9223194,
	}

	rp.WriteObjects(ctx, map[string]any{"images": images})
	// Todo: resolve naming issues (above generally) plz
}

/*
func (svc *PhotoSvc) getImgBook(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx, rp := svc.respond(writer, request)

	// rustle up book and photos

	pbr := photoBookRequest{}
	err := json.NewDecoder(request.Body).Decode(&pbr)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 400, err)
		return
	}
	// Todo: validate plz prolly with decode over in nty
	book, err := svc.Repo.GetBook(ctx, pbr.BookId)
	if err != nil {
		rp.NotOk(ctx, 400, err)
		return
	}

	// Todo: hmmm need a way to scope photos and book, something rotten here :/
	// Todo: params for gets ??

	photos, err := svc.Repo.GetPhotos(ctx)
	if err != nil {
		rp.NotOk(ctx, 500, err)
		return
	}

	// and mush together

	img := image{}
	for _, photo := range photos {
		img = image{
			PhotoId:  photo.Id,
			Source:   photo.Images["large"].Path,
			Width:    photo.Images["large"].Width,
			Height:   photo.Images["large"].Height,
			Thumb:    photo.Images["thumb"].Path,
			ThumbGs:  photo.Images["thumb-gs"].Path, // Todo: choose -_ ffs!
			Lat:      photo.Geo.Lat,
			Lon:      photo.Geo.Lon,
			TakenAt:  photo.TakenAt,
			Featured: book.Featured[photo.Id],
		}
	}

	rp.WriteObjects(ctx, map[string]any{"image": img})
	// Todo: yeah this'n is pure garbage
	// implement url params next you scoundrel
}
*/
