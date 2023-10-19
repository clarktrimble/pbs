package resize_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"xform/entity"
	"xform/resize"
	. "xform/resize"
)

func TestResize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Resize Suite")
}

var _ = Describe("Resize", func() {

	var (
		sizes  Sizes
		files  []entity.PhotoFile
		photos entity.Photos
		dst    string
		err    error
	)

	Describe("resizing  photos", func() {

		BeforeEach(func() {
			sizes = resize.Sizes{
				{Name: "thumb", Scale: 16},
				{Name: "thumb-gs", Scale: 16, Gs: true},
			}
			dst, err = os.MkdirTemp("/tmp", "pbtest")
			Expect(err).ToNot(HaveOccurred())
			files = []entity.PhotoFile{
				{
					Name: "book",
					Path: "../test/data/book.jpg",
				},
			}

			photos = entity.Photos{
				{
					Id:   "GIehp1s",
					Name: "book",
				},
			}
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

		Describe("adding resized image data to photos", func() {
			BeforeEach(func() {
				err = AddImages(photos, dst, sizes)
			})
			When("all is well", func() {
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
					},
					))
				})
			})
		})
	})
})
