package views

import (
	"database/sql"
	"log"

	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

type MenuView struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ViewRepository struct {
	connectionString string
}

func NewViewRepository(connectionString string) ViewRepository {
	return ViewRepository{
		connectionString: connectionString,
	}
}

func (repo ViewRepository) CreateMenu(menuID uuid.UUID, menuName string) error {
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

func (repo ViewRepository) GetMenu(menuID uuid.UUID) (MenuView, error) {
	db, err := sql.Open("postgres", repo.connectionString)
	if err != nil {
		return MenuView{}, err
	}
	defer db.Close()

	var menuView MenuView

	query := `SELECT * FROM menus WHERE id=$1`
	row := db.QueryRow(query, menuID)

	err = row.Scan(&menuView.ID, &menuView.Name)

	if err != nil {
		return MenuView{}, err
	}
	return menuView, nil
}

func (repo ViewRepository) GetAllMenus() ([]MenuView, error) {
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
		err = rows.Scan(&menuView.ID, &menuView.Name)
		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		menuViews = append(menuViews, menuView)
	}

	return menuViews, nil
}

func (repo ViewRepository) DeleteMenu(menuID uuid.UUID) error {
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
