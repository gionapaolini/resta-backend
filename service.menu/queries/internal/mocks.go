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

func (m MockMenuRepository) CreateCategory(categoryID uuid.UUID, categoryName, imageURL string) error {
	args := m.Called(categoryID, categoryName, imageURL)
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
