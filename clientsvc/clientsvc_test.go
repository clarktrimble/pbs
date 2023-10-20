package clientsvc_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "pbs/clientsvc"
	"pbs/entity"
	"pbs/test/mock"
)

func TestClientSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClientSvc Suite")
}

var _ = Describe("ClientSvc", func() {

	Describe("posting photos", func() {
		var (
			ctx    context.Context
			photos entity.Photos
			svc    *ClientSvc
			err    error
		)

		JustBeforeEach(func() {
			err = svc.PostPhotos(ctx, photos)
		})

		When("all is well", func() {
			BeforeEach(func() {
				ctx = context.Background()
				svc = &ClientSvc{
					Client: &mock.Client{},
				}
				photos = entity.Photos{{
					Id:   "GIehp1s",
					Name: "PXL-123",
				}}
			})

			It("should not return an error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
