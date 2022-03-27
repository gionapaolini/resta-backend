package internal

import (
	"net/http"

	"github.com/Resta-Inc/resta/menu/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
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
