package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

// AddKey creates new API key given a master key, and new key specification.
func AddKey(mKey string, spec domain.KeySpec) (string, error) {
	c := cache.Connection()
	defer c.Close()

	id, err := domain.SubjectOf(mKey)
	if err != nil {
		return "", err
	}

	cKey := fmt.Sprintf("auth:user:%s:master", id)
	if k, _ := redis.String(c.Do("GET", cKey)); k != mKey {
		return "", &domain.AuthError{Code: http.StatusForbidden}
	}

	key, _ := domain.CreateKey(spec.Owner)

	c.Send("MULTI")
	c.Send("SADD", fmt.Sprintf("auth:%s:%s:keys", domain.UserType, id), key)

	for _, scope := range spec.Scopes {
		for _, action := range scope.Actions {
			wList := fmt.Sprintf("auth:%s:%s:%s", scope.Type, scope.Id, string(action))
			c.Send("SADD", wList, key)
			c.Send("SADD", fmt.Sprintf("auth:key:%s", key), wList)
		}
	}

	_, err = c.Do("EXEC")
	if err != nil {
		return "", &domain.AuthError{Code: http.StatusInternalServerError}
	}

	return key, nil
}

// FetchKeys retrieves all keys created by user having a provided master key.
func FetchKeys(mKey string) (domain.KeyList, error) {
	var keys domain.KeyList

	c := cache.Connection()
	defer c.Close()

	user, err := domain.SubjectOf(mKey)
	if err != nil {
		return keys, &domain.AuthError{Code: http.StatusForbidden}
	}

	cKey := fmt.Sprintf("auth:user:%s:master", user)
	if cVal, _ := redis.String(c.Do("GET", cKey)); cVal != mKey {
		return keys, &domain.AuthError{Code: http.StatusForbidden}
	}

	cKey = fmt.Sprintf("auth:user:%s:keys", user)
	keys.Keys, err = redis.Strings(c.Do("SMEMBERS", cKey))
	if err != nil {
		return keys, &domain.AuthError{Code: http.StatusInternalServerError}
	}

	return keys, nil
}
