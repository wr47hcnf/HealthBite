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

func registerUser(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/register.tmpl", "static/header.tmpl", "static/navbar.tmpl"))
	pageData := PageData{
		PageTitle: "Register",
	}
	cookie, err := r.Cookie("session_cookie")
	if err == nil {
		parseCookie(cookie, &pageData.UserInfo)
		pageData.PageError = append(pageData.PageError, Error{
			ErrorCode:    3,
			ErrorMessage: fmt.Sprintf("You are already logged in as %s", pageData.UserInfo.Username),
		})
	}
	if r.Method == http.MethodPost {
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

	tmpl.Execute(w, pageData)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/login.tmpl", "static/header.tmpl", "static/navbar.tmpl"))
	pageData := PageData{
		PageTitle: "Register",
	}
	cookie, err := r.Cookie("session_cookie")
	if err == nil {
		parseCookie(cookie, &pageData.UserInfo)
		pageData.PageError = append(pageData.PageError, Error{
			ErrorCode:    3,
			ErrorMessage: fmt.Sprintf("You are already logged in as %s", pageData.UserInfo.Username),
		})
		tmpl.Execute(w, pageData)
		return
	}
	if r.Method == http.MethodPost {
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

		row := Db.QueryRow("SELECT id, password FROM users WHERE username = $1", username)

		var id uuid.UUID
		var hashedpassword string
		err = row.Scan(&id, &hashedpassword)
		if err != nil {

		}
		/*for rows.Next() {
			err := rows.Scan(&id, &hashedpassword)
			if err != nil {
				http.Error(w, "S-a intampinat o problema in citirea bazei de date", http.StatusInternalServerError)
				log.Fatal("Failed to read the login database: ", err)
			}
			if bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password)) == nil {
				expiration := time.Now().Add(24 * time.Hour)
				cookie := http.Cookie{
					Name:     "session_cookie",
					Value:    fmt.Sprint(id),
					Expires:  expiration,
					Path:     "/",
					Secure:   true,
					HttpOnly: true,
				}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} */
		// TODO: Incorrect password logic
	}
	tmpl.Execute(w, pageData)
}
