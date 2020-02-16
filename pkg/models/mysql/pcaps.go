package mysql

import (
	"database/sql"
	"encoding/json"

	"github.com/pavel1337/webfingerprint/pkg/models"
)

type PcapModel struct {
	DB *sql.DB
}

func (m *PcapModel) Insert(path string, subid int, proxy string, cumul [50]int) (int, error) {
	stmt := `INSERT INTO pcaps (path, subid, proxy, bcumul) VALUES (?, ?, ?, ?)`

	pm := models.Pcap{}
	pm.Cumul = cumul

	bcumul, err := json.Marshal(pm.Cumul)
	if err != nil {
		return 0, err
	}

	result, err := m.DB.Exec(stmt, path, subid, proxy, bcumul)
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
	stmt := `SELECT id, path, subid, proxy, bcumul, outlier FROM pcaps ORDER BY subid, proxy`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy, &p.BCumul, &p.Outlier)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(p.BCumul, &p.Cumul)
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

func (m *PcapModel) SetOutlierById(id int) error {
	stmt := `UPDATE pcaps SET outlier = true WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err == sql.ErrNoRows {
		return models.ErrNoRecord
	} else if err != nil {
		return err
	}
	return nil
}

func (m *PcapModel) UnsetOutlierById(id int) error {
	stmt := `UPDATE pcaps SET outlier = false WHERE id = ?`

	_, err := m.DB.Exec(stmt, id)
	if err == sql.ErrNoRows {
		return models.ErrNoRecord
	} else if err != nil {
		return err
	}
	return nil
}

func (m *PcapModel) GetBySubId(subid int) ([]*models.Pcap, error) {
	stmt := `SELECT id, path, subid, proxy, bcumul, outlier FROM pcaps WHERE subid = ?`
	rows, err := m.DB.Query(stmt, subid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy, &p.BCumul, &p.Outlier)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(p.BCumul, &p.Cumul)
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
	stmt := `SELECT id, path, subid, proxy, bcumul, outlier FROM pcaps WHERE subid = ? AND proxy = ?`
	rows, err := m.DB.Query(stmt, subid, proxy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy, &p.BCumul, &p.Outlier)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(p.BCumul, &p.Cumul)
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

func (m *PcapModel) ListProxiesBySubid(subid int) ([]string, error) {
	stmt := `SELECT DISTINCT proxy FROM pcaps WHERE subid = ?`
	rows, err := m.DB.Query(stmt, subid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	proxies := []string{}

	for rows.Next() {
		var proxy string
		err = rows.Scan(&proxy)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

func (m *PcapModel) ListProxies() ([]string, error) {
	stmt := `SELECT DISTINCT proxy FROM pcaps`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	proxies := []string{}

	for rows.Next() {
		var proxy string
		err = rows.Scan(&proxy)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, proxy)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

// SELECT DISTINCT proxy FROM pcaps;

func (m *PcapModel) GetBySubIdAndProxyAndNotOutlier(subid int, proxy string) ([]*models.Pcap, error) {
	stmt := `SELECT id, path, subid, proxy, bcumul, outlier FROM pcaps WHERE subid = ? AND proxy = ? AND outlier = false`
	rows, err := m.DB.Query(stmt, subid, proxy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pcaps := []*models.Pcap{}

	for rows.Next() {
		p := &models.Pcap{}

		err = rows.Scan(&p.ID, &p.Path, &p.SubId, &p.Proxy, &p.BCumul, &p.Outlier)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(p.BCumul, &p.Cumul)
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
