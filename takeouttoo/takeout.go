// Package takeout soaks up from goog's takeout
package takeouttoo

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	pth "path"
	"strconv"
	"strings"
	"time"

	"xform/entity"

	"github.com/clarktrimble/hondo"
	"github.com/pkg/errors"
)

type Geo struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
	Alt float64 `json:"altitude"`
}

type Taken struct {
	Epoch string `json:"timestamp"`
}

// Photo as seen in Takeout files
type Photo struct {
	Title string `json:"title"`
	Taken Taken  `json:"photoTakenTime"`
	Geo   Geo    `json:"geoData"`
}

// FromFiles creates photos given ...
// func FromFiles(jsonPath, resizePath string, sizes resize.Sizes) (photos entity.Photos, err error) {
func FromFiles(jsonPath string) (photos entity.Photos, err error) {

	fmt.Printf(">>> %s\n", jsonPath)

	paths, err := findJson(jsonPath)
	if err != nil {
		return
	}

	photos = entity.Photos{}
	for _, path := range paths {

		var photo entity.Photo
		photo, err = decode(pth.Join(jsonPath, pth.Dir(path)), pth.Base(path))
		if err != nil {
			return
		}

		if !strings.HasSuffix(photo.Path, ".jpg") {
			fmt.Printf("skipping %s\n", photo.Path)
			continue
		}

		_, err = os.Stat(photo.Path)
		if err != nil {
			fmt.Printf("also skipping %s\n", photo.Path)
			continue
		}
		photos = append(photos, photo)
	}

	return
}

func findJson(root string) (paths []string, err error) {

	// Todo: maybe walk is overkill here??

	paths = []string{}
	err = fs.WalkDir(os.DirFS(root), ".", func(path string, entry fs.DirEntry, err error) error {

		if err != nil {
			err = errors.Wrapf(err, "walker called with err")
			return err
		}

		if strings.HasSuffix(path, ".json") {
			paths = append(paths, path)
		}

		return nil
	})

	err = errors.Wrapf(err, "failed to walk")
	return
}

// func decode(dir, name string, sizes resize.Sizes) (photo entity.Photo, err error) {
func decode(dir, name string) (photo entity.Photo, err error) {

	file, err := os.Open(pth.Join(dir, name))
	if err != nil {
		err = errors.Wrapf(err, "failed to open")
		return
	}

	to := &Photo{}
	err = json.NewDecoder(file).Decode(to)
	if err != nil {
		err = errors.Wrapf(err, "failed to decode")
		return
	}

	epoch, err := strconv.ParseInt(to.Taken.Epoch, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse epoch")
		return
	}

	split := strings.Split(to.Title, ".")

	// Todo: yeah split this out somewheres??

	photo = entity.Photo{
		Id:      hondo.Rand(7),
		Name:    split[0],
		Path:    pth.Join(dir, to.Title),
		TakenAt: time.Unix(epoch, 0).UTC(),
		Geo: entity.Geo{
			Lat: to.Geo.Lat,
			Lon: to.Geo.Lon,
			Alt: to.Geo.Alt,
		},
	}

	// Todo: helper ??

	//for _, size := range sizes {
	//fmt.Printf(">>> %s %s\n", photo.Name, size.Name)
	//}

	//fmt.Printf(">>> decoding .. %s\n", photo.Path)
	/*
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
	*/

	return
}

// Todo: move below to resize pls

/*
var (
	baseUrl = "http://tartu/photo/resized"
	suffix  = "png"
)

func AddResize(photos entity.Photos, resizePath string, sizes resize.Sizes) (err error) {

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
*/

/*
func ParseMeta(dirPath string) (photos []entity.Photo, err error) {

	jsonFiles := []string{}
	err = fs.WalkDir(os.DirFS(dirPath), ".", func(fn string, entry fs.DirEntry, err error) error {

		// Todo: walk is hairball, not likely to work on subdirs yet :/

		if err != nil {
			return err
		}

		if strings.HasSuffix(fn, "json") {
			// Todo: use ext
			jsonFiles = append(jsonFiles, fn)
		}

		return nil
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to walk")
		return
	}

	photos = []entity.Photo{}
	for _, jsonFile := range jsonFiles {
		photo, err := DecodeMeta(path.Join(dirPath, jsonFile))
		if err != nil {
			return err
		}

		err = os.Stat("/path/to/whatever")
		if err != nil {
			return err
		}
		if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
		}
		fmt.Printf(">>> %#v\n", photo)
		photos = append(photos, photo)
	}
	return
}
*/
