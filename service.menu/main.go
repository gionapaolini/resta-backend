package main

import (
	"log"
	"net/http"

	"github.com/Resta-Inc/resta/menu/internal"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
)

func main() {
	utils.Time = clock.New()

	eventStore, err := eventutils.NewEventStore("esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000")
	if err != nil {
		panic(err)
	}

	entityRepository := eventutils.NewEntityRepository(eventStore)

	api := internal.NewCommandsApi(entityRepository)

	http.Handle("/", api.Router)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
