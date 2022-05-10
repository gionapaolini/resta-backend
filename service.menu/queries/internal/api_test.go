package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetMenu(t *testing.T) {
	// Arrange
	menu := MenuView{
		ID:   utils.GenerateNewUUID(),
		Name: "TestName",
	}
	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("GetMenu", menu.ID).
		Return(menu, nil)

	app := fiber.New()
	SetupApi(app, mockMenuRepository, "", "")

	url := fmt.Sprintf("/menus/%s", menu.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	response, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var menuResponse MenuView
	err = json.Unmarshal(response, &menuResponse)

	require.Equal(t, menu, menuResponse)
}

func TestGetMenus(t *testing.T) {
	// Arrange
	menus := []MenuView{
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName1",
		},
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName2",
		},
	}
	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("GetAllMenus").
		Return(menus, nil)

	app := fiber.New()
	SetupApi(app, mockMenuRepository, "", "")

	url := "/menus"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	response, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var menuResponse []MenuView
	err = json.Unmarshal(response, &menuResponse)

	require.Equal(t, menus, menuResponse)
}

func TestGetCategoriesByIDsApi(t *testing.T) {
	// Arrange
	categories := []CategoryView{
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName1",
		},
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName2",
		},
	}
	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("GetCategoriesByIDs", []uuid.UUID{
			categories[0].ID,
			categories[1].ID,
		}).
		Return(categories, nil)

	app := fiber.New()
	SetupApi(app, mockMenuRepository, "", "")

	url := fmt.Sprintf("/categories/by-ids?id=%s,%s", categories[0].ID, categories[1].ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	response, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var categoriesResponse []CategoryView
	err = json.Unmarshal(response, &categoriesResponse)

	require.Equal(t, categories[0].ID, categoriesResponse[0].ID)
	require.Equal(t, categories[0].Name, categoriesResponse[0].Name)
}

func TestGetCategoriesByIDsApi_ShouldHaveImageURL(t *testing.T) {
	// Arrange
	categories := []CategoryView{
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName1",
		},
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName2",
		},
	}

	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("GetCategoriesByIDs", []uuid.UUID{
			categories[0].ID,
			categories[1].ID,
		}).
		Return(categories, nil)

	app := fiber.New()

	err := os.MkdirAll("./resources/images/categories", 0755)
	os.OpenFile(fmt.Sprintf("./resources/images/categories/%s.jpg", categories[0].ID), os.O_RDONLY|os.O_CREATE, 0666)
	os.OpenFile(fmt.Sprintf("./resources/images/categories/%s.jpg", categories[1].ID), os.O_RDONLY|os.O_CREATE, 0666)
	defer os.RemoveAll("./resources")

	SetupApi(app, mockMenuRepository, "./resources", "http://localhost:10001")

	url := fmt.Sprintf("/categories/by-ids?id=%s,%s", categories[0].ID, categories[1].ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	response, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var categoriesResponse []CategoryView
	err = json.Unmarshal(response, &categoriesResponse)

	require.Equal(t, categoriesResponse[0].ImageURL, fmt.Sprintf("http://localhost:10001/images/categories/%s.jpg", categories[0].ID))
	require.Equal(t, categoriesResponse[1].ImageURL, fmt.Sprintf("http://localhost:10001/images/categories/%s.jpg", categories[1].ID))
}

func TestGetSubCategoriesByIDsApi(t *testing.T) {
	// Arrange
	categories := []SubCategoryView{
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName1",
		},
		{
			ID:   utils.GenerateNewUUID(),
			Name: "TestName2",
		},
	}
	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("GetSubCategoriesByIDs", []uuid.UUID{
			categories[0].ID,
			categories[1].ID,
		}).
		Return(categories, nil)

	app := fiber.New()
	SetupApi(app, mockMenuRepository, "", "")

	url := fmt.Sprintf("/subcategories/by-ids?id=%s,%s", categories[0].ID, categories[1].ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	resp, _ := app.Test(request, 100000)

	// Assert
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	response, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var subCategoriesResponse []SubCategoryView
	err = json.Unmarshal(response, &subCategoriesResponse)

	require.Equal(t, categories[0].ID, subCategoriesResponse[0].ID)
	require.Equal(t, categories[0].Name, subCategoriesResponse[0].Name)
}
