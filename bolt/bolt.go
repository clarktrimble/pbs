// Package bolt impletments the Repo interface with a bbolt backend.
package bolt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"

	"pbs/entity"
)

// Would be nice if this package had less bolt stuff and could focus on adapting
// a particular type.  Factoring json encode/decode would help isolate the
// bolty closures, but looks like we'll need to hack brackets and commas to
// data returned for Photos before decode, yay!
//
// Update: yes, see Upsert/Get Book below, still thorny for plurality (photos).
// Leaving this here for now as a reasonable stab at bolt.  Upstream interface will
// sheild others in any case.

const (
	photoBkt string = "photo"
	bookBkt  string = "book"
)

// Config is the bolt configuration.
type Config struct {
	Path string `json:"path" desc:"path to db file (inculsive)" required:"true"`
}

// Bolt represents a bbolt db.
type Bolt struct {
	db *bbolt.DB
}

// New creates an instance from config.
func (cfg *Config) New() (blt *Bolt, err error) {

	blt, err = New(cfg.Path)
	return
}

// New creates a Bolt instance and opens its db file.
func New(path string) (blt *Bolt, err error) {

	db, err := bbolt.Open(path, 0644, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		err = errors.Wrapf(err, fmt.Sprintf("failed to open db: %s", path))
		return
	}

	blt = &Bolt{
		db: db,
	}

	return
}

// Close closes the db.
func (blt *Bolt) Close() (err error) {

	err = blt.db.Close()
	err = errors.Wrapf(err, "failed to close bolt db")
	return
}

// UpsertPhotos writes photos overwriting those with same name.
func (blt *Bolt) UpsertPhotos(photos entity.Photos) (err error) {

	err = blt.db.Update(func(tx *bbolt.Tx) error {

		bkt, err := writeBucket(tx, photoBkt)
		if err != nil {
			return err
		}

		for _, photo := range photos {

			data, err := photo.Encode()
			if err != nil {
				return err
			}
			err = bkt.Put([]byte(photo.Name), data)
			if err != nil {
				err = errors.Wrapf(err, "failed to put photo")
				return err
			}
		}

		return nil
	})

	return
}

// GetPhotots gets all photos.
func (blt *Bolt) GetPhotos() (photos entity.Photos, err error) {

	err = blt.db.View(func(tx *bbolt.Tx) error {

		bkt, err := readBucket(tx, photoBkt)
		if err != nil {
			return err
		}

		cursor := bkt.Cursor()
		photos = entity.Photos{}

		for key, data := cursor.First(); key != nil; key, data = cursor.Next() {

			photo, err := entity.DecodePhoto(data)
			if err != nil {
				return err
			}

			photos = append(photos, photo)
		}

		return nil
	})

	return
}

// UpsertBook writes a book overwriting if same id.
func (blt *Bolt) UpsertBook(book entity.Book) (err error) {

	err = blt.upsertObj(book.Id, book)
	return
}

// GetBook gets a book.
func (blt *Bolt) GetBook(id string) (book entity.Book, err error) {

	book = entity.Book{}
	err = blt.getObj(id, &book)
	return
}

// unexported

func (blt *Bolt) upsertObj(id string, obj any) (err error) {

	err = blt.db.Update(func(tx *bbolt.Tx) error {

		bkt, err := writeBucket(tx, bookBkt)
		if err != nil {
			return err
		}

		data, err := json.Marshal(obj)
		if err != nil {
			err = errors.Wrapf(err, fmt.Sprintf("somehow failed to encode obj: %#v", obj))
			return err
		}
		err = bkt.Put([]byte(id), data)
		if err != nil {
			err = errors.Wrapf(err, "failed to put book")
			return err
		}

		return nil
	})

	return
}

func (blt *Bolt) getObj(id string, obj any) (err error) {

	err = blt.db.View(func(tx *bbolt.Tx) error {

		bkt, err := readBucket(tx, bookBkt)
		if err != nil {
			return err
		}

		data := bkt.Get([]byte(id))
		if data == nil {
			err = errors.Errorf("no data found for id: %s", id)
			return err
		}

		err = json.Unmarshal(data, obj)
		if err != nil {
			err = errors.Wrapf(err, fmt.Sprintf("failed to decode into: %#v", obj))
			return err
		}

		return nil
	})

	return
}

func readBucket(tx *bbolt.Tx, bucket string) (bkt *bbolt.Bucket, err error) {

	bkt = tx.Bucket([]byte(bucket))
	if bkt == nil {
		err = errors.Errorf(fmt.Sprintf("bucket: %s not found", bucket))
	}
	return
}

func writeBucket(tx *bbolt.Tx, bucket string) (bkt *bbolt.Bucket, err error) {

	// tx must be opened for write here, garh

	bkt, err = tx.CreateBucketIfNotExists([]byte(bucket))
	err = errors.Wrapf(err, "failed to get/create bucket")
	return
}
