package photosvc

import (
	"net/http"
)

type Router interface {
	Set(method, path string, handler http.HandlerFunc) (err error)
}
