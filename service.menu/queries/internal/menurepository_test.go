package internal

import (
	"fmt"
	"testing"

	"github.com/Resta-Inc/resta/pkg/utils"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

var pgConnectionString string = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

func TestCreateMenu(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	menuName := "TestMenu"

	viewRepository := NewMenuRepository(pgConnectionString)
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
	viewRepository := NewMenuRepository(pgConnectionString)
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

func TestEnableMenu(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	menuName := "TestMenu"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteMenu(menuID)
	err := viewRepository.CreateMenu(menuID, menuName)

	// Act
	err = viewRepository.EnableMenu(menuID)

	// Assert
	require.NoError(t, err)
	returnedMenu, err := viewRepository.GetMenu(menuID)
	require.NoError(t, err)
	require.True(t, returnedMenu.IsEnabled)
}

func TestDisableMenu(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	menuName := "TestMenu"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteMenu(menuID)
	err := viewRepository.CreateMenu(menuID, menuName)
	err = viewRepository.EnableMenu(menuID)

	// Act
	err = viewRepository.DisableMenu(menuID)

	// Assert
	require.NoError(t, err)
	returnedMenu, err := viewRepository.GetMenu(menuID)
	require.NoError(t, err)
	require.False(t, returnedMenu.IsEnabled)
}

func TestChangeMenuName(t *testing.T) {
	// Arrange
	menuID := utils.GenerateNewUUID()
	menuName := "TestMenu"
	newMenuName := "NewMenuName"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteMenu(menuID)
	err := viewRepository.CreateMenu(menuID, menuName)

	// Act
	err = viewRepository.ChangeMenuName(menuID, newMenuName)

	// Assert
	require.NoError(t, err)
	returnedMenu, err := viewRepository.GetMenu(menuID)
	require.NoError(t, err)
	require.Equal(t, newMenuName, returnedMenu.Name)
}

func TestCreateCategory(t *testing.T) {
	// Arrange
	categoryID, categoryName, imageURL :=
		utils.GenerateNewUUID(), "TestCategory", "test.com"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteCategory(categoryID)

	// Act
	err := viewRepository.CreateCategory(categoryID, categoryName, imageURL)

	// Assert
	require.NoError(t, err)
	returnedCategory, err := viewRepository.GetCategory(categoryID)
	require.NoError(t, err)
	require.Equal(t, categoryID, returnedCategory.ID)
	require.Equal(t, categoryName, returnedCategory.Name)
	require.Equal(t, imageURL, returnedCategory.ImageURL)
}

func TestAddCategoryToMenu(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	menuID, menuName, categoryID, categoryName, imageURL :=
		utils.GenerateNewUUID(),
		"TestMenu",
		utils.GenerateNewUUID(),
		"TestCategory",
		"test.com"

	defer viewRepository.DeleteCategory(categoryID)
	defer viewRepository.RemoveCategoryFromMenu(menuID, categoryID)
	defer viewRepository.DeleteMenu(menuID)
	viewRepository.CreateMenu(menuID, menuName)
	viewRepository.CreateCategory(categoryID, categoryName, imageURL)

	// Act
	err := viewRepository.AddCategoryToMenu(menuID, categoryID)

	// Assert
	require.NoError(t, err)
	returnedMenu, err := viewRepository.GetMenu(menuID)
	require.NoError(t, err)
	require.Len(t, returnedMenu.CategoriesIDs, 1)
	require.Equal(t, returnedMenu.CategoriesIDs[0], categoryID)
}

func TestGetMenu_ShouldHaveCategoriesIDPopulated(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	menuID, menuName, categoryID1, categoryID2, categoryName, imageURL :=
		utils.GenerateNewUUID(),
		"TestMenu",
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestCategory",
		"test.com"

	defer viewRepository.DeleteCategory(categoryID1)
	defer viewRepository.DeleteCategory(categoryID2)
	defer viewRepository.RemoveCategoryFromMenu(menuID, categoryID1)
	defer viewRepository.RemoveCategoryFromMenu(menuID, categoryID2)
	defer viewRepository.DeleteMenu(menuID)
	viewRepository.CreateMenu(menuID, menuName)
	viewRepository.CreateCategory(categoryID1, categoryName, imageURL)
	viewRepository.CreateCategory(categoryID2, categoryName, imageURL)
	viewRepository.AddCategoryToMenu(menuID, categoryID1)
	viewRepository.AddCategoryToMenu(menuID, categoryID2)

	// Act
	returnedMenu, err := viewRepository.GetMenu(menuID)

	// Assert
	require.NoError(t, err)
	require.Len(t, returnedMenu.CategoriesIDs, 2)
	require.Contains(t, returnedMenu.CategoriesIDs, categoryID1)
	require.Contains(t, returnedMenu.CategoriesIDs, categoryID2)
}

func TestGetMenu_WhenNoCategory_ShouldHaveEmptyCategoryList(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	menuID, menuName := utils.GenerateNewUUID(), "TestMenu"

	defer viewRepository.DeleteMenu(menuID)
	viewRepository.CreateMenu(menuID, menuName)

	// Act
	returnedMenu, err := viewRepository.GetMenu(menuID)

	// Assert
	require.NoError(t, err)
	require.Len(t, returnedMenu.CategoriesIDs, 0)
}

func TestGetCategoriesByIDs(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	categoryID1, categoryID2, categoryName, imageURL :=
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestCategory",
		"test.com"

	defer viewRepository.DeleteCategory(categoryID1)
	defer viewRepository.DeleteCategory(categoryID2)
	viewRepository.CreateCategory(categoryID1, categoryName, imageURL)
	viewRepository.CreateCategory(categoryID2, categoryName, imageURL)

	// Act
	categories, err := viewRepository.GetCategoriesByIDs([]uuid.UUID{categoryID1, categoryID2})

	// Assert
	require.NoError(t, err)
	require.Len(t, categories, 2)
	require.Equal(t, categories[0].ID, categoryID1)
	require.Equal(t, categories[1].ID, categoryID2)
}
