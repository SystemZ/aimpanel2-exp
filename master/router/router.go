package router

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	v1 := router.PathPrefix("/v1").Subrouter()

	for _, route := range routes {
		v1.HandleFunc(route.Pattern, route.HandlerFunc).Name(route.Name).Methods(route.Method)
	}

	return router
}
