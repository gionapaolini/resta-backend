package internal

import (
	"database/sql"
	"log"

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

	query := `SELECT * FROM menus WHERE id=$1`
	row := db.QueryRow(query, menuID)

	err = row.Scan(
		&menuView.ID,
		&menuView.Name,
		&menuView.IsEnabled,
	)

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

	query := `SELECT * FROM menus`
	rows, err := db.Query(query)
	defer rows.Close()

	for rows.Next() {
		var menuView MenuView
		err = rows.Scan(
			&menuView.ID,
			&menuView.Name,
			&menuView.IsEnabled,
		)
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

func (repo MenuRepository) CreateCategory(categoryID uuid.UUID, categoryName, imageURL string) error {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	query := `INSERT INTO categories ("id", "name", "image_url") VALUES ($1, $2, $3)`
	_, err = db.Exec(query, categoryID, categoryName, imageURL)
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

	query := `SELECT * FROM categories WHERE id=$1`
	row := db.QueryRow(query, categoryID)

	err = row.Scan(
		&categoryView.ID,
		&categoryView.Name,
		&categoryView.ImageURL,
	)

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

func (repo MenuRepository) GetMenuCategoriesIDs(menuID uuid.UUID) ([]uuid.UUID, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return []uuid.UUID{}, err
	}
	defer db.Close()

	var categoriesIDs []uuid.UUID

	query := `SELECT category_id FROM menus_categories WHERE menu_id=$1`
	rows, err := db.Query(query, menuID)
	defer rows.Close()

	for rows.Next() {
		var categoryID uuid.UUID
		err = rows.Scan(&categoryID)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		categoriesIDs = append(categoriesIDs, categoryID)
	}

	return categoriesIDs, nil
}
