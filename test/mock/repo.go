package mock

import (
	"fmt"
	"xform/entity"
)

type Repo struct{}

func (repo *Repo) GetPhotos() (photos entity.Photos, err error) {

	photos = entity.Photos{
		{
			Id:   "asdf",
			Name: "PXL-123",
		},
	}

	return
}

func (repo *Repo) UpsertPhotos(photos entity.Photos) (err error) {

	// Todo: check photos

	return
}

func (repo *Repo) GetBook(bookId string) (book entity.Book, err error) {

	if bookId != "book001" {

		err = fmt.Errorf("mock does not like bookId: %s", bookId)
		return
	}

	book = entity.Book{
		Id:       "book001",
		Featured: map[string]bool{"asdf": true},
	}

	return
}

func (repo *Repo) UpsertBook(book entity.Book) (err error) {

	// Todo: check book

	return
}
