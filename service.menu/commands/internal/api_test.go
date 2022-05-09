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

	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewMenu(t *testing.T) {
	// Arrange
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("entities.Menu")).
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
	menu := entities.NewMenu()
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyMenu(), menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu entities.Menu) bool {
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
	menu := entities.NewMenu().Enable()
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyMenu(), menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu entities.Menu) bool {
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
	menu := entities.NewMenu()
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyMenu(), menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(menu entities.Menu) bool {
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
	menu := entities.NewMenu()
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyMenu(), menu.ID).
		Return(menu, nil)

	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("Category")).
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
	menu := entities.NewMenu()
	category := entities.NewCategory(menu.ID)
	newName := "NewCategoryName"
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyCategory(), category.ID).
		Return(category, nil)

	mockEntityRepository.
		On("SaveEntity", mock.MatchedBy(
			func(category entities.Category) bool {
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
	menu := entities.NewMenu()
	category := entities.NewCategory(menu.ID)
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyCategory(), category.ID).
		Return(category, nil)

	app := fiber.New()
	SetupApi(app, mockEntityRepository, "./resources")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("image", "test_category_image.jpg")
	require.NoError(t, err)
	file, err := os.Open("test_category_image.jpg")
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
