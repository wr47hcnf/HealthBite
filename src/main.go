package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	defer Db.Close()

	http.HandleFunc("/inregistrare", registrationPage)
	http.HandleFunc("/registeruser", registerUser)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_cookie")

		if err == nil {

			userID := cookie.Value

			if userID == "" {
				fmt.Fprintln(w, "No cookie")
			}
		}

		tmpl, err := template.ParseFiles("static/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, time.Now().Format(time.TimeOnly))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}
