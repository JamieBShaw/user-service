package model

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (u *User) IsAdmin() bool {
	if u == nil {
		return false
	}
	return u.Admin
}

func ToJson(rw http.ResponseWriter, status int, body interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(body)
}

func (u *User) Validate() error {
	ok, err := validateUser(u)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("error validating user")
	}
	return nil
}

func validateUser(user *User) (bool, error) {
	if user == nil {
		return false, errors.New("user is empty")
	}
	if user.ID <= 0 {
		return false, errors.New("user id is invalid")
	}
	if user.Username == "" {
		return false, errors.New("username is empty")
	}
	if len(user.Username) > 12 {
		return false, errors.New("username is to long")
	}
	if len(user.Username) < 3 {
		return false, errors.New("username is too short")
	}
	return true, nil
}
