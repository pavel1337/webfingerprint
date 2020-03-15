package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pavel1337/webfingerprint/pkg/models"
)

type application struct {
	flags    Flags
	errorLog *log.Logger
	infoLog  *log.Logger
	db       *gorm.DB
}

func main() {
	f := ParseFlags()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	db, err := models.Init(f.Dsn)
	if err != nil {
		errorLog.Println(err)
		return
	}
	defer db.Close()

	app := &application{
		flags:    f,
		errorLog: errorLog,
		infoLog:  infoLog,
		db:       db,
	}

	if f.Mode == "prepare" {
		app.Prepare()
	} else if f.Mode == "capture" {
		app.Pcaps()
	} else if f.Mode == "visualize" {
		app.Visualize()
	} else if f.Mode == "export" {
		app.Export()
	} else if f.Mode == "learn" {
		app.Learn()
	} else {
		Usage()
	}
}
