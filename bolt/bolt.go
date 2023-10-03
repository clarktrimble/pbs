// Package bolt impletments the Repo interface with a bbolt backend.
package bolt

import (
	"context"
	"fmt"
	"time"
	"xform/entity"

	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

// Bolt represents a bbolt db.
type Bolt struct {
	db     *bbolt.DB
	bucket []byte
}

// New creates a Bolt instance and opens its db file.
func New(path, bucket string) (blt *Bolt, err error) {

	db, err := bbolt.Open(path, 0644, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		err = errors.Wrapf(err, fmt.Sprintf("failed to open db: %s", path))
		return
	}

	blt = &Bolt{
		db:     db,
		bucket: []byte(bucket),
	}

	return
}

func (blt *Bolt) Close() {

	err := blt.db.Close()
	if err != nil {
		panic(err) // Todo: log??
	}
}

func (blt *Bolt) UpsertPhotos(ctx context.Context, photos entity.Photos) (err error) {

	err = blt.db.Update(func(tx *bbolt.Tx) error {

		bkt, err := writeBucket(tx, "photo")
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

func (blt *Bolt) GetPhotos(ctx context.Context) (photos entity.Photos, err error) {

	err = blt.db.View(func(tx *bbolt.Tx) error {

		bkt, err := readBucket(tx, "photo")
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

func (blt *Bolt) UpsertBook(ctx context.Context, book entity.Book) (err error) {

	err = blt.db.Update(func(tx *bbolt.Tx) error {

		bkt, err := writeBucket(tx, "book")
		if err != nil {
			return err
		}

		data, err := book.Encode()
		if err != nil {
			return err
		}
		err = bkt.Put([]byte(book.Id), data)
		if err != nil {
			err = errors.Wrapf(err, "failed to put book")
			return err
		}

		return nil
	})

	return
}

func (blt *Bolt) GetBook(ctx context.Context, id string) (book entity.Book, err error) {

	err = blt.db.View(func(tx *bbolt.Tx) error {

		bkt, err := readBucket(tx, "book")
		if err != nil {
			return err
		}

		data := bkt.Get([]byte(id))
		//fmt.Printf(">>>%s<<< for %s\n\n\n", data, id)
		// Todo: not found??
		book, err = entity.DecodeBook(data)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

// unexported

func readBucket(tx *bbolt.Tx, bucket string) (bkt *bbolt.Bucket, err error) {

	bkt = tx.Bucket([]byte(bucket))
	if bkt == nil {
		err = errors.Errorf(fmt.Sprintf("bucket: %s not found", bucket))
	}
	return
}

func writeBucket(tx *bbolt.Tx, bucket string) (bkt *bbolt.Bucket, err error) {

	bkt, err = tx.CreateBucketIfNotExists([]byte(bucket))
	err = errors.Wrapf(err, "failed to get/create bucket")
	return
}
