package entity_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"pbs/entity"
	. "pbs/entity"
)

func TestEntity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entity Suite")
}

var _ = Describe("Entity PhotoBook", func() {

	Describe("handling a get photos request", func() {
		var (
			photos Photos
			book   Book
			pb     PhotoBook
		)

		JustBeforeEach(func() {
			pb = New(photos, book)
		})

		When("all is well", func() {
			BeforeEach(func() {
				photos = Photos{{
					Id:   "Vrf8tZP",
					Name: "PXL_20230728_234642453",
					Geo: Geo{
						Lat: 36.0131361,
						Lon: -95.8911583,
					},
					Images: map[string]Image{
						"large":    {Width: 768, Height: 1020, Path: "http://tartu/photo/resized/PXL_20230728_234642453-large.png"},
						"thumb":    {Path: "http://tartu/photo/resized/PXL_20230728_234642453-thumb.png"},
						"thumb-gs": {Path: "http://tartu/photo/resized/PXL_20230728_234642453-thumb-gs.png"},
					},
				}}
				book = Book{
					Id:       "testo",
					Featured: map[string]bool{"Vrf8tZP": true},
				}
			})

			It("should respond with photos", func() {
				Expect(pb).To(Equal(entity.PhotoBook{{
					PhotoId:  "Vrf8tZP",
					Source:   "http://tartu/photo/resized/PXL_20230728_234642453-large.png",
					Width:    768,
					Height:   1020,
					Thumb:    "http://tartu/photo/resized/PXL_20230728_234642453-thumb.png",
					ThumbGs:  "http://tartu/photo/resized/PXL_20230728_234642453-thumb-gs.png",
					Lat:      36.0131361,
					Lon:      -95.8911583,
					TakenAt:  time.Time{},
					Featured: true,
				}}))
			})
		})
	})
})
