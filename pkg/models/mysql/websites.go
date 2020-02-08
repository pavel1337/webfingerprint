package mysql

import (
	"database/sql"

	"github.com/pavel1337/webfingerprint/pkg/models"
)

type WebsiteModel struct {
	DB *sql.DB
}

func (m *WebsiteModel) Insert(title string) (int, error) {
	stmt := `INSERT INTO websites (title) VALUES (?)`

	result, err := m.DB.Exec(stmt, title)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *WebsiteModel) List() ([]*models.Website, error) {
	stmt := `SELECT id, title FROM websites ORDER BY title`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	websites := []*models.Website{}

	for rows.Next() {
		w := &models.Website{}

		err = rows.Scan(&w.ID, &w.Title)
		if err != nil {
			return nil, err
		}
		websites = append(websites, w)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return websites, nil
}

func (m *WebsiteModel) GetByTitle(title string) (*models.Website, error) {
	stmt := `SELECT id, title FROM websites WHERE title = ?`
	row := m.DB.QueryRow(stmt, title)
	w := &models.Website{}
	err := row.Scan(&w.ID, &w.Title)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return w, nil
}
