package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Resta-Inc/resta/menu/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewMenu(t *testing.T) {
	// Arrange
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("entities.Menu")).
		Return(nil)

	api := NewCommandsApi(mockEntityRepository)
	recorder := httptest.NewRecorder()

	url := "/menus"
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	// Act
	api.Router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestGetMenu(t *testing.T) {
	// Arrange
	menu := eventutils.CommitEvents(entities.NewMenu())
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("GetEntity", entities.EmptyMenu(), menu.GetID()).
		Return(menu, nil)

	api := NewCommandsApi(mockEntityRepository)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/menus/%s", menu.GetID())
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// Act
	api.Router.ServeHTTP(recorder, request)

	// Assert
	require.Equal(t, http.StatusOK, recorder.Code)

	response, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	var menuResponse MenuResponse

	err = json.Unmarshal(response, &menuResponse)

	require.Equal(t, MapToMenuResponse(menu.(entities.Menu)), menuResponse)
}
