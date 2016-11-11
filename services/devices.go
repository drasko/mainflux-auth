package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

// AddDeviceKey adds new device key based on provided access spec. Bear in mind
// that the device key can be created only when identified as "master", i.e.
// by providing a master key.
func AddDeviceKey(userId, devId, key string, access domain.AccessSpec) (string, error) {
	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("users:%s:master", userId)
	masterKey, _ := redis.String(c.Do("GET", cKey))

	if masterKey == "" {
		return "", &domain.ServiceError{Code: http.StatusNotFound}
	}

	if key != masterKey {
		return "", &domain.ServiceError{Code: http.StatusForbidden}
	}

	if valid := access.Validate(); !valid {
		return "", &domain.ServiceError{Code: http.StatusBadRequest}
	}

	newKey, err := domain.CreateKey(userId, &access)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	cKey = fmt.Sprintf("users:%s:devices:%s:keys", userId, devId)
	_, err = c.Do("SADD", cKey, newKey)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return newKey, nil
}

// FetchDeviceKeys retrieves device keys. It can be accessed only by providing
// a valid master key.
func FetchDeviceKeys(userId, devId, key string) (domain.KeyList, error) {
	var list domain.KeyList

	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("users:%s:master", userId)
	mKey, _ := redis.String(c.Do("GET", cKey))

	if mKey == "" {
		return list, &domain.ServiceError{Code: http.StatusNotFound}
	}

	if key != mKey {
		return list, &domain.ServiceError{Code: http.StatusForbidden}
	}

	cKey = fmt.Sprintf("users:%s:devices:%s:keys", userId, devId)
	keys, err := redis.Strings(c.Do("SMEMBERS", cKey))
	if err != nil {
		return list, &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	list.Keys = keys
	return list, nil
}
