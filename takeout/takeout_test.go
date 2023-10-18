package takeout_test

import (
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"xform/entity"
	. "xform/takeout"
)

func TestTakeout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Takeout Suite")
}

var _ = Describe("Takeout", func() {

	var (
		tos     Takeouts
		err     error
		root    string
		pattern string
	)

	Describe("getting data from filesystem", func() {
		BeforeEach(func() {
			root = "../test/data"
			pattern = "PXL_20230709_185504673"

			tos, err = ScanTakeout(root, pattern)
			Expect(err).ToNot(HaveOccurred())

			rand.Seed(1) //nolint:staticcheck // just for unit
		})

		When("all is well", func() {
			It("should return Takeouts", func() {
				Expect(tos).To(Equal(Takeouts{
					{
						Title:     "PXL_20230709_185504673.jpg",
						ImagePath: "../test/data/PXL_20230709_185504673.jpg",
						Taken: Taken{
							Epoch: "1688928904",
						},
						Geo: Geo{
							Lat: 57.29834439999999,
							Lon: 24.4072528,
							Alt: 17.49,
						},
					},
				}))
			})
		})

		Describe("converting takeouts to Photos", func() {
			var (
				photos entity.Photos
			)

			JustBeforeEach(func() {
				photos, err = tos.Photos()
			})

			When("all is well", func() {
				It("should return a corresponding Photo", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(photos).To(HaveLen(1))

					ts, err := time.Parse(time.RFC3339, "2023-07-09T18:55:04Z")
					Expect(err).ToNot(HaveOccurred())

					Expect(photos[0]).To(Equal(entity.Photo{
						Id:      "GIehp1s",
						Name:    "PXL_20230709_185504673",
						TakenAt: ts,
						Geo: entity.Geo{
							Lat: 57.29834439999999,
							Lon: 24.4072528,
							Alt: 17.49,
						},
					}))
				})
			})
		})

		Describe("converting Takeouts to PhotoFiles", func() {
			var (
				pfs []entity.PhotoFile
			)

			JustBeforeEach(func() {
				pfs = tos.PhotoFiles()
			})

			When("all is well", func() {
				It("should return a corresponding PhotoFile", func() {
					Expect(pfs).To(HaveLen(1))
					Expect(pfs[0]).To(Equal(entity.PhotoFile{
						Name: "PXL_20230709_185504673",
						Path: "../test/data/PXL_20230709_185504673.jpg",
					}))
				})
			})
		})

	})
})
