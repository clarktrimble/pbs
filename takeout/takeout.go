// Package takeout soaks up from goog's takeout
package takeout

import (
	"encoding/json"
	"io/fs"
	"os"
	pth "path"
	"strconv"
	"strings"
	"time"

	"xform/entity"

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

type Photo struct {
	Title string `json:"title"`
	Taken Taken  `json:"photoTakenTime"`
	Geo   Geo    `json:"geoData"`
}

func Find(root string) (photos []entity.Photo, err error) {

	paths, err := findJson(root)
	if err != nil {
		return
	}

	for _, path := range paths {

		var photo entity.Photo
		photo, err = decode(pth.Join(root, pth.Dir(path)), pth.Base(path))
		if err != nil {
			return
		}

		if !strings.HasSuffix(photo.Path, ".jpg") {
			continue
		}

		_, err = os.Stat(photo.Path)
		if err != nil {
			continue
		}
		photos = append(photos, photo)
	}

	return
}

func findJson(root string) (paths []string, err error) {

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

	photo = entity.Photo{
		Name:    split[0],
		Path:    pth.Join(dir, to.Title),
		TakenAt: time.Unix(epoch, 0).UTC(),
		Geo: entity.Geo{
			Lat: to.Geo.Lat,
			Lon: to.Geo.Lon,
			Alt: to.Geo.Alt,
		},
	}

	return
}

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
