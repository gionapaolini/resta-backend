package internal

import (
	"time"

	"github.com/gofrs/uuid"
)

type MenuView struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	IsEnabled     bool        `json:"isEnabled"`
	CategoriesIDs []uuid.UUID `json:"categoriesIDs"`
	CreatedAt     time.Time   `json:"createdAt"`
}

type CategoryView struct {
	ID               uuid.UUID   `json:"id"`
	Name             string      `json:"name"`
	ImageURL         string      `json:"imageURL"`
	SubCategoriesIDs []uuid.UUID `json:"subCategoriesIDs"`
	CreatedAt        time.Time   `json:"createdAt"`
}

type SubCategoryView struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	ImageURL     string      `json:"imageURL"`
	MenuItemsIDs []uuid.UUID `json:"menuItemsIDs"`
	CreatedAt    time.Time   `json:"createdAt"`
}

type MenuItemView struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}
