package internal

import "github.com/gofrs/uuid"

type MenuView struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	IsEnabled     bool        `json:"isEnabled"`
	CategoriesIDs []uuid.UUID `json:"categoriesIDs"`
}

type CategoryView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
