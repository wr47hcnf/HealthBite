package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
	fmt.Printf("HealthBite backend server\n(C) 2024 Patrick Covaci a.k.a Ty3r0X\n%s\n", time.Now())

	err := dbConnect()

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
