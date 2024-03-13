package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func dbConnect() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
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

	if !tableExists {
		_, err = Db.Exec(`
			CREATE TABLE users (
				id UUID PRIMARY KEY,
				username VARCHAR(50) UNIQUE NOT NULL,
				password VARCHAR(100) NOT NULL
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'users' table ", err)
			return err
		}

		log.Println("Succesfully initialized 'user' table")
	}

	// Check 'userdata'
	tableExists = false
	err = Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'userdata')").Scan(&tableExists)
	if err != nil {
		log.Fatal("Failed to query db: ", err)
	}
	if !tableExists {
		_, err = Db.Exec(`
			CREATE TABLE userdata (
				uid UUID PRIMARY KEY,
				fname VARCHAR(20),
				lname VARCHAR(20),
				age INT,
				profilepic VARCHAR(80),
				FOREIGN KEY (uid) REFERENCES users(id)
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'userdata' table", err)
			return err
		}

		log.Println("Succesfully initialized 'userdata' table")
	}

	// Check 'nutritional_info' data type
	tableExists = false
	err = Db.QueryRow("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'nutritional_info')").Scan(&tableExists)
	if err != nil {
		log.Fatal("Failed to query db: ", err)
	}
	if !tableExists {
		_, err = Db.Exec(`
			CREATE TYPE nutritional_info AS (
    				name TEXT,
    				value TEXT
			)
		`)
		if err != nil {
			log.Fatal("Failed to initialize 'nutritional_info' data type ", err)
		}
		log.Println("Initialized nutritional_info data type")
	}

	// Check 'productdata'
	tableExists = false
	err = Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'productdata')").Scan(&tableExists)
	if err != nil {
		log.Fatal("Failed to query db: ", err)
	}
	if !tableExists {
		_, err = Db.Exec(`
			CREATE TABLE productdata (
				prodid UUID PRIMARY KEY,
				barcode VARCHAR(20),
				name VARCHAR(30),
				brand VARCHAR(20),
				pic VARCHAR(50),
				weight VARCHAR(10),
				calories INT,
				nutritional_info nutritional_info[],
				additives VARCHAR(20)[],
				allergens VARCHAR(20)[]
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'productdata' table", err)
			return err
		}

		log.Println("Succesfully initialized 'productdata' table")
	}

	// Check 'productreview'
	tableExists = false
	err = Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'productreview')").Scan(&tableExists)
	if err != nil {
		log.Fatal("Failed to query db: ", err)
	}
	if !tableExists {
		_, err = Db.Exec(`
			CREATE TABLE productreview (
				prodid UUID,
				userid UUID,
				message VARCHAR(200),
				score INT,
				date DATE,
				FOREIGN KEY (prodid) REFERENCES productdata(prodid),
				FOREIGN KEY (userid) REFERENCES users(id)
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'productreview' table", err)
			return err
		}

		log.Println("Succesfully initialized 'productreview' table")
	}

	// Check 'userallergens'
	//
	/*
		tableExists = false
		err = Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'userallergens')").Scan(&tableExists)
		if tableExists == false {
			_, err = Db.Exec(`
				CREATE TABLE productreview (
					userid UUID KEY,
					allergen VARCHAR(20)
				)
			`)

			if err != nil {
				log.Fatal("Failed to initialize 'userallergens' table", err)
				return err
			}

			log.Println("Succesfully initialized 'userallergens' table")
		}
	*/
	return nil
}
