package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func homePage(w http.ResponseWriter, r *http.Request) {
	pageData := PageData{
		PageTitle: "Homepage",
	}
	cookie, err := r.Cookie("session_cookie")
	if err == nil {
		err := parseCookie(cookie, &pageData.UserInfo)
		if err != nil {
			log.Printf("Failed to parse cookie for %s: %s", r.RemoteAddr, err)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    2,
				ErrorMessage: "failed to parse cookie",
			})
		}
	}
	tmpl, err := template.ParseFiles(
		"static/index.tmpl",
		"static/header.tmpl",
		"static/error.tmpl",
		"static/navbar.tmpl",
		"static/footer.tmpl",
	)
	if err != nil {
		log.Print("Failed to parse files: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, err := Db.Query("SELECT name, brand, pic FROM productdata")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var product ProductData
		err := rows.Scan(
			&product.ProdName,
			&product.ProdBrand,
			&product.ProdImage,
		)
		if err != nil {
			log.Fatal(err)
		}
		pageData.Products = append(pageData.Products, product)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Print("Failed to render page: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

		age, err := strconv.Atoi(r.FormValue("age"))
		if err != nil {
			log.Printf("Could not update profile for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Invalid age!",
			})
			tmpl.Execute(w, pageData)
			return
		}
		calories, err := strconv.Atoi(r.FormValue("targetCalories"))
		if err != nil {
			log.Printf("Could not update profile for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Invalid calories value!",
			})
			tmpl.Execute(w, pageData)
			return
		}

		pageData.UserDetails = UserData{
			FirstName:      r.Form.Get("fname"),
			LastName:       r.Form.Get("lname"),
			Age:            age,
			ProfilePic:     r.FormValue("profilepic"),
			TargetCalories: calories,
			Email:          r.FormValue("email"),
			Location:       r.FormValue("location"),
			Allergens:      strings.Split(r.FormValue("allergens"), ","),
		}

		allergens := r.Form.Get("allergens")

		if pageData.UserDetails.ProfilePic != "" {
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
					ErrorMessage: "Could not update profile pic! (S3)",
				})
			}
			s3url := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, aws_region, key)
			_, err = Db.Exec("UPDATE userdata SET profilepic = $1 WHERE uid = $2", s3url, pageData.UserInfo.ID)
			if err != nil {
				log.Printf("Failed to upload pfp to db for %s", r.RemoteAddr)
				pageData.PageError = append(pageData.PageError, Error{
					ErrorCode:    1,
					ErrorMessage: "Could not update profile pic! (DB)",
				})
			}
		}
		if allergens != "" {
			allergen_slice := strings.Split(allergens, ",")
			result, err := Db.Exec("UPDATE userdata SET allergens = $1 WHERE uid = $2", pq.Array(allergen_slice), pageData.UserInfo.ID)
			log.Print(result)
			if err != nil {
				log.Printf("Failed to update allergens for %s", r.RemoteAddr)
				pageData.PageError = append(pageData.PageError, Error{
					ErrorCode:    1,
					ErrorMessage: "Could not update allergens!",
				})
			}
		}
		query := "UPDATE userdata SET"
		updates := ""
		if pageData.UserDetails.FirstName != "" {
			updates += fmt.Sprintf(" fname = '%s',", pageData.UserDetails.FirstName)
		}
		if pageData.UserDetails.LastName != "" {
			updates += fmt.Sprintf(" lname = '%s',", pageData.UserDetails.LastName)
		}
		if pageData.UserDetails.Email != "" {
			updates += fmt.Sprintf(" email = '%s',", pageData.UserDetails.Email)
		}
		if pageData.UserDetails.Location != "" {
			updates += fmt.Sprintf(" fname = '%s',", pageData.UserDetails.Location)
		}
		updates += fmt.Sprintf(" target_calories = '%d',", pageData.UserDetails.TargetCalories)
		updates = updates[:len(updates)-1]

		query += updates + fmt.Sprintf(" WHERE uid = %s", pageData.UserInfo.ID)
		_, err = Db.Exec(query)
		if err != nil {
			fmt.Print("Failed to modify userdata table: ", err)
			http.Error(w, "Error updating data in database", http.StatusInternalServerError)
			return
		}

	}
	tmpl.Execute(w, pageData)
}
