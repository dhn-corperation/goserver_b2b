package databasepool

import (
	"database/sql"
	config "mycs/src/kaoconfig"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() {
	db, err := sql.Open(config.Conf.DB, config.Conf.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	DB = db

}
