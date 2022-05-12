package internal

import (
	"database/sql"
	"log"
	"strings"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

type IMenuRepository interface {
	CreateMenu(menuID uuid.UUID, menuName string) error
	GetMenu(menuID uuid.UUID) (MenuView, error)
	GetAllMenus() ([]MenuView, error)
	DeleteMenu(menuID uuid.UUID) error
	EnableMenu(menuID uuid.UUID) error
	DisableMenu(menuID uuid.UUID) error
	ChangeMenuName(menuID uuid.UUID, newName string) error
	CreateCategory(categoryID uuid.UUID, categoryName string) error
	AddCategoryToMenu(menuID, categoryID uuid.UUID) error
	GetCategoriesByIDs(categoriesIDs []uuid.UUID) ([]CategoryView, error)
	ChangeCategoryName(categoryID uuid.UUID, newName string) error
	CreateSubCategory(subCategoryID uuid.UUID, subCategoryName string) error
	GetSubCategoriesByIDs(subCategoriesIDs []uuid.UUID) ([]SubCategoryView, error)
	AddSubCategoryToCategory(categoryID, subCategoryID uuid.UUID) error
	CreateMenuItem(menuItemID uuid.UUID, menuItemName string) error
	AddMenuItemToSubCategory(subCategoryID, menuItemID uuid.UUID) error
	GetMenuItemsByIDs(menuItemsIDs []uuid.UUID) ([]MenuItemView, error)
}

type MenuRepository struct {
	connectionString string
}

func NewMenuRepository(connectionString string) MenuRepository {
	return MenuRepository{
		connectionString: connectionString,
	}
}

func (repo MenuRepository) CreateMenu(menuID uuid.UUID, menuName string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO menus ("id", "name") VALUES ($1, $2)`
	_, err = db.Exec(query, menuID, menuName)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) GetMenu(menuID uuid.UUID) (MenuView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return MenuView{}, err
	}
	defer db.Close()

	var menuView MenuView

	query := `
		SELECT m.*, array_agg(mc.category_id) AS ids
		FROM menus m
		LEFT JOIN menus_categories mc ON m.id = mc.menu_id
		WHERE m.id=$1
		GROUP BY m.id;
	`
	row := db.QueryRow(query, menuID)
	var categoriesIDs []uint8
	err = row.Scan(
		&menuView.ID,
		&menuView.Name,
		&menuView.IsEnabled,
		&categoriesIDs,
	)
	menuView.CategoriesIDs = convertUint8ToUUIDSlice(categoriesIDs)
	if err != nil {
		return MenuView{}, err
	}
	return menuView, nil
}

func (repo MenuRepository) GetAllMenus() ([]MenuView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return []MenuView{}, err
	}
	defer db.Close()

	var menuViews []MenuView

	query := `
		SELECT m.*, array_agg(mc.category_id) AS ids
		FROM menus m
		LEFT JOIN menus_categories mc ON m.id = mc.menu_id
		GROUP BY m.id;
	`
	rows, err := db.Query(query)
	defer rows.Close()

	for rows.Next() {
		var menuView MenuView
		var categoriesIDs []uint8
		err = rows.Scan(
			&menuView.ID,
			&menuView.Name,
			&menuView.IsEnabled,
			&categoriesIDs,
		)
		menuView.CategoriesIDs = convertUint8ToUUIDSlice(categoriesIDs)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		menuViews = append(menuViews, menuView)
	}

	return menuViews, nil
}

func (repo MenuRepository) DeleteMenu(menuID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM menus WHERE id=$1`
	_, err = db.Exec(query, menuID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) EnableMenu(menuID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE menus SET is_enabled=TRUE WHERE id=$1`
	_, err = db.Exec(query, menuID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) DisableMenu(menuID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE menus SET is_enabled=FALSE WHERE id=$1`
	_, err = db.Exec(query, menuID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) ChangeMenuName(menuID uuid.UUID, newName string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE menus SET name=$2 WHERE id=$1`
	_, err = db.Exec(query, menuID, newName)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) CreateCategory(categoryID uuid.UUID, categoryName string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO categories ("id", "name") VALUES ($1, $2)`
	_, err = db.Exec(query, categoryID, categoryName)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) DeleteCategory(categoryID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM categories WHERE id=$1`
	_, err = db.Exec(query, categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) GetCategory(categoryID uuid.UUID) (CategoryView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return CategoryView{}, err
	}
	defer db.Close()

	var categoryView CategoryView

	query := `
		SELECT m.*, array_agg(mc.subcategory_id) AS ids
		FROM categories m
		LEFT JOIN category_subcategories mc ON m.id = mc.category_id
		WHERE m.id=$1
		GROUP BY m.id;
	`
	row := db.QueryRow(query, categoryID)
	var subCategoriesIDs []uint8
	err = row.Scan(
		&categoryView.ID,
		&categoryView.Name,
		&subCategoriesIDs,
	)

	categoryView.SubCategoriesIDs = convertUint8ToUUIDSlice(subCategoriesIDs)

	if err != nil {
		return CategoryView{}, err
	}
	return categoryView, nil
}

func (repo MenuRepository) AddCategoryToMenu(menuID, categoryID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO menus_categories ("menu_id", "category_id") VALUES ($1, $2)`
	_, err = db.Exec(query, menuID, categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) RemoveCategoryFromMenu(menuID, categoryID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM menus_categories WHERE menu_id=$1 AND category_id=$2`
	_, err = db.Exec(query, menuID, categoryID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) GetCategoriesByIDs(categoriesIDs []uuid.UUID) ([]CategoryView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return []CategoryView{}, err
	}
	defer db.Close()

	idListString := makeStringList(categoriesIDs)

	query := `
		SELECT m.*, array_agg(mc.subcategory_id) AS ids
		FROM categories m
		LEFT JOIN category_subcategories mc ON m.id = mc.category_id
		WHERE id IN(` + idListString + `)
		GROUP BY m.id;
	`

	rows, err := db.Query(query)
	defer rows.Close()

	categories := []CategoryView{}

	for rows.Next() {
		var categoryView CategoryView
		var subCategoriesIDs []uint8
		err = rows.Scan(
			&categoryView.ID,
			&categoryView.Name,
			&subCategoriesIDs,
		)

		categoryView.SubCategoriesIDs = convertUint8ToUUIDSlice(subCategoriesIDs)

		if err != nil {
			return []CategoryView{}, err
		}
		categories = append(categories, categoryView)
	}

	return categories, nil
}

func (repo MenuRepository) ChangeCategoryName(categoryID uuid.UUID, newName string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE categories SET name=$2 WHERE id=$1`
	_, err = db.Exec(query, categoryID, newName)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) CreateSubCategory(subCategoryID uuid.UUID, subCategoryName string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO subcategories ("id", "name") VALUES ($1, $2)`
	_, err = db.Exec(query, subCategoryID, subCategoryName)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) DeleteSubCategory(subCategoryID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM subcategories WHERE id=$1`
	_, err = db.Exec(query, subCategoryID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) GetSubCategory(subCategoryID uuid.UUID) (SubCategoryView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return SubCategoryView{}, err
	}
	defer db.Close()

	var subCategoryView SubCategoryView

	query := `
		SELECT m.*, array_agg(mc.menuitem_id) AS ids
		FROM subcategories m
		LEFT JOIN subcategory_menuitems mc ON m.id = mc.subcategory_id
		WHERE m.id=$1
		GROUP BY m.id;
	`
	row := db.QueryRow(query, subCategoryID)
	var menuItemsIDs []uint8
	err = row.Scan(
		&subCategoryView.ID,
		&subCategoryView.Name,
		&menuItemsIDs,
	)

	subCategoryView.MenuItemsIDs = convertUint8ToUUIDSlice(menuItemsIDs)

	if err != nil {
		return SubCategoryView{}, err
	}
	return subCategoryView, nil
}

func (repo MenuRepository) GetSubCategoriesByIDs(subCategoriesIDs []uuid.UUID) ([]SubCategoryView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return []SubCategoryView{}, err
	}
	defer db.Close()

	idListString := makeStringList(subCategoriesIDs)

	query := `
		SELECT m.*, array_agg(mc.menuitem_id) AS ids
		FROM subcategories m
		LEFT JOIN subcategory_menuitems mc ON m.id = mc.subcategory_id
		WHERE id IN(` + idListString + `)
		GROUP BY m.id;
	`

	rows, err := db.Query(query)
	defer rows.Close()

	subCategories := []SubCategoryView{}

	for rows.Next() {
		var subCategoryView SubCategoryView
		var menuItemsIDs []uint8
		err = rows.Scan(
			&subCategoryView.ID,
			&subCategoryView.Name,
			&menuItemsIDs,
		)

		subCategoryView.MenuItemsIDs = convertUint8ToUUIDSlice(menuItemsIDs)

		if err != nil {
			return []SubCategoryView{}, err
		}
		subCategories = append(subCategories, subCategoryView)
	}

	return subCategories, nil
}

func (repo MenuRepository) RemoveSubCategoryFromCategory(categoryID, subCategoryID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM category_subcategories WHERE category_id=$1 AND subcategory_id=$2`
	_, err = db.Exec(query, categoryID, subCategoryID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) AddSubCategoryToCategory(categoryID, subCategoryID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO category_subcategories ("category_id", "subcategory_id") VALUES ($1, $2)`
	_, err = db.Exec(query, categoryID, subCategoryID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) CreateMenuItem(menuItemID uuid.UUID, menuItemName string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO menuitems ("id", "name") VALUES ($1, $2)`
	_, err = db.Exec(query, menuItemID, menuItemName)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) DeleteMenuItem(menuItemID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM menuitems WHERE id=$1`
	_, err = db.Exec(query, menuItemID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) GetMenuItem(menuItemID uuid.UUID) (MenuItemView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return MenuItemView{}, err
	}
	defer db.Close()

	var menuItemView MenuItemView

	query := `SELECT * FROM menuitems WHERE id=$1`
	row := db.QueryRow(query, menuItemID)

	err = row.Scan(
		&menuItemView.ID,
		&menuItemView.Name,
	)

	if err != nil {
		return MenuItemView{}, err
	}
	return menuItemView, nil
}

func (repo MenuRepository) GetMenuItemsByIDs(menuItemsIDs []uuid.UUID) ([]MenuItemView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return []MenuItemView{}, err
	}
	defer db.Close()

	idListString := makeStringList(menuItemsIDs)

	query := `
		SELECT *
		FROM menuitems
		WHERE id IN(` + idListString + `);
	`

	rows, err := db.Query(query)
	defer rows.Close()

	menuItems := []MenuItemView{}

	for rows.Next() {
		var menuItemView MenuItemView
		err = rows.Scan(
			&menuItemView.ID,
			&menuItemView.Name,
		)

		if err != nil {
			return []MenuItemView{}, err
		}
		menuItems = append(menuItems, menuItemView)
	}

	return menuItems, nil
}

func (repo MenuRepository) RemoveMenuItemFromSubCategory(subCategoryID, menuItemID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM subcategory_menuitems WHERE subcategory_id=$1 AND menuitem_id=$2`
	_, err = db.Exec(query, subCategoryID, menuItemID)
	if err != nil {
		return err
	}
	return nil
}

func (repo MenuRepository) AddMenuItemToSubCategory(subCategoryID, menuItemID uuid.UUID) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO subcategory_menuitems ("subcategory_id", "menuitem_id") VALUES ($1, $2)`
	_, err = db.Exec(query, subCategoryID, menuItemID)
	if err != nil {
		return err
	}
	return nil
}

// helpers

func convertUint8ToUUIDSlice(categoriesIDs []uint8) (res []uuid.UUID) {
	categoriesIDsString := string(categoriesIDs)
	if categoriesIDsString == "{NULL}" {
		return
	}
	categoriesString := strings.TrimPrefix(categoriesIDsString, "{")
	categoriesString = strings.TrimSuffix(categoriesString, "}")
	categories := strings.Split(string(categoriesString), ",")
	for _, v := range categories {
		id, _ := uuid.FromString(v)
		res = append(res, id)
	}
	return
}

func makeStringList(categoriesIDs []uuid.UUID) string {
	idListString := ""
	for _, v := range categoriesIDs {
		idListString = idListString + "'" + v.String() + "',"
	}
	if len(idListString) > 0 {
		idListString = strings.TrimSuffix(idListString, ",")
	}
	return idListString
}
