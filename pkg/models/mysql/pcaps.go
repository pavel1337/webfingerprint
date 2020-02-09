package mysql

import (
	"database/sql"

	"github.com/pavel1337/webfingerprint/pkg/models"
)

type PcapModel struct {
	DB *sql.DB
}

func (m *PcapModel) Insert(path string, subid int, proxy string) (int, error) {
	stmt := `INSERT INTO pcaps (path, subid, proxy) VALUES (?, ?, ?)`

	result, err := m.DB.Exec(stmt, path, subid, proxy)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *PcapModel) List() ([]*models.Pcap, error) {
	stmt := `SELECT id, path, subid, proxy FROM pcaps ORDER BY subid, proxy`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy)
		if err != nil {
			return nil, err
		}
		pcaps = append(pcaps, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pcaps, nil
}

func (m *PcapModel) GetBySubId(subid int) ([]*models.Pcap, error) {
	stmt := `SELECT id, path, subid, proxy FROM pcaps WHERE subid = ?`
	rows, err := m.DB.Query(stmt, subid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy)
		if err != nil {
			return nil, err
		}
		pcaps = append(pcaps, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pcaps, nil
}

func (m *PcapModel) GetBySubIdAndProxy(subid int, proxy string) ([]*models.Pcap, error) {
	stmt := `SELECT id, path, subid, proxy FROM pcaps WHERE subid = ? AND proxy = ?`
	rows, err := m.DB.Query(stmt, subid, proxy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy)
		if err != nil {
			return nil, err
		}
		pcaps = append(pcaps, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pcaps, nil
}
