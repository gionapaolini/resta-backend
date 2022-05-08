package internal

import (
	"errors"

	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type Api struct {
	repository eventutils.IEntityRepository
}

func SetupApi(app *fiber.App, repo eventutils.IEntityRepository) {
	api := Api{
		repository: repo,
	}
	api.setupRoutes(app)
}

func (api Api) setupRoutes(app *fiber.App) {
	app.Post("/menus", api.CreateNewMenu)
	app.Post("/menus/:id/enable", api.EnableMenu)
	app.Post("/menus/:id/disable", api.DisableMenu)
	app.Post("/menus/:id/change-name", api.ChangeMenuName)

	app.Post("/categories", api.CreateNewCategory)
	app.Post("/categories/:id/change-name", api.ChangeCategoryName)
}

func (api Api) CreateNewMenu(c *fiber.Ctx) error {
	menu := entities.NewMenu()
	err := api.repository.SaveEntity(menu)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the new menu. Please try again later")
	}
	c.SendStatus(fiber.StatusCreated)
	return nil
}

func (api Api) EnableMenu(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid menu id")
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	menu = menu.(entities.Menu).Enable()
	err = api.repository.SaveEntity(menu)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusOK)
	return nil
}

func (api Api) DisableMenu(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid menu id")
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	menu = menu.(entities.Menu).Disable()
	err = api.repository.SaveEntity(menu)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusOK)
	return nil
}

type ChangeMenuNameRequest struct {
	NewName string `json:"newName"`
}

func (api Api) ChangeMenuName(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid menu id")
	}

	reqBody := new(ChangeMenuNameRequest)
	if err := c.BodyParser(reqBody); err != nil {
		return err
	}

	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	menu = menu.(entities.Menu).ChangeName(reqBody.NewName)
	err = api.repository.SaveEntity(menu)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusOK)
	return nil
}

type CreateNewCategoryRequest struct {
	MenuID string `json:"menuID"`
}

func (api Api) CreateNewCategory(c *fiber.Ctx) error {
	reqBody := new(CreateNewCategoryRequest)
	if err := c.BodyParser(reqBody); err != nil {
		return err
	}
	menuID := uuid.FromStringOrNil(reqBody.MenuID)
	if menuID == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "menuID is not valid")
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), menuID)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	if menu.(entities.Menu).IsDeleted {
		return fiber.NewError(fiber.StatusNotFound, "Menu not found")
	}
	category := entities.NewCategory(menuID)
	err = api.repository.SaveEntity(category)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusCreated)
	return nil
}

type ChangeCategoryNameRequest struct {
	NewName string `json:"newName"`
}

func (api Api) ChangeCategoryName(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category id")
	}

	reqBody := new(ChangeMenuNameRequest)
	if err := c.BodyParser(reqBody); err != nil {
		return err
	}

	category, err := api.repository.GetEntity(entities.EmptyCategory(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Category not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the category, please try again later.")
		}
	}
	category = category.(entities.Category).ChangeName(reqBody.NewName)
	err = api.repository.SaveEntity(category)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusOK)
	return nil
}
