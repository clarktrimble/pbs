package takeout_test

import (
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

	// Todo: test ScanTakeout on it's own an synthdata here
	// Todo: test with errors

	var (
		tos     Takeouts
		err     error
		root    string
		pattern string
	)

	BeforeEach(func() {
		root = "../test/data"
		pattern = "PXL_20230709_185504673"
		tos, err = ScanTakeout(root, pattern)
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("getting data from filesystem", func() {
		var (
			photos entity.Photos
		)

		JustBeforeEach(func() {
			photos, err = tos.Photos()
		})

		When("nothing in ctx", func() {
			It("should return an empty slice", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(photos).To(HaveLen(1))

				// Todo: remove hax
				photos[0].Id = "RVvhSQc"
				photos[0].TakenAt = time.Time{}

				Expect(photos[0]).To(Equal(entity.Photo{
					Id:   "RVvhSQc",
					Name: "PXL_20230709_185504673",
					Geo: entity.Geo{
						Lat: 57.29834439999999,
						Lon: 24.4072528,
						Alt: 17.49,
					},
				}))
			})
		})
	})

	Describe("getting data from filesystem", func() {
		var (
			pfs []entity.PhotoFile
		)

		JustBeforeEach(func() {
			pfs = tos.PhotoFiles()
		})

		When("nothing in ctx", func() {
			It("should return an empty slice", func() {
				Expect(pfs).To(HaveLen(1))
				Expect(pfs[0]).To(Equal(entity.PhotoFile{
					Name: "PXL_20230709_185504673",
					Path: "../test/data/PXL_20230709_185504673.jpg",
				}))
			})
		})
	})
})
