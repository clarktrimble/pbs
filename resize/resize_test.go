package resize_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pbs/entity"
	. "pbs/resize"
)

func TestResize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resize Suite")
}

var _ = Describe("Resize", func() {

	Describe("resizing  photos", func() {
		var (
			sizes Sizes
			files []entity.PhotoFile
			dst   string
			err   error
		)

		BeforeEach(func() {
			sizes = Sizes{
				{Name: "thumb", Scale: 16},
				{Name: "thumb-gs", Scale: 16, Gs: true},
			}
			files = []entity.PhotoFile{
				{
					Name: "book",
					Path: "../test/data/book.jpg",
				},
			}
			dst, err = os.MkdirTemp("/tmp", "pbtest")
			Expect(err).ToNot(HaveOccurred())
		})

		BeforeEach(func() {
			err = sizes.ResizePhotos(dst, files)
		})

		When("all is well", func() {
			It("should put files in destination", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(fmt.Sprintf("%s/book-%s.png", dst, sizes[0].Name)).To(BeARegularFile())
				Expect(fmt.Sprintf("%s/book-%s.png", dst, sizes[1].Name)).To(BeARegularFile())
			})
		})

		// nested test "adding resized" is convienient in that it can depend on the
		// above test having completed, but might get wierd if errrors were tested above
		Describe("adding resized image data to photos", func() {
			var (
				photos entity.Photos
			)

			JustBeforeEach(func() {
				err = sizes.AddImages(photos, dst)
			})

			When("all is well", func() {
				BeforeEach(func() {
					photos = entity.Photos{
						{
							Id:   "GIehp1s",
							Name: "book",
						},
					}
				})

				It("should populate Images", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(photos[0].Images).To(Equal(map[string]entity.Image{
						"thumb": entity.Image{
							SizeName: "thumb",
							Width:    38,
							Height:   25,
							Path:     "book-thumb.png",
						},
						"thumb-gs": entity.Image{
							SizeName: "thumb-gs",
							Width:    38,
							Height:   25,
							Path:     "book-thumb-gs.png",
						},
					}))
				})
			})
		})
	})
})
