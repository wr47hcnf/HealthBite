package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "healthbite"
	password = "healthbite"
	dbname   = "healthbite"
)

func dbConnect() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		return err
	}
	Db = db
	return nil
}

func dbInit() error {
	var tableExists bool
	err := Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)

	if err != nil {
		log.Fatal("Failed to query the database", err)
		return err
	}

	if tableExists {
		log.Println("Database already initialized, continuing")
		return nil
	}

	_, err = Db.Exec(`
	CREATE TABLE users (
		id UUID PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL)
   `)

	if err != nil {
		log.Fatal("Failed to initialize database", err)
		return err
	}

	log.Println("Succesfully initialized database")
	return nil
}
