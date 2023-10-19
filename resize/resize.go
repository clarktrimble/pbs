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

var (
	pngExt string = ".png"
)

// Size represents a plan to scale down an image.
type Size struct {
	Name  string
	Scale int
	Gs    bool
}

// Sizes is a multiplicity of Size.
type Sizes []Size

// ResizePhotos resizes a slice of photos.
func (sizes Sizes) ResizePhotos(dst string, photos []entity.PhotoFile) (err error) {

	for _, photo := range photos {
		err = sizes.Resize(dst, photo)
		if err != nil {
			return
		}
	}

	return
}

// Resize resizes a photo.
func (sizes Sizes) Resize(dst string, photo entity.PhotoFile) (err error) {

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
		}

		out := path.Join(dst, fmt.Sprintf("%s-%s%s", photo.Name, size.Name, pngExt))

		err = imgio.Save(out, sized, imgio.PNGEncoder())
		if err != nil {
			err = errors.Wrapf(err, "failed to save image")
			return
		}
		fmt.Printf(".")
	}

	return
}

// AddImages adds resized image data to photos.
func (sizes Sizes) AddImages(photos entity.Photos, resizePath string) (err error) {

	for i, photo := range photos {

		images := map[string]entity.Image{}
		for _, size := range sizes {

			var wd, ht int

			filename := fmt.Sprintf("%s-%s%s", photo.Name, size.Name, pngExt)
			path := fmt.Sprintf("%s/%s", resizePath, filename)
			wd, ht, err = getSize(path)
			if err != nil {
				return
			}

			images[size.Name] = entity.Image{
				SizeName: size.Name,
				Width:    wd,
				Height:   ht,
				Path:     filename,
			}
		}

		photos[i].Images = images
	}

	return
}

// unexported

func getSize(imagePath string) (wd, ht int, err error) {

	rdr, err := os.Open(imagePath)
	if err != nil {
		err = errors.Wrapf(err, "failed to open image")
		return
	}
	defer rdr.Close()

	// note: much faster than imgio (but seems to depend on it)
	cfg, _, err := image.DecodeConfig(rdr)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode config")
		return
	}

	wd = cfg.Width
	ht = cfg.Height

	return
}
