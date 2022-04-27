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
