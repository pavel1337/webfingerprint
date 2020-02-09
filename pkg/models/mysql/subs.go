package mysql

import (
	"database/sql"

	"github.com/pavel1337/webfingerprint/pkg/models"
)

type SubModel struct {
	DB *sql.DB
}

func (m *SubModel) Insert(title string, websiteid int) (int, error) {
	stmt := `INSERT INTO subs (title, websiteid) VALUES (?, ?)`

	result, err := m.DB.Exec(stmt, title, websiteid)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SubModel) List() ([]*models.Sub, error) {
	stmt := `SELECT id, title, websiteid FROM subs ORDER BY title`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := []*models.Sub{}

	for rows.Next() {
		s := &models.Sub{}

		err = rows.Scan(&s.ID, &s.Title, &s.WebsiteId)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (m *SubModel) GetById(subid int) (*models.Sub, error) {
	stmt := `SELECT id, title, websiteid FROM subs WHERE id = ?`

	row := m.DB.QueryRow(stmt, subid)
	s := &models.Sub{}
	err := row.Scan(&s.ID, &s.Title, &s.WebsiteId)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

func (m *SubModel) GetByWebsiteId(websiteid int) ([]*models.Sub, error) {
	stmt := `SELECT id, title, websiteid FROM subs WHERE websiteid = ?`
	rows, err := m.DB.Query(stmt, websiteid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := []*models.Sub{}

	for rows.Next() {
		s := &models.Sub{}

		err = rows.Scan(&s.ID, &s.Title, &s.WebsiteId)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
