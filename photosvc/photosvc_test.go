package photosvc

// note: not using "_test" pkg for access to handlers
// Todo: beef up repo mock to check upsert calls and params

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"xform/test/mock"
)

func TestPhotoSvc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PhotoSvc Suite")
}

var _ = Describe("PhotoSvc", func() {

	var (
		writer  *httptest.ResponseRecorder
		request *http.Request
		lgr     *mock.Logger
		svc     *PhotoSvc
		err     error
	)

	BeforeEach(func() {
		writer = httptest.NewRecorder()

		lgr = mock.NewLogger()
		svc = &PhotoSvc{
			Logger: lgr,
			Repo:   &mock.Repo{},
		}
	})

	Describe("handling a get photos request", func() {

		JustBeforeEach(func() {
			svc.getPhotos(writer, request)
		})

		When("all is well", func() {
			BeforeEach(func() {
				request, err = http.NewRequest("GET", "example.com", &bytes.Buffer{})
				Expect(err).ToNot(HaveOccurred())
			})

			It("should respond with photos", func() {
				serialized := `{"photos":[{"Id":"asdf","Name":"PXL-123","TakenAt":"0001-01-01T00:00:00Z",`
				serialized += `"Geo":{"Lat":0,"Lon":0,"Alt":0},"Images":null}]}`

				Expect(writer.Code).To(Equal(200))
				Expect(writer.Body.String()).To(Equal(serialized))
			})
		})
	})

	Describe("handling an upsert photos request", func() {
		JustBeforeEach(func() {
			svc.upsertPhotos(writer, request)
		})

		When("all is well", func() {
			BeforeEach(func() {
				buf := bytes.NewBufferString(`[{"id": "asdf", "name": "PXL-123"}]`)
				request, err = http.NewRequest("POST", "example.com", buf)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should respond with OK", func() {
				Expect(writer.Code).To(Equal(200))
				Expect(writer.Body.String()).To(Equal(`{"status":"ok"}`))
			})
		})
	})

	Describe("handling an upsert book request", func() {
		JustBeforeEach(func() {
			svc.upsertBook(writer, request)
		})

		When("all is well", func() {
			BeforeEach(func() {
				buf := bytes.NewBufferString(`{"id": "book001"}`)
				request, err = http.NewRequest("POST", "example.com", buf)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should respond with OK", func() {
				Expect(writer.Code).To(Equal(200))
				Expect(writer.Body.String()).To(Equal(`{"status":"ok"}`))
			})
		})
	})

	Describe("handling a set featured request", func() {
		JustBeforeEach(func() {
			svc.setFeatured(writer, request)
		})

		When("all is well", func() {
			BeforeEach(func() {
				buf := bytes.NewBufferString(`{"book_id": "book001", "photo_id": "asdf", "featured": true}`)
				request, err = http.NewRequest("POST", "example.com", buf)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should respond with OK", func() {
				Expect(writer.Code).To(Equal(200))
				Expect(writer.Body.String()).To(Equal(`{"status":"ok"}`))
			})
		})
	})

	Describe("handling a get photobook request", func() {
		JustBeforeEach(func() {
			svc.getPhotoBook(writer, request)
		})

		When("all is well", func() {
			BeforeEach(func() {
				request, err = http.NewRequest("POST", "example.com", &bytes.Buffer{})
				Expect(err).ToNot(HaveOccurred())

				ctx := chi.NewRouteContext()
				ctx.URLParams.Add("bookId", "book001")
				request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, ctx))
			})

			It("should respond with OK", func() {
				serial := `{"images":[{"photo_id":"asdf","src":"","width":0,"height":0,"thumb":"","thumb_gs":"",`
				serial += `"lat":0,"lon":0,"taken_at":"0001-01-01T00:00:00Z","featured":true}]}`
				Expect(writer.Code).To(Equal(200))
				Expect(writer.Body.String()).To(Equal(serial))
			})
		})
	})
})
