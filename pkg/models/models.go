package models

import "github.com/jinzhu/gorm"

type Website struct {
	gorm.Model
	Hostname string `gorm:"unique"`
	Subs     []Sub
}

type Sub struct {
	gorm.Model
	Link      string `gorm:"unique"`
	WebsiteID uint
	Pcaps     []Pcap
}

type Pcap struct {
	gorm.Model
	Path    string `gorm:"unique"`
	SubID   uint
	Proxy   string
	BCumul  []byte `gorm:"size:32000"`
	Outlier bool
}

func Init(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Migrate the schema
	db.AutoMigrate(&Website{})
	db.AutoMigrate(&Sub{})
	db.AutoMigrate(&Pcap{})
	return db, nil
}

func DistinctProxiesBySub(id uint, db *gorm.DB) ([]string, error) {
	var proxies []string
	rows, err := db.Raw("SELECT DISTINCT proxy FROM pcaps WHERE sub_id = ?", id).Rows()
	if err != nil {
		return proxies, err
	}
	defer rows.Close()
	for rows.Next() {
		var proxy string
		err := rows.Scan(&proxy)
		if err != nil {
			return proxies, err
		}
		proxies = append(proxies, proxy)
	}
	return proxies, nil
}
