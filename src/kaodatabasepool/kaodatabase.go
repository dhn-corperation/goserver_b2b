package databasepool

import (
	"fmt"
	"database/sql"
	// config "mycs/src/kaoconfig"
	"log"

	// _ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDatabase() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "210.114.225.58", "5432", "postgres", "dhn7985!", "kakao")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	DB = db

}
