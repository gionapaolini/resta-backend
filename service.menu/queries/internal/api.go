package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type Api struct {
	menuRepository IMenuRepository
	Router         *mux.Router
}

func SetupApi(router *mux.Router, repo IMenuRepository) {
	api := Api{
		menuRepository: repo,
	}
	api.setupRoutes(router)
}

func (api Api) setupRoutes(r *mux.Router) {
	r.HandleFunc("/menus/{id}", api.GetMenu)
	r.HandleFunc("/menus", api.GetAllMenus)
}

func (api Api) GetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.menuRepository.GetMenu(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong when trying to find the menu, please try again later.", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(menu)
}

func (api Api) GetAllMenus(w http.ResponseWriter, r *http.Request) {
	menu, err := api.menuRepository.GetAllMenus()
	if err != nil {
		http.Error(w, "Something went wrong when trying to find the menus, please try again later.", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(menu)
}
