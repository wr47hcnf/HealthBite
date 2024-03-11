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
	PageTitle string
	UserInfo  User
	PageError []Error
}

type UserData struct {
	FirstName string
	LastName  string
	Age       int
}
