package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func parseCookie(cookie *http.Cookie, userdata *User) error {
	if cookie == nil {
		*userdata = User{
			IsLogged: false,
		}
		return nil

	} else {
		uid := cookie.Value
		row := Db.QueryRow("SELECT username FROM users WHERE id = $1", uid)

		var username string
		err := row.Scan(&username)

		if err != nil {
			return err
		}

		parsedUUID, err := uuid.Parse(uid)

		if err != nil {
			return err
		}

		*userdata = User{
			IsLogged: true,
			Username: username,
			ID:       parsedUUID,
		}
	}
	return nil
}

func profilePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"static/profile_page.tmpl",
		"static/error.tmpl",
		"static/header.tmpl",
		"static/navbar.tmpl",
		"static/footer.tmpl",
	))
	pageData := PageData{
		PageTitle: "Profile",
	}
	cookie, err := r.Cookie("session_cookie")
	if err != nil {
		pageData.PageError = append(pageData.PageError, Error{
			ErrorCode:    3,
			ErrorMessage: "You must be logged in!",
		})
		tmpl.Execute(w, pageData)
		return
	}
	err = parseCookie(cookie, &pageData.UserInfo)
	if err != nil {
		log.Printf("Failed to parse cookie for %s", r.RemoteAddr)
	}
	tmpl.Execute(w, pageData)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"static/add_product.tmpl",
		"static/error.tmpl",
		"static/header.tmpl",
		"static/navbar.tmpl",
		"static/footer.tmpl",
	))
	pageData := PageData{
		PageTitle: "Add product",
	}
	cookie, err := r.Cookie("session_cookie")
	if err == nil {
		err = parseCookie(cookie, &pageData.UserInfo)
		if err != nil {
			log.Printf("Failed to parse cookie for %s", r.RemoteAddr)
		}
	}
	tmpl.Execute(w, pageData)
}
