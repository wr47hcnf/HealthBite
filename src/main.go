package main

import (
	"log"
	"net/http"
)

func main() {
	defer Db.Close()

	http.HandleFunc("/register", registerUser)
	http.HandleFunc("/login", loginUser)
	http.HandleFunc("/logout", logoutUser)
	http.HandleFunc("/profile", profilePage)
	http.HandleFunc("/addproduct", addProduct)
	http.HandleFunc("/product", viewProduct)

	http.HandleFunc("/", homePage)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
