package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/aimpanel2/master/middleware"
	"net/http"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	v1 := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc

		if route.AuthRequired {
			handler = middleware.AuthMiddleware(route.HandlerFunc)
		}

		v1.Path(route.Pattern).Handler(handler).Name(route.Name).Methods(route.Method)
	}

	return router
}
