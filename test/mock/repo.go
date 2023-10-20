package mock

import (
	"fmt"
	"pbs/entity"
)

// hand tooled mock is a little thin for repo
// gomock or moq from this point perhaps

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

	if len(photos) != 1 || photos[0].Id != "asdf" {
		err = fmt.Errorf("mock does not like photos: %#v", photos)
		return
	}

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

	if book.Id != "book001" {
		err = fmt.Errorf("mock does not like book: %#v", book)
		return
	}

	return
}
