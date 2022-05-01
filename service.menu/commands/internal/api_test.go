package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewMenu(t *testing.T) {
	// Arrange
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("entities.Menu")).
		Return(nil)
	router := mux.NewRouter()
	SetupApi(router, mockEntityRepository)
	recorder := httptest.NewRecorder()

	url := "/menus"
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusCreated, recorder.Code)
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
	router := mux.NewRouter()
	SetupApi(router, mockEntityRepository)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/menus/%s/enable", menu.ID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)
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
	router := mux.NewRouter()
	SetupApi(router, mockEntityRepository)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/menus/%s/disable", menu.ID)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)
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
	router := mux.NewRouter()
	SetupApi(router, mockEntityRepository)
	recorder := httptest.NewRecorder()
	jsonBody := `{"newName": "NewMenuName"}`
	url := fmt.Sprintf("/menus/%s/change-name", menu.ID)
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)
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
	router := mux.NewRouter()
	SetupApi(router, mockEntityRepository)
	recorder := httptest.NewRecorder()
	jsonBody := fmt.Sprintf(`{"menuID": "%s"}`, menu.ID)
	url := "/categories"
	request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(jsonBody))
	require.NoError(t, err)

	// Act
	router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusCreated, recorder.Code)
	mockEntityRepository.AssertExpectations(t)
}
