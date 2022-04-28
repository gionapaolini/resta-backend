package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type Api struct {
	repository eventutils.IEntityRepository
}

func SetupApi(router *mux.Router, repo eventutils.IEntityRepository) {
	api := Api{
		repository: repo,
	}
	api.setupRoutes(router)
}

func (api Api) setupRoutes(r *mux.Router) {
	r.HandleFunc("/menus", api.CreateNewMenu).Methods("POST")
	r.HandleFunc("/menus/{id}/enable", api.EnableMenu).Methods("POST")
	r.HandleFunc("/menus/{id}/disable", api.DisableMenu).Methods("POST")
	r.HandleFunc("/menus/{id}/change-name", api.ChangeMenuName).Methods("POST")
}

func (api Api) CreateNewMenu(w http.ResponseWriter, r *http.Request) {
	menu := entities.NewMenu()
	err := api.repository.SaveEntity(menu)
	if err != nil {
		http.Error(w, "Something went wrong when saving the new menu. Please try again later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (api Api) EnableMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "Invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			http.Error(w, "Menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong when trying to find the menu, please try again later.", http.StatusInternalServerError)
		}
		return
	}
	menu = menu.(entities.Menu).Enable()
	err = api.repository.SaveEntity(menu)
	if err != nil {
		http.Error(w, "Something went wrong when saving the changes. Please try again later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api Api) DisableMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "Invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			http.Error(w, "Menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong when trying to find the menu, please try again later.", http.StatusInternalServerError)
		}
		return
	}
	menu = menu.(entities.Menu).Disable()
	err = api.repository.SaveEntity(menu)
	if err != nil {
		http.Error(w, "Something went wrong when saving the changes. Please try again later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api Api) ChangeMenuName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "Invalid menu id", http.StatusBadRequest)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var body map[string]string
	err := json.Unmarshal(reqBody, &body)
	if err != nil || body["newName"] == "" {
		http.Error(w, "newName property is not valid", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			http.Error(w, "Menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong when trying to find the menu, please try again later.", http.StatusInternalServerError)
		}
		return
	}
	menu = menu.(entities.Menu).ChangeName(body["newName"])
	err = api.repository.SaveEntity(menu)
	if err != nil {
		http.Error(w, "Something went wrong when saving the changes. Please try again later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
