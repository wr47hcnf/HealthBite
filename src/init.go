package main

import (
	"fmt"
	"net/http"
	"time"
)

func init() {

	fmt.Printf("HealthBite backend server\n2024 Patrick Covaci\n%s\n", time.Now())

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}
