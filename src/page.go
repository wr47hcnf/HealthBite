package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func parseCookie(cookie *http.Cookie, userdata *User) error {
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
		pageData.PageError = append(pageData.PageError, Error{
			ErrorCode:    1,
			ErrorMessage: "Could not parse the cookie!",
		})
		tmpl.Execute(w, pageData)
		return
	}
	if r.Method == http.MethodPost {
		err = r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			log.Printf("Failed to parse profile form for %s", r.RemoteAddr)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pfpFile, pfpHeader, err := r.FormFile("profilepic")
		if err != nil {
			log.Printf("Failed to parse pfp for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Could not parse profile pic!",
			})
			tmpl.Execute(w, pageData)
			return
		}
		pfpext := filepath.Ext(pfpHeader.Filename)
		key := "user/" + pageData.UserInfo.ID.String() + pfpext
		_, err = Svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   pfpFile,
		})
		if err != nil {
			log.Printf("Failed to upload pfp to aws for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Could not update profile pic!",
			})
			tmpl.Execute(w, pageData)
			return
		}
		s3url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, aws_region, key)
		_, err = Db.Exec("UPDATE userdata SET profilepic = $1 WHERE uid = $2", s3url, pageData.UserInfo.ID)
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
	if r.Method == http.MethodPost {
		err = r.ParseMultipartForm(10 << 20) // 10MB
		if err != nil {
			log.Printf("Failed to parse product form for %s", r.RemoteAddr)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pfpFile, pfpHeader, err := r.FormFile("profilepic")
		if err != nil {
			log.Printf("Failed to parse pfp for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Could not parse profile pic!",
			})
			tmpl.Execute(w, pageData)
			return
		}
		pfpext := filepath.Ext(pfpHeader.Filename)
		key := "user/" + pageData.UserInfo.ID.String() + pfpext
		_, err = Svc.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			Body:   pfpFile,
		})
		if err != nil {
			log.Printf("Failed to upload pfp to aws for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Could not update profile pic!",
			})
			tmpl.Execute(w, pageData)
			return
		}
		s3url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, aws_region, key)
		_, err = Db.Exec("UPDATE userdata SET profilepic = $1 WHERE uid = $2", s3url, pageData.UserInfo.ID)
	}
	tmpl.Execute(w, pageData)
}
