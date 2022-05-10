package internal

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type Api struct {
	menuRepository IMenuRepository
	resourceHost   string
}

func SetupApi(app *fiber.App, repo IMenuRepository, resourcePath, resourceHost string) {
	api := Api{
		menuRepository: repo,
		resourceHost:   resourceHost,
	}
	api.setupRoutes(app, resourcePath)
}

func (api Api) setupRoutes(app *fiber.App, resourcePath string) {
	app.Get("/menus/:id", api.GetMenu)
	app.Get("/menus", api.GetAllMenus)

	app.Get("/categories/by-ids", api.GetCategoriesByIDs)

	app.Get("/subcategories/by-ids", api.GetSubCategoriesByIDs)

	path := filepath.Join(resourcePath, "images/categories")
	app.Static("/images/categories", path)
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

	categories = populateCategoryImageURL(categories, api)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the categories, please try again later.")
	}
	return c.JSON(categories)
}

func (api Api) GetSubCategoriesByIDs(c *fiber.Ctx) error {
	ids := c.Query("id")
	uuids := []uuid.UUID{}
	for _, v := range strings.Split(ids, ",") {
		parsedID := uuid.FromStringOrNil(v)
		if parsedID == uuid.Nil {
			fmtError := fmt.Sprintf("invalid subcategory id: %s", v)
			return fiber.NewError(fiber.StatusBadRequest, fmtError)
		}
		uuids = append(uuids, parsedID)
	}
	subcategories, err := api.menuRepository.GetSubCategoriesByIDs(uuids)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the subcategories, please try again later.")
	}
	return c.JSON(subcategories)
}

func populateCategoryImageURL(categories []CategoryView, api Api) (populatedCategories []CategoryView) {
	for _, v := range categories {
		v.ImageURL = fmt.Sprintf("%s/images/categories/%s.jpg", api.resourceHost, v.ID)
		populatedCategories = append(populatedCategories, v)
	}
	return
}
