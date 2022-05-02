package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	SetupApi(router, mockMenuRepository)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/menus/%s", menu.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)

	response, err := ioutil.ReadAll(recorder.Body)
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

	router := mux.NewRouter()
	SetupApi(router, mockMenuRepository)
	recorder := httptest.NewRecorder()

	url := "/menus"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)

	response, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var menuResponse []MenuView
	err = json.Unmarshal(response, &menuResponse)

	require.Equal(t, menus, menuResponse)
}

func TestGetCategoriesByIDsApi(t *testing.T) {
	// Arrange
	categories := []CategoryView{
		{
			ID:       utils.GenerateNewUUID(),
			Name:     "TestName1",
			ImageURL: "Jello1",
		},
		{
			ID:       utils.GenerateNewUUID(),
			Name:     "TestName2",
			ImageURL: "Jello2",
		},
	}
	mockMenuRepository := new(MockMenuRepository)
	mockMenuRepository.
		On("GetCategoriesByIDs", []uuid.UUID{
			categories[0].ID,
			categories[1].ID,
		}).
		Return(categories, nil)

	router := mux.NewRouter()
	SetupApi(router, mockMenuRepository)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/categories/by-ids?id=%s,%s", categories[0].ID, categories[1].ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)

	response, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var categoriesResponse []CategoryView
	err = json.Unmarshal(response, &categoriesResponse)

	require.Equal(t, categories, categoriesResponse)
}
