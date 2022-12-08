package internal

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Resta-Inc/resta/menu/commands/internal/entities2"
	"github.com/Resta-Inc/resta/pkg/eventutils2"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewMenu(t *testing.T) {
	// Arrange
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("*entities2.Menu")).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	url := "/menus"
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestEnableMenu(t *testing.T) {
	// Arrange
	menu := entities2.NewMenu()
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Menu{}, menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu *entities2.Menu) bool {
				return menu.IsEnabled()
			},
		)).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	url := fmt.Sprintf("/menus/%s/enable", menu.ID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestDisableMenu(t *testing.T) {
	// Arrange
	menu := entities2.NewMenu()
	menu.Enable()
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Menu{}, menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu *entities2.Menu) bool {
				return !menu.IsEnabled()
			},
		)).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	url := fmt.Sprintf("/menus/%s/disable", menu.ID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestChangeMenuName(t *testing.T) {
	// Arrange
	menu := entities2.NewMenu()
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Menu{}, menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu *entities2.Menu) bool {
				return menu.GetName() == "NewMenuName"
			},
		)).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	jsonBody := `{"newName": "NewMenuName"}`
	url := fmt.Sprintf("/menus/%s/change-name", menu.ID)
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	request.Header.Add("content-type", "application/json")
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestNewCategory(t *testing.T) {
	// Arrange
	menu := entities2.NewMenu()
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Menu{}, menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("*entities2.Category")).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	jsonBody := fmt.Sprintf(`{"menuID": "%s"}`, menu.ID)
	url := "/categories"
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	request.Header.Add("content-type", "application/json")
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestChangeCategoryName(t *testing.T) {
	// Arrange
	menu := entities2.NewMenu()
	category := entities2.NewCategory(menu.ID)
	newName := "NewCategoryName"
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Category{}, category.ID).
		Return(category, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(category *entities2.Category) bool {
				return category.GetName() == newName
			},
		)).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	jsonBody := fmt.Sprintf(`{"newName": "%s"}`, newName)
	url := fmt.Sprintf("/categories/%s/change-name", category.ID)
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	request.Header.Add("content-type", "application/json")
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestUploadCategoryPicture(t *testing.T) {
	// Arrange
	err := os.MkdirAll("./resources/images/categories", 0755)
	defer os.RemoveAll("./resources")
	require.NoError(t, err)
	menu := entities2.NewMenu()
	category := entities2.NewCategory(menu.ID)
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Category{}, category.ID).
		Return(category, nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "./resources")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("image", "test.jpg")
	require.NoError(t, err)
	file, err := os.Open("test.jpg")
	require.NoError(t, err)
	_, err = io.Copy(fw, file)
	require.NoError(t, err)
	writer.Close()

	url := fmt.Sprintf("/categories/%s/upload-image", category.ID)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	request.Header.Set("Content-Type", writer.FormDataContentType())
	require.NoError(t, err)

	// Act
	resp, err := app.Test(request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
	imageName := fmt.Sprintf("%s.jpg", category.ID)
	path := filepath.Join("./resources", "images/categories", imageName)
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func TestNewSubCategory(t *testing.T) {
	// Arrange
	category := entities2.NewCategory(utils.GenerateNewUUID())
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.Category{}, category.ID).
		Return(category, nil)

	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("*entities2.SubCategory")).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	jsonBody := fmt.Sprintf(`{"categoryID": "%s"}`, category.ID)
	url := "/subcategories"
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	request.Header.Add("content-type", "application/json")
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestNewMenuItem(t *testing.T) {
	// Arrange
	subCategory := entities2.NewSubCategory(utils.GenerateNewUUID())
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.SubCategory{}, subCategory.ID).
		Return(subCategory, nil)

	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("*entities2.MenuItem")).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	jsonBody := fmt.Sprintf(`{"subCategoryID": "%s"}`, subCategory.ID)
	url := "/menuitems"
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	request.Header.Add("content-type", "application/json")
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}

func TestUploadSubCategoryPicture(t *testing.T) {
	// Arrange
	err := os.MkdirAll("./resources/images/subcategories", 0755)
	defer os.RemoveAll("./resources")
	require.NoError(t, err)
	subcategory := entities2.NewSubCategory(utils.GenerateNewUUID())
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.SubCategory{}, subcategory.ID).
		Return(subcategory, nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "./resources")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("image", "test.jpg")
	require.NoError(t, err)
	file, err := os.Open("test.jpg")
	require.NoError(t, err)
	_, err = io.Copy(fw, file)
	require.NoError(t, err)
	writer.Close()

	url := fmt.Sprintf("/subcategories/%s/upload-image", subcategory.ID)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body.Bytes()))
	request.Header.Set("Content-Type", writer.FormDataContentType())
	require.NoError(t, err)

	// Act
	resp, err := app.Test(request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
	imageName := fmt.Sprintf("%s.jpg", subcategory.ID)
	path := filepath.Join("./resources", "images/subcategories", imageName)
	_, err = os.Stat(path)
	require.NoError(t, err)
}

func TestChangeMenuItemName(t *testing.T) {
	// Arrange
	menuItem := entities2.NewMenuItem(utils.GenerateNewUUID())
	mockEntityRepository := new(eventutils2.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", &entities2.MenuItem{}, menuItem.ID).
		Return(menuItem, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu *entities2.MenuItem) bool {
				return menu.GetName() == "NewMenuItemName"
			},
		)).
		Return(nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "")

	jsonBody := `{"newName": "NewMenuItemName"}`
	url := fmt.Sprintf("/menuitems/%s/change-name", menuItem.ID)
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	request.Header.Add("content-type", "application/json")
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEntityRepository.AssertExpectations(t)
}
