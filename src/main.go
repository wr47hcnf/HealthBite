package main

import (
	"html/template"
	"net/http"
	"time"
)

func main() {
	timp_curent := time.Now()

	defer Db.Close()

	http.HandleFunc("/inregistrare", registrationPage)
	http.HandleFunc("/registeruser", registerUser)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("static/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, timp_curent.Format(time.TimeOnly))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
