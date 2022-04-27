package internal

import (
	"database/sql"
	"log"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

type MenuView struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsEnabled bool      `json:"isEnabled"`
}

type IMenuRepository interface {
	CreateMenu(menuID uuid.UUID, menuName string) error
	GetMenu(menuID uuid.UUID) (MenuView, error)
	GetAllMenus() ([]MenuView, error)
	DeleteMenu(menuID uuid.UUID) error
	EnableMenu(menuID uuid.UUID) error
	DisableMenu(menuID uuid.UUID) error
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

	query := `UPDATE menus SET isEnabled=TRUE WHERE id=$1`
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

	query := `UPDATE menus SET isEnabled=FALSE WHERE id=$1`
	_, err = db.Exec(query, menuID)
	if err != nil {
		return err
	}
	return nil
}
