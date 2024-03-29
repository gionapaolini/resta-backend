package internal

import (
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMenuRepository struct {
	mock.Mock
}

func (m MockMenuRepository) CreateMenu(menuID uuid.UUID, menuName string) error {
	args := m.Called(menuID, menuName)
	return args.Error(0)
}

func (m MockMenuRepository) GetMenu(menuID uuid.UUID) (MenuView, error) {
	args := m.Called(menuID)
	menuView, _ := args.Get(0).(MenuView)
	return menuView, args.Error(1)
}

func (m MockMenuRepository) GetAllMenus() ([]MenuView, error) {
	args := m.Called()
	menuViews, _ := args.Get(0).([]MenuView)
	return menuViews, args.Error(1)
}

func (m MockMenuRepository) DeleteMenu(menuID uuid.UUID) error {
	args := m.Called(menuID)
	return args.Error(0)
}

func (m MockMenuRepository) EnableMenu(menuID uuid.UUID) error {
	args := m.Called(menuID)
	return args.Error(0)
}

func (m MockMenuRepository) DisableMenu(menuID uuid.UUID) error {
	args := m.Called(menuID)
	return args.Error(0)
}

func (m MockMenuRepository) ChangeMenuName(menuID uuid.UUID, newName string) error {
	args := m.Called(menuID, newName)
	return args.Error(0)
}

func (m MockMenuRepository) CreateCategory(categoryID uuid.UUID, categoryName string) error {
	args := m.Called(categoryID, categoryName)
	return args.Error(0)
}

func (m MockMenuRepository) GetCategory(categoryID uuid.UUID) (CategoryView, error) {
	args := m.Called(categoryID)
	categoryView, _ := args.Get(0).(CategoryView)
	return categoryView, args.Error(1)
}

func (m MockMenuRepository) AddCategoryToMenu(menuID, categoryID uuid.UUID) error {
	args := m.Called(menuID, categoryID)
	return args.Error(0)
}

func (m MockMenuRepository) GetCategoriesByIDs(categoriesIDs []uuid.UUID) ([]CategoryView, error) {
	args := m.Called(categoriesIDs)
	categoriesViews, _ := args.Get(0).([]CategoryView)
	return categoriesViews, args.Error(1)
}

func (m MockMenuRepository) ChangeCategoryName(categoryID uuid.UUID, newName string) error {
	args := m.Called(categoryID, newName)
	return args.Error(0)
}

func (m MockMenuRepository) CreateSubCategory(subCategoryID uuid.UUID, subCategoryName string) error {
	args := m.Called(subCategoryID, subCategoryName)
	return args.Error(0)
}

func (m MockMenuRepository) GetSubCategoriesByIDs(subCategoriesIDs []uuid.UUID) ([]SubCategoryView, error) {
	args := m.Called(subCategoriesIDs)
	subCategoriesViews, _ := args.Get(0).([]SubCategoryView)
	return subCategoriesViews, args.Error(1)
}

func (m MockMenuRepository) AddSubCategoryToCategory(categoryID, subCategoryID uuid.UUID) error {
	args := m.Called(categoryID, subCategoryID)
	return args.Error(0)
}

func (m MockMenuRepository) CreateMenuItem(menuItemID uuid.UUID, menuItemName string) error {
	args := m.Called(menuItemID, menuItemName)
	return args.Error(0)
}

func (m MockMenuRepository) AddMenuItemToSubCategory(subCategoryID, menuItemID uuid.UUID) error {
	args := m.Called(subCategoryID, menuItemID)
	return args.Error(0)
}

func (m MockMenuRepository) GetMenuItemsByIDs(menuItemsIDs []uuid.UUID) ([]MenuItemView, error) {
	args := m.Called(menuItemsIDs)
	menuItemsViews, _ := args.Get(0).([]MenuItemView)
	return menuItemsViews, args.Error(1)
}
