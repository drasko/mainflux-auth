package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth-server/cache"
	"github.com/mainflux/mainflux-auth-server/domain"
)

func AddDeviceKey(uid, dev, key string, payload domain.Payload) (string, error) {
	c := cache.Connection()
	defer c.Close()

	cKey := "users:" + uid
	masterKey, _ := redis.String(c.Do("HGET", cKey, "masterKey"))

	if masterKey == "" {
		return "", &domain.ServiceError{Code: http.StatusNotFound}
	}

	if key != masterKey {
		return "", &domain.ServiceError{Code: http.StatusForbidden}
	}

	newKey, err := domain.CreateKey(uid, &payload)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	cKey = fmt.Sprintf("users:%s:devices:%s:keys", uid, dev)
	_, err = c.Do("SADD", cKey, newKey)
	if err != nil {
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return newKey, nil
}
