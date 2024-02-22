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
		http.Error(w, "Eroare procesare formular inregistrare: ", http.StatusBadRequest)
		log.Print(fmt.Sprintf("Error parsing registration form for %s: %s ", r.RemoteAddr, err))
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if usernameRegex := regexp.MustCompile(`[a-zA-Z0-9_]{3,20}$`); !usernameRegex.MatchString(username) {
		http.Error(w, "Ati introdus un username invalid!", http.StatusBadRequest)
		log.Printf("%s attempted to enter an illegal username", r.RemoteAddr)
		return
	}

	if passwordRegex := regexp.MustCompile(`.{6,30}$`); !passwordRegex.MatchString(password) {
		http.Error(w, "Ati introdus o parola invalida!", http.StatusBadRequest)
		log.Printf("%s attempted to enter an illegal password", r.RemoteAddr)
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Parola nu a putut fi procesata!", http.StatusInternalServerError)
		log.Print(fmt.Sprintf("Failed to hash password for session %s", r.RemoteAddr))
		return
	}

	userID, err := uuid.NewRandom()

	_, err = Db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", userID, username, hashedpassword)

	if err != nil {
		http.Error(w, "Credentialele nu au putut fi inserate in baza de date!", http.StatusInternalServerError)
		log.Printf("Could not insert user data in the database for %s: %s", r.RemoteAddr, err)
		return
	}

	fmt.Fprintf(w, "User ID:%s\nSalted password: %s\n", userID, hashedpassword)
}