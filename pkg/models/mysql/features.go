package mysql

import (
	"database/sql"
	"encoding/json"

	"github.com/pavel1337/webfingerprint/pkg/models"
)

type FeatureModel struct {
	DB *sql.DB
}

func (m *FeatureModel) Insert(pcapid int, featureset [50]int) (int, error) {
	stmt := `INSERT INTO features (pcapid, featureset) VALUES (?, ?)`

	fs := models.FeatureSetJson{}
	fs.Cumul = featureset

	b, err := json.Marshal(fs)
	if err != nil {
		return 0, err
	}

	result, err := m.DB.Exec(stmt, pcapid, b)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *FeatureModel) List() ([]*models.Feature, error) {
	stmt := `SELECT id, pcapid, featureset FROM features ORDER BY pcapid`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	features := []*models.Feature{}

	for rows.Next() {
		f := &models.Feature{}

		err = rows.Scan(&f.ID, &f.PcapID, &f.FeatureSet)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(f.FeatureSet, &f.FeatureSetJson)
		if err != nil {
			return nil, err
		}
		features = append(features, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return features, nil
}
