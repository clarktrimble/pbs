// Package resize provides for resizing images
package resize

import (
	"fmt"
	"path"
	"xform/entity"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/pkg/errors"
)

// Todo: short post on mutating golang slices, may need to look back in log..
// Todo: pull takeout from tarball!!

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
