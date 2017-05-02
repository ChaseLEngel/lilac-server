package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func logRequest(method string, uri string, name string) {
	log.Infof("%v %v %v\n", method, uri, name)
}

func routeSetup(handler http.Handler, name string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logRequest(r.Method, r.RequestURI, name)
			handler.ServeHTTP(w, r)
		},
	)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(routeSetup(handler, route.Name))
	}
	return router
}
