package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	defer Db.Close()

	http.HandleFunc("/register", registerUser)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		pageData := PageData{
			PageTitle: "Homepage",
		}

		cookie, err := r.Cookie("session_cookie")

		if cookie != nil {
			err := parseCookie(cookie, &pageData.UserInfo)
			if err != nil {
				log.Printf("Failed to parse cookie for %s: %s", r.RemoteAddr, err)
				pageData.PageError = append(pageData.PageError, Error{
					ErrorCode:    1,
					ErrorMessage: "failed to parse cookie",
				})
			}
		}

		tmpl, err := template.ParseFiles("static/index.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, pageData)

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
