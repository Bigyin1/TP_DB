package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"

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

	query, err := ioutil.ReadFile("scheme.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(string(query))
	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}

	query, err = ioutil.ReadFile("init.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(string(query))
	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}

	return
}
