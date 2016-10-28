package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth-server/cache"
	"github.com/mainflux/mainflux-auth-server/domain"
)

// CheckPermissions checks whether or not a provided key has sufficient
// permissions to perform specified action on a specified resource.
func CheckPermissions(req domain.AccessRequest) error {
	if valid := req.Validate(); !valid {
		return &domain.ServiceError{Code: http.StatusBadRequest}
	}

	content, err := domain.ContentOf(req.Key)
	if err != nil {
		return err
	}

	c := cache.Connection()
	defer c.Close()

	err = checkUserKeys(c, req, content)
	if err != nil {
		return checkDeviceKeys(c, req, content)
	}

	return err
}

func checkUserKeys(c redis.Conn, req domain.AccessRequest, content domain.KeyData) error {
	userId := content.Subject

	cKey := fmt.Sprintf("users:%s", userId)
	mKey, err := redis.String(c.Do("HGET", cKey, "masterKey"))
	if err != nil {
		return &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	if mKey == req.Key {
		return nil
	}

	cKey = fmt.Sprintf("users:%s:keys", userId)
	exists, err := redis.Bool(c.Do("SISMEMBER", cKey, req.Key))
	if err != nil {
		return &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	if !exists {
		return &domain.ServiceError{Code: http.StatusForbidden}
	}

	return checkScopes(req, content)
}

func checkDeviceKeys(c redis.Conn, req domain.AccessRequest, content domain.KeyData) error {
	userId := content.Subject

	devId := req.Id
	if req.Resource == "channel" {
		devId = req.Device
	}

	cKey := fmt.Sprintf("users:%s:devices:%s:keys", userId, devId)
	exists, err := redis.Bool(c.Do("SISMEMBER", cKey, req.Key))
	if err != nil {
		return &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	if !exists {
		return &domain.ServiceError{Code: http.StatusForbidden}
	}

	return checkScopes(req, content)
}

func checkScopes(req domain.AccessRequest, content domain.KeyData) error {
	for _, s := range content.Scopes {
		if s.Resource == req.Resource && s.Id == req.Id && strings.Contains(s.Actions, req.Action) {
			return nil
		}
	}

	return &domain.ServiceError{Code: http.StatusForbidden}
}
