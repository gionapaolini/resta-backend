package internal

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type Api struct {
	menuRepository IMenuRepository
	Router         *mux.Router
}

func SetupApi(app *fiber.App, repo IMenuRepository) {
	api := Api{
		menuRepository: repo,
	}
	api.setupRoutes(app)
}

func (api Api) setupRoutes(app *fiber.App) {
	app.Get("/menus/:id", api.GetMenu)
	app.Get("/menus", api.GetAllMenus)
	app.Get("/categories/by-ids", api.GetCategoriesByIDs)
}

func (api Api) GetMenu(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid menu id")
	}
	menu, err := api.menuRepository.GetMenu(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	return c.JSON(menu)
}

func (api Api) GetAllMenus(c *fiber.Ctx) error {
	menus, err := api.menuRepository.GetAllMenus()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menus, please try again later.")
	}
	return c.JSON(menus)
}

func (api Api) GetCategoriesByIDs(c *fiber.Ctx) error {
	ids := c.Query("id")
	uuids := []uuid.UUID{}
	for _, v := range strings.Split(ids, ",") {
		parsedID := uuid.FromStringOrNil(v)
		if parsedID == uuid.Nil {
			fmtError := fmt.Sprintf("invalid category id: %s", v)
			return fiber.NewError(fiber.StatusBadRequest, fmtError)
		}
		uuids = append(uuids, parsedID)
	}
	categories, err := api.menuRepository.GetCategoriesByIDs(uuids)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the categories, please try again later.")
	}
	return c.JSON(categories)
}
