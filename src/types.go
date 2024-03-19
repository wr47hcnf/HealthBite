package main

import "github.com/google/uuid"

type Error struct {
	ErrorCode    int
	ErrorMessage string
}

type User struct {
	IsLogged   bool
	ID         uuid.UUID
	Username   string
	Profilepic string
}

type PageData struct {
	PageTitle   string
	UserInfo    User
	UserDetails UserData
	Products    []ProductData
	PageError   []Error
}

type UserData struct {
	FirstName      string
	LastName       string
	Age            int
	ProfilePic     string
	TargetCalories int
	Email          string
	Location       string
	Allergens      []string
}

type ProductData struct {
	ProdID          uuid.UUID
	ProdName        string
	ProdBarcode     string
	ProdBrand       string
	ProdLocation    string
	ProdCalories    string
	ProdImage       string
	ProdWeight      int
	NutritionalInfo string
	ProdAllergens   []string
	ProdAdditives   []string
}
