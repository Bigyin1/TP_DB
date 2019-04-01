package db

import (
	"database/sql"
	"fmt"

	//
	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func Init() (db *Database, err error) {

	db = &Database{
		DB: nil,
	}
	dbCredentials := "user=docker password=docker dbname=docker sslmode=disable"

	if db.DB, err = sql.Open("postgres", dbCredentials); err != nil {
		fmt.Println("db/Init cant open:" + err.Error())
		return
	}

	db.DB.SetMaxOpenConns(50)

	if err = db.DB.Ping(); err != nil {
		fmt.Println("db/Init cant access:" + err.Error())
		return
	}

	return
}
