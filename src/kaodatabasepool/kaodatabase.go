package databasepool

import (
	"database/sql"
	"log"
	config "mycs/src/kaoconfig"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDatabase() {
	db, err := sql.Open(config.Conf.DB, config.Conf.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	DB = db

}
