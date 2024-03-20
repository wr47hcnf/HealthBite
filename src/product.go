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

func addProduct(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"static/product_submission.tmpl",
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
		productUUID, err := uuid.NewRandom()
		if err != nil {
			log.Printf("Failed to generate product uuid for %s", r.RemoteAddr)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Could not generate product ID!",
			})
			tmpl.Execute(w, pageData)
			return
		}
		productpic := r.Form.Get("productPhoto")
		var s3url string
		if productpic != "" {
			pfpFile, pfpHeader, err := r.FormFile("productPhoto")
			if err != nil {
				log.Printf("Failed to parse pfp for %s", r.RemoteAddr)
				pageData.PageError = append(pageData.PageError, Error{
					ErrorCode:    1,
					ErrorMessage: "Could not parse product pic!",
				})
				tmpl.Execute(w, pageData)
				return
			}
			pfpext := filepath.Ext(pfpHeader.Filename)
			key := "product/" + pageData.Products[0].ProdID.String() + pfpext
			_, err = Svc.PutObject(&s3.PutObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(key),
				Body:   pfpFile,
			})
			if err != nil {
				log.Printf("Failed to upload pfp to aws for %s", r.RemoteAddr)
				pageData.PageError = append(pageData.PageError, Error{
					ErrorCode:    1,
					ErrorMessage: "Could not update product pic! (S3)",
				})
			}
			s3url = fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, aws_region, key)
		}
		weight, err := strconv.Atoi(r.FormValue("productWeight"))
		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Could not update product weight! %s", err),
			})
		}
		fat, err := strconv.Atoi(r.FormValue("productFat"))
		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Could not update product fat! %s", err),
			})
		}
		sodium, err := strconv.Atoi(r.FormValue("productSodium"))
		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Could not update product sodium! %s", err),
			})
		}
		carbs, err := strconv.Atoi(r.FormValue("productCarbohydrates"))
		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Could not update product carbohydrates! %s", err),
			})
		}
		protein, err := strconv.Atoi(r.FormValue("productProtein"))
		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Could not update product protein! %s", err),
			})
		}
		price, err := strconv.Atoi(r.FormValue("productPrice"))
		if err != nil {
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: fmt.Sprintf("Could not update product price! %s", err),
			})
		}
		pageData.Products = append(pageData.Products, ProductData{
			ProdID:        productUUID,
			ProdName:      r.FormValue("productName"),
			ProdBarcode:   r.FormValue("productBarcode"),
			ProdLocation:  r.FormValue("productLocation"),
			ProdBrand:     r.FormValue("productBrand"),
			ProdCalories:  r.FormValue("productCalories"),
			ProdWeight:    weight,
			ProdFat:       fat,
			ProdSodium:    sodium,
			ProdCarbs:     carbs,
			ProdProtein:   protein,
			ProdPrice:     price,
			ProdImage:     s3url,
			ProdAdditives: strings.Split(r.FormValue("productAdditives"), ","),
			ProdAllergens: strings.Split(r.FormValue("productAllergens"), ","),
		})
		result, err := Db.Exec("INSERT INTO productdata (prodid, barcode, name, brand, price, pic, location, weight, calories, fat, sodium, carbohydrates, protein, additives, allergens) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
			pageData.Products[0].ProdID,
			pageData.Products[0].ProdBarcode,
			pageData.Products[0].ProdName,
			pageData.Products[0].ProdBrand,
			pageData.Products[0].ProdPrice,
			pageData.Products[0].ProdImage,
			pageData.Products[0].ProdLocation,
			pageData.Products[0].ProdWeight,
			pageData.Products[0].ProdCalories,
			pageData.Products[0].ProdFat,
			pageData.Products[0].ProdSodium,
			pageData.Products[0].ProdCarbs,
			pageData.Products[0].ProdProtein,
			pq.Array(pageData.Products[0].ProdAdditives),
			pq.Array(pageData.Products[0].ProdAllergens),
		)
		log.Print(result)
		if err != nil {
			log.Print(err)
			pageData.PageError = append(pageData.PageError, Error{
				ErrorCode:    1,
				ErrorMessage: "Could not update product!",
			})
		}
	}
	pageData.PageError = append(pageData.PageError, Error{
		ErrorCode:    5,
		ErrorMessage: "Product" + pageData.Products[0].ProdName + "added!",
	})
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Print("Failed to render page: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func viewProduct(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"static/product_page.tmpl",
		"static/error.tmpl",
		"static/header.tmpl",
		"static/navbar.tmpl",
		"static/footer.tmpl",
	))
	productParam := r.URL.Query().Get("pid")
	rows, err := Db.Query(`SELECT 
	barcode, name, brand, pic, location, weight, calories, fat, sodium, carbohydrates, protein, additives, allergens
	FROM productdata WHERE pid=$1`, productParam)
	if err != nil {
		log.Fatal(err)
	}
	pageData := PageData{
		PageTitle: "Product Details",
	}
	defer rows.Close()
	for rows.Next() {
		var product ProductData
		err := rows.Scan(
			&product.ProdBarcode,
			&product.ProdName,
			&product.ProdBrand,
			&product.ProdImage,
			&product.ProdLocation,
			&product.ProdWeight,
			&product.ProdCalories,
			&product.ProdFat,
			&product.ProdSodium,
			&product.ProdCarbs,
			&product.ProdProtein,
			pq.Array(&product.ProdAdditives),
			pq.Array(&product.ProdAllergens),
		)
		if err != nil {
			log.Fatal(err)
		}
		pageData.Products = append(pageData.Products, product)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	pageData = PageData{
		PageTitle: pageData.Products[0].ProdName,
	}
	err = tmpl.Execute(w, pageData)
	if err != nil {
		log.Print("Failed to render page: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func searchProduct(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"static/products_list.tmpl",
		"static/error.tmpl",
		"static/header.tmpl",
		"static/navbar.tmpl",
		"static/footer.tmpl",
	))
	productParam := r.URL.Query().Get("product")
	pageData := PageData{
		PageTitle: productParam,
	}
	rows, err := Db.Query("SELECT prodid, name, brand, pic, location, calories FROM productdata")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var product ProductData
		err := rows.Scan(
			&product.ProdID,
			&product.ProdName,
			&product.ProdBrand,
			&product.ProdImage,
			&product.ProdLocation,
			&product.ProdCalories,
		)
		if err != nil {
			log.Fatal(err)
		}
		product.ProdName = strings.ToLower(product.ProdName)
		productParam = strings.ToLower(productParam)
		if strings.Contains(product.ProdName, productParam) {
			pageData.Products = append(pageData.Products, product)
		}
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
