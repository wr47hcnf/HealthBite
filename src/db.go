package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

const (
	host     = "172.31.26.63"
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
	// Check 'users'
	var tableExists bool
	err := Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)

	if err != nil {
		log.Fatal("Failed to query the database", err)
		return err
	}

	if tableExists == false {
		_, err = Db.Exec(`
			CREATE TABLE users (
				id UUID PRIMARY KEY,
				username VARCHAR(50) UNIQUE NOT NULL,
				password VARCHAR(100) NOT NULL
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'users' database", err)
			return err
		}

		log.Println("Succesfully initialized 'user' database")
	}

	// Check 'userdata'
	tableExists = false
	err = Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if tableExists == false {
		_, err = Db.Exec(`
			CREATE TABLE userdata (
				uid UUID PRIMARY KEY,
				fname VARCHAR(10),
				lname VARCHAR(10),
				age INT,
				profilepic VARCHAR(10),
				FOREIGN KEY (uid) REFERENCES users(id)
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'userdata' database", err)
			return err
		}

		log.Println("Succesfully initialized 'userdata' database")
	}
	return nil
}
