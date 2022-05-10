package internal

import "github.com/gofrs/uuid"

type MenuView struct {
	ID            uuid.UUID   `json:"id"`
	Name          string      `json:"name"`
	IsEnabled     bool        `json:"isEnabled"`
	CategoriesIDs []uuid.UUID `json:"categoriesIDs"`
}

type CategoryView struct {
	ID               uuid.UUID   `json:"id"`
	Name             string      `json:"name"`
	ImageURL         string      `json:"imageURL"`
	SubCategoriesIDs []uuid.UUID `json:"subCategoriesIDs"`
}

type SubCategoryView struct {
	ID           uuid.UUID   `json:"id"`
	Name         string      `json:"name"`
	ImageURL     string      `json:"imageURL"`
	MenuItemsIDs []uuid.UUID `json:"menuItemsIDs"`
}

type MenuItemView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
