package bolt_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pbs/bolt"
	"pbs/entity"
)

func TestBolt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bolt Suite")
}

var _ = Describe("Bolt", func() {

	var (
		err error
		blt *Bolt
	)

	BeforeEach(func() {

		file, err := os.CreateTemp("", "test.*.db")
		Expect(err).ToNot(HaveOccurred())
		filename := file.Name()
		file.Close()

		blt, err = New(filename)
		Expect(err).ToNot(HaveOccurred())

		DeferCleanup(func() {
			err = blt.Close()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("upserting and getting book", func() {
		var (
			book    entity.Book
			bookToo entity.Book
			bookId  string
		)

		JustBeforeEach(func() {
			err = blt.UpsertBook(book)
			Expect(err).ToNot(HaveOccurred())
			bookToo, err = blt.GetBook(bookId)
		})

		BeforeEach(func() {
			book = entity.Book{
				Id:       "asdf",
				Featured: map[string]bool{"123": true},
			}
			bookId = "asdf"
		})

		When("all is well", func() {
			It("make a round trip", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(bookToo).To(Equal(book))
			})
		})

		When("book is not found", func() {
			BeforeEach(func() {
				bookId = "grrk"
			})

			It("is not found error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("no data found for id: grrk"))
			})
		})
	})

	Describe("upserting and getting photos", func() {
		var (
			photos    entity.Photos
			photosToo entity.Photos
		)

		JustBeforeEach(func() {
			err = blt.UpsertPhotos(photos)
			Expect(err).ToNot(HaveOccurred())
			photosToo, err = blt.GetPhotos()
		})

		When("all is well", func() {
			BeforeEach(func() {
				photos = entity.Photos{{
					Id:   "asdf",
					Name: "PXL-123",
				}}
			})

			It("make a round trip", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(photosToo).To(Equal(photos))
			})
		})
	})
})
