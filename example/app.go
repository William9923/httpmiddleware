package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/William9923/httpmiddleware"
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

func Logging(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logger := log.Default()
		logger.Println("start http request...")
		h(w, r, ps)
	}
}

func Authentication(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

		ctx := req.Context()

		req = req.WithContext(ctx)
		fn(w, req, p)
	}
}

func main() {

	middlewares := httpmiddleware.New()
	middlewares.Use(Logging)

	router := httprouter.New()
	router.GET("/", middlewares.Wrap(Index))
	router.GET("/ping", ping)

	log.Fatal(http.ListenAndServe(":8080", router))
}
