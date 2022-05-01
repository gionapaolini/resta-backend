package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/queries/internal"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

var pgConnectionString string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

const eventStoreConnectionString = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"

func main() {
	utils.Time = clock.New()
	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)
	eventHandler := eventutils.NewEventHandler(db, "menu.queries")
	menuRepository := internal.NewMenuRepository(pgConnectionString)
	menuEventHandler := internal.NewMenuEventHandler(menuRepository)
	eventHandler.HandleEvent("MenuCreated", menuEventHandler.HandleMenuCreated)
	eventHandler.HandleEvent("MenuEnabled", menuEventHandler.HandleMenuEnabled)
	eventHandler.HandleEvent("MenuDisabled", menuEventHandler.HandleMenuDisabled)
	eventHandler.HandleEvent("MenuNameChanged", menuEventHandler.HandleMenuNameChanged)
	eventHandler.HandleEvent("CategoryCreated", menuEventHandler.HandleCategoryCreated)
	eventHandler.HandleEvent("CategoryAddedToMenu", menuEventHandler.HandleCategoryAddedToMenu)
	eventHandler.Start()
	router := mux.NewRouter()
	internal.SetupApi(router, menuRepository)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":10001", nil))
}
