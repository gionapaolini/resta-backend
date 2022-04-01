package internal

import (
	"net/http"

	"github.com/Resta-Inc/resta/menu/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gorilla/mux"
)

type CommandsApi struct {
	repository eventutils.IEntityRepository
}

func SetupCommandsApi(router *mux.Router, repo eventutils.IEntityRepository) {
	api := CommandsApi{
		repository: repo,
	}
	api.setupRoutes(router)
}

func (api CommandsApi) setupRoutes(r *mux.Router) {
	r.HandleFunc("/menus", api.CreateNewMenu).Methods("POST")
}

func (api CommandsApi) CreateNewMenu(w http.ResponseWriter, r *http.Request) {
	menu := entities.NewMenu()
	api.repository.SaveEntity(menu)
	w.WriteHeader(http.StatusCreated)
}
