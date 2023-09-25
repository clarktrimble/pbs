// Package resize provides for resizing images
package resize

import (
	"fmt"
	"image"
	"os"
	"path"

	"xform/entity"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/pkg/errors"
)

type Size struct {
	Scale int
	Grey  bool
}

type Sizes []Size

func (sizes Sizes) Bulk(tmpDir string, photos []entity.Photo) (err error) {

	for i := range photos {
		err = sizes.Resize(tmpDir, &photos[i])
		if err != nil {
			return
		}
		fmt.Printf(">>> inner: %#v\n", photos[i])
	}

	return
}

// Todo: short post on mutating golang slices, may need to look back in log..

func (sizes Sizes) Resize(tmpDir string, photo *entity.Photo) (err error) {

	// mutate photo w/h, eek!

	//bounds := img.Bounds()
	//photo.Width = bounds.Dx()
	//photo.Height = bounds.Dy()
	rdr, err := os.Open(photo.Path)
	if err != nil {
		err = errors.Wrapf(err, "failed to open reader")
		return
	}
	defer rdr.Close()

	cfg, _, err := image.DecodeConfig(rdr)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode config")
		return
	}
	photo.Width = cfg.Width
	photo.Height = cfg.Height

	if tmpDir == "no save" {
		return
	}
	img, err := imgio.Open(photo.Path)
	if err != nil {
		err = errors.Wrapf(err, "failed to open image")
		return
	}

	for _, size := range sizes {

		sized := transform.Resize(img, photo.Width/size.Scale, photo.Height/size.Scale, transform.CatmullRom)

		var gs string
		if size.Grey {
			gs = "-gs"
			sized = effect.Grayscale(img)
		}

		out := path.Join(tmpDir, fmt.Sprintf("%s-%d%s.png", photo.Name, size.Scale, gs))

		err = imgio.Save(out, sized, imgio.PNGEncoder())
		if err != nil {
			err = errors.Wrapf(err, "failed to save image")
			return
		}
		fmt.Printf(".")
	}

	return
}

/*
func (size Size) Resize(tmpDir string, width, height, scale int) (err error) {
}

func Bulk(tmpDir string, photos *[]entity.Photo, scale int) (err error) {

	for _, photo := range *photos {

		out := path.Join(tmpDir, fmt.Sprintf("%s-%d-gs.png", photo.Name, scale))
		resize(photo.Path, out, scale)
	}

	return
}

func resizeOld(in, out string, scale int) {

	img, err := imgio.Open(in)
	if err != nil {
		//err = errors.Wrapf(err, "failed to open image")
		fmt.Printf("error: %s\n", err)
		return
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	gs := effect.Grayscale(img)
	sized := transform.Resize(gs, w/scale, h/scale, transform.CatmullRom)

	err = imgio.Save(out, sized, imgio.PNGEncoder())
	if err != nil {
		//err = errors.Wrapf(err, "failed to save image")
		fmt.Printf("error: %s\n", err)
		return
	}

	fmt.Printf(".")
}
*/
