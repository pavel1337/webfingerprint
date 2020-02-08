package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pavel1337/webfingerprint/pkg/models/mysql"
)

type application struct {
	config   Config
	errorLog *log.Logger
	infoLog  *log.Logger
	websites *mysql.WebsiteModel
	subs     *mysql.SubModel
}

func main() {
	c := ParseConfig()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(c.Dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := &application{
		config:   c,
		errorLog: errorLog,
		infoLog:  infoLog,
		websites: &mysql.WebsiteModel{DB: db},
		subs:     &mysql.SubModel{DB: db},
	}
	if c.Mode == "prepare" {
		app.Websites()
		app.Subs()
	} else {
		Usage()
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
