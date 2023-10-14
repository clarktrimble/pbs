// Package takeout soaks up from goog's takeout.
package takeout

// Todo: test on walkable dir or ditch walk plz
//       so look at takeout struct from gz again
// Todo: promote loadTakeout to public and skip there

import (
	"encoding/json"
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

// Taken is when snap was snapped.
type Taken struct {
	Epoch string `json:"timestamp"`
}

// Takeout represents a photo as seen in Takeout json files.
type Takeout struct {
	Title     string `json:"title"`
	Taken     Taken  `json:"photoTakenTime"`
	Geo       Geo    `json:"geoData"`
	ImagePath string `json:"-"`
}

// Takeouts are a multiplicity of Takeouts.
type Takeouts []Takeout

// ScanTakeout scans a takeout folder for metdatums.
func ScanTakeout(root, pattern string) (tos Takeouts, err error) {

	jsonExt := ".json"
	jpgExt := ".jpg"
	tos = Takeouts{}

	jsonFiles, err := findFiles(root, jsonExt, pattern)
	if err != nil {
		return
	}

	for _, jsonFile := range jsonFiles {

		var to Takeout
		to, err = decodeFile(path.Join(root, jsonFile))
		if err != nil {
			return
		}

		if !validImage(to.ImagePath, jpgExt) {
			// Todo: error
			continue
		}

		tos = append(tos, to)
	}

	return
}

// PhotoFiles converts Takeouts to PhotoFiles.
func (tos Takeouts) PhotoFiles() (pfs []entity.PhotoFile) {

	pfs = []entity.PhotoFile{}
	for _, to := range tos {

		pfs = append(pfs, entity.PhotoFile{
			Name: strings.TrimSuffix(to.Title, path.Ext(to.Title)),
			Path: to.ImagePath,
		})
	}

	return
}

// Photos converts Takeouts to Photos.
func (tos Takeouts) Photos() (photos entity.Photos, err error) {

	photos = entity.Photos{}
	for _, to := range tos {

		var photo entity.Photo
		photo, err = takeoutToEntity(to)
		if err != nil {
			return
		}

		photos = append(photos, photo)
	}

	return
}

// unexported

func validImage(pth, ext string) (valid bool) {

	valid = true

	_, err := os.Stat(pth)
	if err != nil {
		valid = false
	}

	return
}

func findFiles(root, ext, pattern string) (paths []string, err error) {

	paths = []string{}

	re, err := regexp.Compile(pattern)
	if err != nil {
		err = errors.Wrapf(err, "failed to compile regex")
		return
	}

	err = fs.WalkDir(os.DirFS(root), ".", func(pth string, entry fs.DirEntry, inerr error) error {

		//fmt.Printf(">>> walk: %s %s %#v\n", root, pth, entry)
		if inerr != nil {
			inerr = errors.Wrapf(inerr, "walker called with error")
			return inerr
		}

		if path.Ext(pth) != ext || !re.MatchString(pth) {
			return nil
		}

		paths = append(paths, pth)
		return nil
	})
	err = errors.Wrapf(err, "failed to walk")
	return
}

func decodeFile(name string) (to Takeout, err error) {

	file, err := os.Open(name)
	if err != nil {
		err = errors.Wrapf(err, "failed to open")
		return
	}
	defer file.Close()

	to = Takeout{}
	err = json.NewDecoder(file).Decode(&to)
	err = errors.Wrapf(err, "failed to decode")
	if err != nil {
		err = errors.Wrapf(err, "failed to decode")
		return
	}

	// takeout convention is image and json files have almost the same path
	to.ImagePath = strings.TrimSuffix(name, path.Ext(name))
	return
}

// func takeoutToEntity(to Photo, dir string) (photo entity.Photo, err error) {
func takeoutToEntity(to Takeout) (photo entity.Photo, err error) {

	epoch, err := strconv.ParseInt(to.Taken.Epoch, 10, 64)
	if err != nil {
		err = errors.Wrapf(err, "failed to parse epoch")
		return
	}

	split := strings.Split(to.Title, ".")

	photo = entity.Photo{
		Id:   hondo.Rand(7), // Todo: remove rand seed in hondo for testing plz
		Name: split[0],
		//Path:    path.Join(dir, to.Title),
		TakenAt: time.Unix(epoch, 0).UTC(),
		Geo: entity.Geo{
			Lat: to.Geo.Lat,
			Lon: to.Geo.Lon,
			Alt: to.Geo.Alt,
		},
	}

	return
}
