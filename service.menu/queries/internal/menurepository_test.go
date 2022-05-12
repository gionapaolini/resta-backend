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
	require.True(t, len(menus) >= 3) //there might already be menus since we use the same DB
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
	categoryID, categoryName :=
		utils.GenerateNewUUID(), "TestCategory"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteCategory(categoryID)

	// Act
	err := viewRepository.CreateCategory(categoryID, categoryName)

	// Assert
	require.NoError(t, err)
	returnedCategory, err := viewRepository.GetCategory(categoryID)
	require.NoError(t, err)
	require.Equal(t, categoryID, returnedCategory.ID)
	require.Equal(t, categoryName, returnedCategory.Name)
}

func TestAddCategoryToMenu(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	menuID, menuName, categoryID, categoryName :=
		utils.GenerateNewUUID(),
		"TestMenu",
		utils.GenerateNewUUID(),
		"TestCategory"

	defer viewRepository.DeleteCategory(categoryID)
	defer viewRepository.RemoveCategoryFromMenu(menuID, categoryID)
	defer viewRepository.DeleteMenu(menuID)
	viewRepository.CreateMenu(menuID, menuName)
	viewRepository.CreateCategory(categoryID, categoryName)

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
	menuID, menuName, categoryID1, categoryID2, categoryName :=
		utils.GenerateNewUUID(),
		"TestMenu",
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestCategory"

	defer viewRepository.DeleteCategory(categoryID1)
	defer viewRepository.DeleteCategory(categoryID2)
	defer viewRepository.RemoveCategoryFromMenu(menuID, categoryID1)
	defer viewRepository.RemoveCategoryFromMenu(menuID, categoryID2)
	defer viewRepository.DeleteMenu(menuID)
	viewRepository.CreateMenu(menuID, menuName)
	viewRepository.CreateCategory(categoryID1, categoryName)
	viewRepository.CreateCategory(categoryID2, categoryName)
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
	categoryID1, categoryID2, categoryName :=
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestCategory"

	defer viewRepository.DeleteCategory(categoryID1)
	defer viewRepository.DeleteCategory(categoryID2)
	viewRepository.CreateCategory(categoryID1, categoryName)
	viewRepository.CreateCategory(categoryID2, categoryName)

	// Act
	categories, err := viewRepository.GetCategoriesByIDs([]uuid.UUID{categoryID1, categoryID2})

	// Assert
	require.NoError(t, err)
	require.Len(t, categories, 2)
	require.Equal(t, categories[0].ID, categoryID1)
	require.Equal(t, categories[1].ID, categoryID2)
}

func TestGetCategoriesByIDs_ShouldHaveSubCategoriesIDPopulated(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	categoryID, subCategoryID1, subCategoryID2, subCategoryName :=
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestSubCategory"

	defer viewRepository.DeleteCategory(categoryID)
	defer viewRepository.DeleteSubCategory(subCategoryID1)
	defer viewRepository.DeleteSubCategory(subCategoryID2)
	defer viewRepository.RemoveSubCategoryFromCategory(categoryID, subCategoryID1)
	defer viewRepository.RemoveSubCategoryFromCategory(categoryID, subCategoryID2)

	viewRepository.CreateCategory(categoryID, subCategoryName)
	viewRepository.CreateSubCategory(subCategoryID1, subCategoryName)
	viewRepository.CreateSubCategory(subCategoryID2, subCategoryName)
	viewRepository.AddSubCategoryToCategory(categoryID, subCategoryID1)
	viewRepository.AddSubCategoryToCategory(categoryID, subCategoryID2)

	// Act
	categories, err := viewRepository.GetCategoriesByIDs([]uuid.UUID{categoryID})

	// Assert
	require.NoError(t, err)
	require.Contains(t, categories[0].SubCategoriesIDs, subCategoryID1)
	require.Contains(t, categories[0].SubCategoriesIDs, subCategoryID2)
}

func TestChangeCategoryName(t *testing.T) {
	// Arrange
	newName := "NewName"
	viewRepository := NewMenuRepository(pgConnectionString)
	categoryID, categoryName :=
		utils.GenerateNewUUID(),
		"TestCategory"

	defer viewRepository.DeleteCategory(categoryID)
	viewRepository.CreateCategory(categoryID, categoryName)

	// Act
	err := viewRepository.ChangeCategoryName(categoryID, newName)

	// Assert
	require.NoError(t, err)
	returnedCategory, err := viewRepository.GetCategory(categoryID)
	require.NoError(t, err)
	require.Equal(t, newName, returnedCategory.Name)
}

func TestCreateSubCategory(t *testing.T) {
	// Arrange
	subCategoryID, subCategoryName :=
		utils.GenerateNewUUID(), "TestSubCategory"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteSubCategory(subCategoryID)

	// Act
	err := viewRepository.CreateSubCategory(subCategoryID, subCategoryName)

	// Assert
	require.NoError(t, err)
	returnedSubCategory, err := viewRepository.GetSubCategory(subCategoryID)
	require.NoError(t, err)
	require.Equal(t, subCategoryID, returnedSubCategory.ID)
	require.Equal(t, subCategoryName, returnedSubCategory.Name)
}

