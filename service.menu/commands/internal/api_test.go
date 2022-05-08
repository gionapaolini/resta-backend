package internal

import (
	"fmt"
	"net/http"
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
	SetupApi(app, mockEntityRepository)

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
	SetupApi(app, mockEntityRepository)

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
	SetupApi(app, mockEntityRepository)

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
	SetupApi(app, mockEntityRepository)

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
	SetupApi(app, mockEntityRepository)

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
	SetupApi(app, mockEntityRepository)

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
