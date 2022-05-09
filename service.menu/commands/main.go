package main

import (
	"os"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
)

const eventStoreConnectionString = "esdb://127.0.0.1:2113?tls=false&keepAliveTimeout=10000&keepAliveInterval=10000"

func main() {
	// TODO: FIX THIS IN A PROPER WAY
	os.Setenv("RESOURCE_FOLDER", "../../resources")
	err := os.MkdirAll("../../resources/images/categories", 0755)
	if err != nil {
		panic(err)
	}
	utils.Time = clock.New()

	settings, _ := esdb.ParseConnectionString(eventStoreConnectionString)
	db, _ := esdb.NewClient(settings)

	eventStore, err := eventutils.NewEventStore(eventStoreConnectionString)
	if err != nil {
		panic(err)
	}

	entityRepository := eventutils.NewEntityRepository(eventStore)

	eventHandler := eventutils.NewEventHandler(db, "menu.commands")
	menuEventHandler := internal.NewMenuEventHandler(entityRepository)
	eventHandler.HandleEvent("CategoryCreated", menuEventHandler.HandleCategoryCreated)
	eventHandler.Start()

	app := fiber.New()
	internal.SetupApi(app, entityRepository)

	app.Listen(":10000")
}
