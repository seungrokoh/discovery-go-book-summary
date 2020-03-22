package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc(htmlPrefix, htmlHandler)
	http.HandleFunc(pathPrefix, apiHandler)
	log.Fatal(http.ListenAndServe(":8887", nil))
}
