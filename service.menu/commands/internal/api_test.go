package internal

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
