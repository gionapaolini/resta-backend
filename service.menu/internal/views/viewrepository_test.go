package views

import (
	"fmt"
	"testing"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/stretchr/testify/require"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

var connectionString string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func TestCreateMenu(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	menuName := "TestMenu"

	viewRepository := NewViewRepository(connectionString)
	defer viewRepository.DeleteMenu(menuID)

	// Act
	err := viewRepository.CreateMenu(menuID, menuName)

	// Assert
	require.NoError(t, err)
	returnedMenu, err := viewRepository.GetMenu(menuID)
	require.NoError(t, err)
	require.Equal(t, menuID, returnedMenu.ID)
	require.Equal(t, menuName, returnedMenu.Name)
}

func TestGetAllMenus(t *testing.T) {
	// Arrange
	menuID1, menuID2, menuID3 := utils.GenerateNewUUID(), utils.GenerateNewUUID(), utils.GenerateNewUUID()
	menuName := "TestMenu"
	viewRepository := NewViewRepository(connectionString)
	defer viewRepository.DeleteMenu(menuID1)
	defer viewRepository.DeleteMenu(menuID2)
	defer viewRepository.DeleteMenu(menuID3)
	_ = viewRepository.CreateMenu(menuID1, menuName)
	_ = viewRepository.CreateMenu(menuID2, menuName)
	_ = viewRepository.CreateMenu(menuID3, menuName)

	// Act
	menus, err := viewRepository.GetAllMenus()

	// Assert
	require.NoError(t, err)
	require.Len(t, menus, 3)
}
