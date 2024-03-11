package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

const (
	host     = "healthbite-db.c7c68asau12b.eu-north-1.rds.amazonaws.com"
	port     = 5432
	user     = "healthbite"
	password = "healthbite"
	dbname   = "healthbite"
)

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

	if tableExists == false {
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
	if tableExists == false {
		_, err = Db.Exec(`
			CREATE TABLE userdata (
				uid UUID PRIMARY KEY,
				fname VARCHAR(10),
				lname VARCHAR(10),
				age INT,
				profilepic VARCHAR(50),
				FOREIGN KEY (uid) REFERENCES users(id)
			)
		`)

		if err != nil {
			log.Fatal("Failed to initialize 'userdata' table", err)
			return err
		}

		log.Println("Succesfully initialized 'userdata' table")
	}

	// Check 'productdata'
	tableExists = false
	err = Db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'productdata')").Scan(&tableExists)
	if tableExists == false {
		_, err = Db.Exec(`
			CREATE TABLE productdata (
				prodid UUID PRIMARY KEY,
				barcode VARCHAR(20),
				name VARCHAR(30),
				brand VARCHAR(20),
				pic VARCHAR(50),
				weight VARCHAR(10),
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
	if tableExists == false {
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
					prodid UUID KEY,
					userid UUID KEY,
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
	*/
	return nil
}
