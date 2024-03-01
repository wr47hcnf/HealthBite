package main

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)

func registerUser(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/register.tmpl", "static/error.tmpl", "static/header.tmpl", "static/navbar.tmpl"))
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
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Invalid username!"),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s attempted to enter an illegal username", r.RemoteAddr)
			return
		}

		if passwordRegex := regexp.MustCompile(`.{6,30}$`); !passwordRegex.MatchString(password) {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Invalid password!"),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s attempted to enter an illegal password", r.RemoteAddr)
			return
		}

		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			http.Error(w, "Cannot process password!", http.StatusInternalServerError)
			log.Print(fmt.Sprintf("Failed to hash password for session %s", r.RemoteAddr))
			return
		}

		userID, err := uuid.NewRandom()

		dbexec, err := Db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", userID, username, hashedpassword)
		log.Print(dbexec)

		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode: 1,
				ErrorMessage: fmt.Sprintf(
					"Credentials cannot be inserted into the database! Perhaps you used a username already in use!"),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			log.Printf("Could not insert user data in the database for %s: %s", r.RemoteAddr, err)
			return
		}

		pageData.PageError = append(pageData.PageError, Error{
			ErrorCode: 4,
			ErrorMessage: fmt.Sprintf(
				"Succesfully created account %s, redirecting to the login page shortly...", username),
		})
		log.Printf("Created account %s! User ID:%s\nSalted password: %s\n", username, userID, hashedpassword)
		err = tmpl.Execute(w, pageData)
		if err != nil {
			log.Print(err)
		}
		time.Sleep(5 * time.Second)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Print(err)
	}
	return
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/login.tmpl", "static/header.tmpl", "static/navbar.tmpl", "static/error.tmpl"))
	pageData := PageData{
		PageTitle: "Login",
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
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Error parsing registration form!"),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			log.Print(fmt.Sprintf("Error parsing registration form for %s: %s ", r.RemoteAddr, err))
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		if usernameRegex := regexp.MustCompile(`[a-zA-Z0-9_]{3,20}$`); !usernameRegex.MatchString(username) {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Invalid username!"),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s attempted to enter an illegal username", r.RemoteAddr)
			return
		}

		if passwordRegex := regexp.MustCompile(`.{6,30}$`); !passwordRegex.MatchString(password) {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Invalid password!"),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s attempted to enter an illegal password", r.RemoteAddr)
			return
		}

		row := Db.QueryRow("SELECT id, password FROM users WHERE username = $1", username)

		var id uuid.UUID
		var hashedpassword string
		err = row.Scan(&id, &hashedpassword)
		if err != nil {
			log.Print(err)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("The username %s could not be found!", username),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
		if err != nil {
			log.Print(err)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Invalid password for username %s", username),
			})
			err = tmpl.Execute(w, pageData)
			if err != nil {
				log.Print(err)
			}
			return
		}
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
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Print(err)
	}
	return
}
