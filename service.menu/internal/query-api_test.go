package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Resta-Inc/resta/menu/internal/dataviews"
	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestGetMenu(t *testing.T) {
	// Arrange
	menu := dataviews.MenuView{
		ID:   utils.GenerateNewUUID(),
		Name: "TestName",
	}
	mockMenuRepository := new(dataviews.MockMenuRepository)
	mockMenuRepository.
		On("GetMenu", menu.ID).
		Return(menu, nil)

	router := mux.NewRouter()
	SetupQueryApi(router, mockMenuRepository)
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

	var menuResponse dataviews.MenuView
	err = json.Unmarshal(response, &menuResponse)

	require.Equal(t, menu, menuResponse)
}
