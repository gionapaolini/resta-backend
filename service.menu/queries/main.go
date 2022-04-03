package main

import (
	"log"
	"net/http"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
	"github.com/gorilla/mux"
)

func main() {
	utils.Time = clock.New()

	router := mux.NewRouter()
	// internal.SetupQueryApi(router, entityRepository)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
