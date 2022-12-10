package internal

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type Api struct {
	repository   eventutils.IEntityRepository
	resourcePath string
}

func SetupApi(app *fiber.App, repo eventutils.IEntityRepository, resourcePath string) {
	api := Api{
		repository:   repo,
		resourcePath: resourcePath,
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
	app.Post("/categories/:id/upload-image", api.UploadCategoryImage)

	app.Post("/subcategories", api.CreateNewSubCategory)
	app.Post("/subcategories/:id/upload-image", api.UploadSubCategoryImage)

	app.Post("/menuitems", api.CreateNewMenuItem)
	app.Post("/menuitems/:id/change-name", api.ChangeMenuItemName)
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
	menu, err := api.repository.GetEntity(&entities.Menu{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	menu.(*entities.Menu).Enable()
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
	menu, err := api.repository.GetEntity(&entities.Menu{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	menu.(*entities.Menu).Disable()
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

	menu, err := api.repository.GetEntity(&entities.Menu{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	menu.(*entities.Menu).ChangeName(reqBody.NewName)
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
	menu, err := api.repository.GetEntity(&entities.Menu{}, menuID)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Menu not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu, please try again later.")
		}
	}
	if menu.(*entities.Menu).IsDeleted {
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

	category, err := api.repository.GetEntity(&entities.Category{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Category not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the category, please try again later.")
		}
	}
	category.(*entities.Category).ChangeName(reqBody.NewName)
	err = api.repository.SaveEntity(category)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusOK)
	return nil
}

func (api Api) UploadCategoryImage(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid category id")
	}

	_, err := api.repository.GetEntity(&entities.Category{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Category not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the category, please try again later.")
		}
	}

	file, err := c.FormFile("image")

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to read the uploaded image.")
	}

	path := filepath.Join(api.resourcePath, "images/categories")
	err = saveFile(id, c, file, path)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to save the uploaded image.")
	}

	c.SendStatus(fiber.StatusOK)
	return nil
}

type CreateNewSubCategoryRequest struct {
	CategoryID string `json:"categoryID"`
}

func (api Api) CreateNewSubCategory(c *fiber.Ctx) error {
	reqBody := new(CreateNewSubCategoryRequest)
	if err := c.BodyParser(reqBody); err != nil {
		return err
	}
	categoryID := uuid.FromStringOrNil(reqBody.CategoryID)
	if categoryID == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "categoryID is not valid")
	}
	category, err := api.repository.GetEntity(&entities.Category{}, categoryID)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "Category not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the category, please try again later.")
		}
	}
	if category.(*entities.Category).IsDeleted {
		return fiber.NewError(fiber.StatusNotFound, "Category not found")
	}
	subcategory := entities.NewSubCategory(categoryID)
	err = api.repository.SaveEntity(subcategory)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusCreated)
	return nil
}

type CreateNewMenuItemRequest struct {
	SubCategoryID string `json:"subCategoryID"`
}

func (api Api) CreateNewMenuItem(c *fiber.Ctx) error {
	reqBody := new(CreateNewMenuItemRequest)
	if err := c.BodyParser(reqBody); err != nil {
		return err
	}
	subCategoryID := uuid.FromStringOrNil(reqBody.SubCategoryID)
	if subCategoryID == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "subCategoryID is not valid")
	}
	subCategory, err := api.repository.GetEntity(&entities.SubCategory{}, subCategoryID)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "SubCategory not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the SubCategory, please try again later.")
		}
	}
	if subCategory.(*entities.SubCategory).IsDeleted {
		return fiber.NewError(fiber.StatusNotFound, "SubCategory not found")
	}
	menuItem := entities.NewMenuItem(subCategoryID)
	err = api.repository.SaveEntity(menuItem)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusCreated)
	return nil
}

func (api Api) UploadSubCategoryImage(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid subcategory id")
	}

	_, err := api.repository.GetEntity(&entities.SubCategory{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "SubCategory not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the SubCategory, please try again later.")
		}
	}

	file, err := c.FormFile("image")

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to read the uploaded image.")
	}

	path := filepath.Join(api.resourcePath, "images/subcategories")
	err = saveFile(id, c, file, path)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to save the uploaded image.")
	}

	c.SendStatus(fiber.StatusOK)
	return nil
}

func saveFile(id uuid.UUID, c *fiber.Ctx, file *multipart.FileHeader, resourcePath string) error {
	imageName := fmt.Sprintf("%s.jpg", id)
	path := filepath.Join(resourcePath, imageName)
	err := c.SaveFile(file, path)
	return err
}

type ChangeMenuItemNameRequest struct {
	NewName string `json:"newName"`
}

func (api Api) ChangeMenuItemName(c *fiber.Ctx) error {
	id := uuid.FromStringOrNil(c.Params("id"))
	if id == uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid menuitem id")
	}

	reqBody := new(ChangeMenuItemNameRequest)
	if err := c.BodyParser(reqBody); err != nil {
		return err
	}

	menuItem, err := api.repository.GetEntity(&entities.MenuItem{}, id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "MenuItem not found")
		} else {
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to find the menu item, please try again later.")
		}
	}
	menuItem.(*entities.MenuItem).ChangeName(reqBody.NewName)
	err = api.repository.SaveEntity(menuItem)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the changes. Please try again later")
	}
	c.SendStatus(fiber.StatusOK)
	return nil
}
