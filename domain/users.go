package domain

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents one Mainflux user account.
type User struct {
	Id        string `json:"id"`
	Username  string `json:"-"`
	Password  string `json:"-"`
	MasterKey string `json:"key"`
}

// CreateUser creates new user account based on provided username and password.
// The account is assigned with one master key - a key with all permissions on
// all owned resources regardless of their type. Provided password in encrypted
// using bcrypt algorithm.
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
