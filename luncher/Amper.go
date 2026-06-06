package main

import (
	"amper/cache/business"
	"amper/cache/database"
	"amper/common/util/ampstrings"
	"amper/controller"
	"amper/properties/application"
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	amperId := business.AmperId()
	if amperId == nil {
		log.Print(fmt.Errorf("identifier is a required property, make sure it is configured, and later never changed"))
		return
	}
	port := ""
	config, errAP := application.Get()
	if errAP == nil {
		port = config.GetString("amper.port", "")
	} else {
		log.Print(fmt.Errorf("make sure application.properties file exists in luncher directory"))
	}

	if !ampstrings.HasValue(&port) {
		port = "7777"
	}
	database.Init()
	//http.HandleFunc("/", indexHandler)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		indexHandler(w, r)
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.AllowAll().Handler(mux)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	controller.Dispatch(&w, r)
}
