package internal

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Resta-Inc/resta/menu/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type Api struct {
	repository eventutils.IEntityRepository
}

func NewApi(repo eventutils.IEntityRepository) Api {
	return Api{
		repository: repo,
	}
}

func (api Api) CreateNewMenu(w http.ResponseWriter, r *http.Request) {
	menu := entities.NewMenu()
	api.repository.SaveEntity(menu)
	w.WriteHeader(http.StatusCreated)
}

func (api Api) GetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			http.NotFound(w, r)
			return
		}
	}
	returnedMenu := MapToMenuResponse(menu.(entities.Menu))
	json.NewEncoder(w).Encode(returnedMenu)
}

type MenuResponse struct {
	ID uuid.UUID
	entities.MenuState
}

func MapToMenuResponse(menu entities.Menu) MenuResponse {
	return MenuResponse{
		ID:        menu.GetID(),
		MenuState: menu.State,
	}
}
