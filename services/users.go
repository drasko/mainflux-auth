package services

import (
	"net/http"

	"github.com/mainflux/mainflux-auth-server/cache"
	"github.com/mainflux/mainflux-auth-server/domain"
)

func RegisterUser(username, password string) (domain.User, error) {
	var user domain.User
	if username == "" || password == "" {
		return user, &domain.AuthError{Code: http.StatusBadRequest}
	}

	c := cache.Connection()
	defer c.Close()

	n, err := c.Do("SADD", "users", username)
	if err != nil {
		return user, &domain.AuthError{Code: http.StatusInternalServerError}
	}

	if n.(int64) == 0 {
		return user, &domain.AuthError{Code: http.StatusConflict}
	}

	user, err = domain.CreateUser(username, password)
	if err != nil {
		return user, err
	}

	//
	// NOTE: consider using MULTI to ensure consistency
	//
	key := "users:" + user.Username
	_, err = c.Do("HMSET", key, "id", user.Id, "password", user.Password, "key", user.MasterKey)
	if err != nil {
		return user, &domain.AuthError{Code: http.StatusInternalServerError}
	}

	return user, nil
}
