// Package takeout soaks up from goog's takeout.
package takeout

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"xform/entity"

	"github.com/clarktrimble/hondo"
	"github.com/pkg/errors"
)

// Geo represents a location.
type Geo struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
	Alt float64 `json:"altitude"`
}

// Take is when snap was snapped.
type Taken struct {
	Epoch string `json:"timestamp"`
}

// Photo is as seen in Takeout json files.
type Photo struct {
	Title string `json:"title"`
	Taken Taken  `json:"photoTakenTime"`
	Geo   Geo    `json:"geoData"`
}

// FromFiles creates photos given ...
func FromFiles(jsonPath string, filter string) (photos entity.Photos, err error) {

	// Todo: yeah would be better to have different logics for path and name (for resize)

	paths, err := findJson(jsonPath, filter)
	if err != nil {
		return
	}

	photos = entity.Photos{}
	for _, pth := range paths {

		var to Photo
		dir := path.Join(jsonPath, path.Dir(pth))
		file := path.Base(pth)

		to, err = decode(dir, file)
		if err != nil {
			return
		}

		var photo entity.Photo
		photo, err = takeoutToEntity(to, dir)
		if err != nil {
			return
		}

		if skip(photo.Path) {
			continue
		}

		photos = append(photos, photo)
	}

	return
}

// unexported

func skip(path string) (skip bool) {

	// Todo: log rather than print herein

	if !strings.HasSuffix(path, ".jpg") {
		fmt.Printf("skipping %s\n", path)
		skip = true
	}

	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("also skipping %s\n", path)
		skip = true
	}

	return
}

func findJson(root, filter string) (paths []string, err error) {

	// Todo: maybe walk is overkill here??

	paths = []string{}
	err = fs.WalkDir(os.DirFS(root), ".", func(path string, entry fs.DirEntry, inerr error) error {

		// Todo: factor this and name/path for resize plz
		// Todo: check error
		re, err := regexp.Compile(filter)
		if err != nil {
			err = errors.Wrapf(err, "failed to compile regex")
			return err
		}
		if !re.MatchString(path) {
			return nil
		}

		if inerr != nil {
			inerr = errors.Wrapf(inerr, "walker called with error")
			return inerr
		}

		if strings.HasSuffix(path, ".json") {
			paths = append(paths, path)
		}

		return nil
	})
	err = errors.Wrapf(err, "failed to walk")
	return
}

func decode(dir, name string) (to Photo, err error) {

	file, err := os.Open(path.Join(dir, name))
	if err != nil {
		err = errors.Wrapf(err, "failed to open")
		return
	}
	defer file.Close()

	to = Photo{}
	err = json.NewDecoder(file).Decode(&to)
	err = errors.Wrapf(err, "failed to decode")
	return
}

func takeoutToEntity(to Photo, dir string) (photo entity.Photo, err error) {

	epoch, err := strconv.ParseInt(to.Taken.Epoch, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse epoch")
		return
	}

	split := strings.Split(to.Title, ".")

	photo = entity.Photo{
		Id:      hondo.Rand(7),
		Name:    split[0],
		Path:    path.Join(dir, to.Title),
		TakenAt: time.Unix(epoch, 0).UTC(),
		Geo: entity.Geo{
			Lat: to.Geo.Lat,
			Lon: to.Geo.Lon,
			Alt: to.Geo.Alt,
		},
	}

	return
}
