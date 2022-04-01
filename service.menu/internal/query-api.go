package internal

import (
	"encoding/json"
	"net/http"

	"github.com/Resta-Inc/resta/menu/internal/dataviews"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type QueryApi struct {
	menuRepository dataviews.IMenuRepository
	Router         *mux.Router
}

func SetupQueryApi(router *mux.Router, repo dataviews.IMenuRepository) {
	api := QueryApi{
		menuRepository: repo,
	}
	api.setupRoutes(router)
}

func (api QueryApi) setupRoutes(r *mux.Router) {
	r.HandleFunc("/menus/{id}", api.GetMenu)
}

func (api QueryApi) GetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "invalid menu id", http.StatusBadRequest)
		return
	}
	menu, err := api.menuRepository.GetMenu(id)
	if err != nil {
		//FIX IT with not found as well
		http.Error(w, "something wrong", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(menu)
}
