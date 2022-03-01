package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte("pong"))
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var name string
	if name = ps.ByName("name"); name == "" {
		name = "gophers"
	}
	fmt.Fprintf(w, "Hello, %s! Welcome to sample web app for httpmiddleware demo...\n", name)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/ping", ping)

	log.Fatal(http.ListenAndServe(":8080", router))
}
