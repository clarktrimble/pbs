package photosvc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"xform/entity"

	"github.com/clarktrimble/delish"
	"github.com/pkg/errors"
)

func NewAndReg(svr *delish.Server, rtr Router) (svc *photoSvc) {

	svc = &photoSvc{
		Server: svr,
	}

	rtr.Set("GET", "/photos", svc.getPhotos)
	rtr.Set("POST", "/photos", svc.addPhotos)
	return
}

// Todo: weird that svr is passed in above but already in struct below?
//rp := svc.Server.NewResponder(writer) found this commented in getAsdfdf below, clue??
// Yeah, prolly rename ^^^ New or Register
// Todo: copy new name to delish example please!!
// Todo: break into New and SetRoutes ??

// unexported

type photoSvc struct {
	Server *delish.Server
	photos []entity.Photo
}

// "largeURL": "http://tartu/photo/resized/PXL_20230707_101846985-4.png",
// "thumbnailURL": "http://tartu/photo/resized/PXL_20230707_101846985-16.png",
// "width": 3072,
// "height": 4080
// Todo: rename pkg swipesvc
type image struct {
	Large  string `json:"largeURL"`
	Thumb  string `json:"thumbnailURL"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func (svc *photoSvc) getPhotos(writer http.ResponseWriter, request *http.Request) {

	// Todo: bloga?
	writer.Header().Add("Access-Control-Allow-Origin", "*")
	//hdr.Add("Access-Control-Allow-Origin", "*")

	ctx := request.Context()
	rp := &delish.Respond{
		Writer: writer,
		Logger: svc.Server.Logger,
	}

	// Todo: xlate photos for swipe front-end

	images := []image{}
	for _, photo := range svc.photos {
		images = append(images, image{
			Large:  fmt.Sprintf("http://tartu/photo/resized/%s-4.png", photo.Name),
			Thumb:  fmt.Sprintf("http://tartu/photo/resized/%s-16.png", photo.Name),
			Width:  photo.Width,
			Height: photo.Height,
		})
	}

	rp.WriteObjects(ctx, map[string]any{"images": images})
}

func (svc *photoSvc) addPhotos(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	rp := &delish.Respond{
		Writer: writer,
		Logger: svc.Server.Logger,
	}

	photos := []entity.Photo{}
	err := json.NewDecoder(request.Body).Decode(&photos)
	if err != nil {
		err = errors.Wrapf(err, "failed decode")
		rp.NotOk(ctx, 500, err)
		return
	}

	svc.photos = photos
	// Todo: add Ok to responder
	rp.Write(ctx, []byte(`{"status":"ok"}`))
}
