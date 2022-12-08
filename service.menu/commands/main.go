package main

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/commands/internal"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := internal.LoadConfig(".")
	utils.Time = clock.New()

	settings, _ := esdb.ParseConnectionString(config.EventStoreConnectionString)
	db, _ := esdb.NewClient(settings)

	eventStore, err := eventutils2.NewEventStore(db)
	if err != nil {
		panic(err)
	}

	entityRepository := eventutils2.NewEntityRepository(eventStore)

	eventHandler := eventutils.NewEventHandler(db, "menu.commands")
	menuEventHandler := internal.NewMenuEventHandler(entityRepository)
	eventHandler.HandleEvent("CategoryCreated", menuEventHandler.HandleCategoryCreated)
	eventHandler.HandleEvent("SubCategoryCreated", menuEventHandler.HandleSubCategoryCreated)
	eventHandler.HandleEvent("MenuItemCreated", menuEventHandler.HandleMenuItemCreated)
	eventHandler.Start()

	app := fiber.New()
	internal.SetupApi(app, entityRepository, config.ResourcePath)

	app.Listen(":10000")
}
