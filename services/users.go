package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth-server/cache"
	"github.com/mainflux/mainflux-auth-server/domain"
)

func RegisterUser(username, password string) (domain.User, error) {
	var user domain.User

	if username == "" || password == "" {
		return user, &domain.ServiceError{Code: http.StatusBadRequest}
	}

	c := cache.Connection()
	defer c.Close()

	cVal, err := redis.Int64(c.Do("SADD", "users", username))
	if err != nil {
		return user, &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	if cVal == 0 {
		return user, &domain.ServiceError{Code: http.StatusConflict}
	}

	user, err = domain.CreateUser(username, password)
	if err != nil {
		return user, err
	}

	//
	// NOTE: consider using MULTI to ensure consistency
	//
	cKey := "users:" + user.Username
	_, err = c.Do("HMSET", cKey, "id", user.Id, "password", user.Password, "masterKey", user.MasterKey)
	if err != nil {
		return user, &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return user, nil
}

func AddUserKey(uid, key string, payload domain.Payload) (string, error) {
	c := cache.Connection()
	defer c.Close()

	cKey := "users:" + uid
	masterKey, _ := redis.String(c.Do("HGET", cKey, "masterKey"))

	if masterKey == "" {
		fmt.Println("JEBITE SE FUDBALEEERIIIIII")
		return "", &domain.ServiceError{Code: http.StatusNotFound}
	}

	if key != masterKey {
		return "", &domain.ServiceError{Code: http.StatusForbidden}
	}

	newKey, err := domain.CreateKey(uid, &payload)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	cKey = cKey + ":keys"
	_, err = c.Do("SADD", cKey, newKey)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return newKey, nil
}
