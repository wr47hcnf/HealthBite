package main

import (
	"github.com/google/uuid"
	"net/http"
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
