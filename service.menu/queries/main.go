package main

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/Resta-Inc/resta/menu/queries/internal"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config := internal.LoadConfig(".")
	utils.Time = clock.New()
	settings, _ := esdb.ParseConnectionString(config.EventStoreConnectionString)
	db, _ := esdb.NewClient(settings)
	eventHandler := eventutils.NewEventHandler(db, "menu.queries")
	menuRepository := internal.NewMenuRepository(config.PostgresConnectionString)
	menuEventHandler := internal.NewMenuEventHandler(menuRepository)
	eventHandler.HandleEvent("MenuCreated", menuEventHandler.HandleMenuCreated)
	eventHandler.HandleEvent("MenuEnabled", menuEventHandler.HandleMenuEnabled)
	eventHandler.HandleEvent("MenuDisabled", menuEventHandler.HandleMenuDisabled)
	eventHandler.HandleEvent("MenuNameChanged", menuEventHandler.HandleMenuNameChanged)
	eventHandler.HandleEvent("CategoryCreated", menuEventHandler.HandleCategoryCreated)
	eventHandler.HandleEvent("CategoryAddedToMenu", menuEventHandler.HandleCategoryAddedToMenu)
	eventHandler.HandleEvent("CategoryNameChanged", menuEventHandler.HandleCategoryNameChanged)
	eventHandler.Start()

	app := fiber.New()
	internal.SetupApi(app, menuRepository)

	app.Listen(":10001")
}
