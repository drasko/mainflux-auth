package domain

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string `json:"id"`
	Username  string `json:"-"`
	Password  string `json:"-"`
	MasterKey string `json:"key"`
}

func CreateUser(username, password string) (User, error) {
	user := User{
		Id:       uuid.NewV4().String(),
		Username: username,
	}

	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, &ServiceError{Code: http.StatusInternalServerError}
	}

	user.Password = string(p)

	// master payload: can do everything on all resources
	masterScope := Scope{"RWX", "*", "*"}
	payload := Payload{Scopes: []Scope{masterScope}}

	user.MasterKey, err = CreateKey(user.Id, &payload)
	if err != nil {
		return user, &ServiceError{Code: http.StatusInternalServerError}
	}

	return user, nil
}
