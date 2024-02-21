package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

func registrationPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/register.html"))
	tmpl.Execute(w, nil)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal("Error parsing registration form for %s: %s ", r.RemoteAddr, err)
		http.Error(w, "Error parsing registration form: ", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	usernameRegex := regexp.MustCompile(`[a-zA-Z0-9_]{3,20}$`)
	passwordRegex := regexp.MustCompile(`.{6,30}$`)

	if !usernameRegex.MatchString(username) {
		log.Fatal("%s attempted to enter an illegal username")
		http.Error(w,)
	}

}