func TestAddSubCategoryToCategory(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	categoryID, categoryName, subCategoryID, subCategoryName :=
		utils.GenerateNewUUID(),
		"TestCategory",
		utils.GenerateNewUUID(),
		"TestSubCategory"

	defer viewRepository.DeleteSubCategory(subCategoryID)
	defer viewRepository.RemoveSubCategoryFromCategory(categoryID, subCategoryID)
	defer viewRepository.DeleteCategory(categoryID)
	viewRepository.CreateCategory(categoryID, categoryName)
	viewRepository.CreateSubCategory(subCategoryID, subCategoryName)

	// Act
	err := viewRepository.AddSubCategoryToCategory(categoryID, subCategoryID)

	// Assert
	require.NoError(t, err)
	returnedCategory, err := viewRepository.GetCategory(categoryID)
	require.NoError(t, err)
	require.Len(t, returnedCategory.SubCategoriesIDs, 1)
	require.Equal(t, returnedCategory.SubCategoriesIDs[0], subCategoryID)
}

func TestCreateMenuItem(t *testing.T) {
	// Arrange
	menuItemID, menuItemName :=
		utils.GenerateNewUUID(), "TestMenuItem"

	viewRepository := NewMenuRepository(pgConnectionString)
	defer viewRepository.DeleteMenuItem(menuItemID)

	// Act
	err := viewRepository.CreateMenuItem(menuItemID, menuItemName)

	// Assert
	require.NoError(t, err)
	returnedMenuItem, err := viewRepository.GetMenuItem(menuItemID)
	require.NoError(t, err)
	require.Equal(t, menuItemID, returnedMenuItem.ID)
	require.Equal(t, menuItemName, returnedMenuItem.Name)
}

func TestAddMenuItemToSubCategory(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	subCategoryID, subCategoryName, menuItemID, menuItemName :=
		utils.GenerateNewUUID(),
		"TestCategory",
		utils.GenerateNewUUID(),
		"TestMenuItem"

	defer viewRepository.DeleteMenuItem(menuItemID)
	defer viewRepository.RemoveMenuItemFromSubCategory(subCategoryID, menuItemID)
	defer viewRepository.DeleteCategory(subCategoryID)
	viewRepository.CreateSubCategory(subCategoryID, subCategoryName)
	viewRepository.CreateMenuItem(menuItemID, menuItemName)

	// Act
	err := viewRepository.AddMenuItemToSubCategory(subCategoryID, menuItemID)

	// Assert
	require.NoError(t, err)
	returnedSubCategory, err := viewRepository.GetSubCategory(subCategoryID)
	require.NoError(t, err)
	require.Len(t, returnedSubCategory.MenuItemsIDs, 1)
	require.Equal(t, returnedSubCategory.MenuItemsIDs[0], menuItemID)
}

func TestGetSubCategoriesByIDs_ShouldHaveMenuItemsIDsPopulated(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	subCategoryID, menuItem1, menuItem2, menuItemName :=
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestSubCategory"

	defer viewRepository.DeleteCategory(subCategoryID)
	defer viewRepository.DeleteSubCategory(menuItem1)
	defer viewRepository.DeleteSubCategory(menuItem2)
	defer viewRepository.RemoveSubCategoryFromCategory(subCategoryID, menuItem1)
	defer viewRepository.RemoveSubCategoryFromCategory(subCategoryID, menuItem2)

	viewRepository.CreateSubCategory(subCategoryID, menuItemName)
	viewRepository.CreateMenuItem(menuItem1, menuItemName)
	viewRepository.CreateMenuItem(menuItem2, menuItemName)
	viewRepository.AddMenuItemToSubCategory(subCategoryID, menuItem1)
	viewRepository.AddMenuItemToSubCategory(subCategoryID, menuItem2)

	// Act
	categories, err := viewRepository.GetSubCategoriesByIDs([]uuid.UUID{subCategoryID})

	// Assert
	require.NoError(t, err)
	require.Contains(t, categories[0].MenuItemsIDs, menuItem1)
	require.Contains(t, categories[0].MenuItemsIDs, menuItem2)
}

func TestGetMenuItemsByIDs(t *testing.T) {
	// Arrange
	viewRepository := NewMenuRepository(pgConnectionString)
	menuItemID1, menuItemID2, menuItemName :=
		utils.GenerateNewUUID(),
		utils.GenerateNewUUID(),
		"TestName"

	defer viewRepository.DeleteMenuItem(menuItemID1)
	defer viewRepository.DeleteMenuItem(menuItemID2)
	viewRepository.CreateMenuItem(menuItemID1, menuItemName)
	viewRepository.CreateMenuItem(menuItemID2, menuItemName)

	// Act
	menuItems, err := viewRepository.GetMenuItemsByIDs([]uuid.UUID{menuItemID1, menuItemID2})

	// Assert
	require.NoError(t, err)
	require.Len(t, menuItems, 2)
	require.Equal(t, menuItems[0].ID, menuItemID1)
	require.Equal(t, menuItems[1].ID, menuItemID2)
}
