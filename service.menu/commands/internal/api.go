package internal

import (
	"encoding/json"
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
		http.Error(w, "Something went wrong on our servers. Please re-try later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (api Api) EnableMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	menu = menu.(entities.Menu).Enable()
	err = api.repository.SaveEntity(menu)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api Api) DisableMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	menu = menu.(entities.Menu).Disable()
	err = api.repository.SaveEntity(menu)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api Api) ChangeMenuName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var body map[string]string
	err := json.Unmarshal(reqBody, &body)
	if err != nil || body["newName"] == "" {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), id)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	menu = menu.(entities.Menu).ChangeName(body["newName"])
	err = api.repository.SaveEntity(menu)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
