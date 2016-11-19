package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

// CheckPermissions determines whether or not a platform request can be
// performed given the provided key.
func CheckPermissions(key string, req domain.AccessRequest) error {
	if valid := req.Validate(); !valid {
		return &domain.AuthError{Code: http.StatusBadRequest}
	}

	_, err := domain.SubjectOf(key)
	if err != nil {
		return err
	}

	wList := fmt.Sprintf("auth:%s:%s:%s", req.Type, req.Id, req.Action)

	c := cache.Connection()
	defer c.Close()

	if allowed, _ := redis.Bool(c.Do("SISMEMBER", wList, key)); !allowed {
		return &domain.AuthError{Code: http.StatusForbidden}
	}

	return nil
}
