package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Resta-Inc/resta/menu/commands/internal/entities"
	"github.com/Resta-Inc/resta/pkg/eventutils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type Api struct {
	repository eventutils.IEntityRepository
}

func SetupApiOLD(router *mux.Router, repo eventutils.IEntityRepository) {
	api := Api{
		repository: repo,
	}
	api.setupRoutesOLD(router)
}

func SetupApi(app *fiber.App, repo eventutils.IEntityRepository) {
	api := Api{
		repository: repo,
	}
	api.setupRoutes(app)
}

func (api Api) setupRoutes(app *fiber.App) {
	app.Post("/menus", api.CreateNewMenu)
	// app.Post("/menus/{id}/enable", api.EnableMenu)
	// app.Post("/menus/{id}/disable", api.DisableMenu)
	// app.Post("/menus/{id}/change-name", api.ChangeMenuName)

	// app.Post("/categories", api.CreateNewCategory)
	// app.Post("/categories/{id}/change-name", api.ChangeCategoryName)
}

func (api Api) setupRoutesOLD(r *mux.Router) {
	// r.HandleFunc("/menus", api.CreateNewMenu).Methods("POST")
	r.HandleFunc("/menus/{id}/enable", api.EnableMenu).Methods("POST")
	r.HandleFunc("/menus/{id}/disable", api.DisableMenu).Methods("POST")
	r.HandleFunc("/menus/{id}/change-name", api.ChangeMenuName).Methods("POST")

	r.HandleFunc("/categories", api.CreateNewCategory).Methods("POST")
	r.HandleFunc("/categories/{id}/change-name", api.ChangeCategoryName).Methods("POST")
}

func (api Api) CreateNewMenu(c *fiber.Ctx) error {
	menu := entities.NewMenu()
	err := api.repository.SaveEntity(menu)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when saving the new menu. Please try again later")
	}
	c.SendStatus(fiber.StatusCreated)
	return nil
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

func (api Api) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var body map[string]string
	err := json.Unmarshal(reqBody, &body)
	menuID := uuid.FromStringOrNil(body["menuID"])
	if err != nil || body["menuID"] == "" || menuID == uuid.Nil {
		http.Error(w, "menuID is not valid", http.StatusBadRequest)
		return
	}
	menu, err := api.repository.GetEntity(entities.EmptyMenu(), menuID)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			http.Error(w, "Menu not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong when trying to find the menu, please try again later.", http.StatusInternalServerError)
		}
		return
	}
	if menu.(entities.Menu).IsDeleted {
		http.Error(w, "Menu not found", http.StatusNotFound)
	}
	category := entities.NewCategory(menuID)
	err = api.repository.SaveEntity(category)
	if err != nil {
		http.Error(w, "Something went wrong when saving the changes. Please try again later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (api Api) ChangeCategoryName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.FromStringOrNil(vars["id"])
	if id == uuid.Nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var body map[string]string
	err := json.Unmarshal(reqBody, &body)
	if err != nil || body["newName"] == "" {
		http.Error(w, "newName property is not valid", http.StatusBadRequest)
		return
	}
	category, err := api.repository.GetEntity(entities.EmptyCategory(), id)
	if err != nil {
		if errors.Is(err, eventutils.ErrEntityNotFound) {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, "Something went wrong when trying to find the category, please try again later.", http.StatusInternalServerError)
		}
		return
	}
	category = category.(entities.Category).ChangeName(body["newName"])
	err = api.repository.SaveEntity(category)
	if err != nil {
		http.Error(w, "Something went wrong when saving the changes. Please try again later", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
