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

// Todo: short post on mutating golang slices, may need to look back in log..
// Todo: pull takeout from tarball!!
// Todo: look at using multiple cores!

type Size struct {
	Name  string
	Scale int
	Gs    bool
}

type Sizes []Size

func (sizes Sizes) BulkResize(dst string, photos []entity.Photo) (err error) {

	for _, photo := range photos {
		err = sizes.Resize(dst, photo)
		if err != nil {
			return
		}
	}

	return
}

func (sizes Sizes) Resize(dst string, photo entity.Photo) (err error) {

	img, err := imgio.Open(photo.Path)
	if err != nil {
		err = errors.Wrapf(err, "failed to open image")
		return
	}

	bounds := img.Bounds()
	wd := bounds.Dx()
	ht := bounds.Dy()

	fmt.Printf("-")
	for _, size := range sizes {

		sized := transform.Resize(img, wd/size.Scale, ht/size.Scale, transform.CatmullRom)
		if size.Gs {
			sized = effect.Grayscale(sized)
			// Todo: would save time to combo yeah?
		}

		out := path.Join(dst, fmt.Sprintf("%s-%s.png", photo.Name, size.Name))

		err = imgio.Save(out, sized, imgio.PNGEncoder())
		if err != nil {
			err = errors.Wrapf(err, "failed to save image")
			return
		}
		fmt.Printf(".")
	}

	return
}

var (
	baseUrl = "http://tartu/photo/resized"
	suffix  = "png"
)

func AddResize(photos entity.Photos, resizePath string, sizes Sizes) (err error) {

	for i, photo := range photos {

		images := map[string]entity.Image{}
		for _, size := range sizes {

			var wd, ht int

			filename := fmt.Sprintf("%s-%s.%s", photo.Name, size.Name, suffix)
			url := fmt.Sprintf("%s/%s", baseUrl, filename)
			path := fmt.Sprintf("%s/%s", resizePath, filename)
			wd, ht, err = getSize(path)
			if err != nil {
				return
			}

			images[size.Name] = entity.Image{
				Width:  wd,
				Height: ht,
				Path:   url, // Todo:!
			}
		}

		photos[i].Images = images
	}

	return
}

func getSize(imagePath string) (wd, ht int, err error) {

	rdr, err := os.Open(imagePath)
	if err != nil {
		err = errors.Wrapf(err, "failed to open image")
		return
	}
	defer rdr.Close()

	cfg, _, err := image.DecodeConfig(rdr)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode config")
		return
	}

	wd = cfg.Width
	ht = cfg.Height

	return
}
