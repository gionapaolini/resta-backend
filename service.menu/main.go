package main

import (
	"log"
	"net/http"

	"github.com/Resta-Inc/resta/menu/internal"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
	"github.com/gorilla/mux"
)

func main() {
	utils.Time = clock.New()
	eventStore, err := eventutils.NewEventStore("esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000")
	if err != nil {
		panic(err)
	}
	entityRepository := eventutils.NewEntityRepository(eventStore)
	api := internal.NewApi(entityRepository)

	r := mux.NewRouter()
	r.HandleFunc("/menu", api.CreateNewMenu).Methods("POST")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
