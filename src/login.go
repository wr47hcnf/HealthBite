package main

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
		log.Fatal(fmt.Sprintf("Error parsing registration form for %s: %s ", r.RemoteAddr, err))
		http.Error(w, "Eroare procesare formular inregistrare: ", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if usernameRegex := regexp.MustCompile(`[a-zA-Z0-9_]{3,20}$`); !usernameRegex.MatchString(username) {
		log.Printf("%s attempted to enter an illegal username", r.RemoteAddr)
		http.Error(w, "Ati introdus un username invalid!", http.StatusBadRequest)
	}

	if passwordRegex := regexp.MustCompile(`.{6,30}$`); !passwordRegex.MatchString(password) {
		log.Printf("%s attempted to enter an illegal password", r.RemoteAddr)
		http.Error(w, "Ati introdus o parola invalida!", http.StatusBadRequest)
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to hash password for session %s", r.RemoteAddr))
		http.Error(w, "Parola nu a putut fi procesata!", http.StatusInternalServerError)
	}

	userID, err := uuid.NewRandom()

	_, err = Db.Exec("INSERT INTO users (user_id, username, password) VALUES ($1, $2, $3)", userID, username, hashedpassword)

	fmt.Fprintf(w, "Salted password: %s\n", hashedpassword)
}
