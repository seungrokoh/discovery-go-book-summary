package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


const (
	apiPathPrefix  = "/api/v1/task/"
	htmlPathPrefix = "/task/"
	idPattern      = "/{id:[0-9]+}"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix(htmlPathPrefix).
		Path(idPattern).
		Methods("GET").HandlerFunc(htmlHandler)

	s := r.PathPrefix(apiPathPrefix).Subrouter()
	s.HandleFunc(idPattern, apiGetHandler).Methods("GET")
	s.HandleFunc(idPattern, apiPutHandler).Methods("PUT")
	s.HandleFunc("/", apiPostHandler).Methods("POST")
	s.HandleFunc(idPattern, apiDeleteHandler).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8887", nil))
}
