package main

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "healthbite"
	password = "healthbite"
	dbname   = "healthbite"
)

var Db *sql.DB

func init_db() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	Db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer Db.Close()

	err = Db.Ping()

	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		panic(err)
	}
}
