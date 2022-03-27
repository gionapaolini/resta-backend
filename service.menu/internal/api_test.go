package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewMenu(t *testing.T) {
	mockEntityRepository := new(eventutils.MockEntityRepository)
	mockEntityRepository.
		On("SaveEntity", mock.AnythingOfType("entities.Menu")).
		Return(nil)

	api := NewApi(mockEntityRepository)
	recorder := httptest.NewRecorder()

	url := "/menus"
	request, err := http.NewRequest(http.MethodPost, url, nil)
	require.NoError(t, err)

	api.router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusCreated, recorder.Code)
}
