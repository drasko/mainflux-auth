package services

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

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

	cKey := fmt.Sprintf("auth:%s:%s:master", domain.UserType, id)
	if cVal, _ := redis.String(c.Do("GET", cKey)); cVal != mKey {
		return "", &domain.AuthError{Code: http.StatusForbidden}
	}

	key, _ := domain.CreateKey(spec.Owner)

	c.Send("MULTI")
	c.Send("SADD", fmt.Sprintf("auth:%s:%s:keys", domain.UserType, id), key)

	keyList := fmt.Sprintf("auth:keys:%s", key)

	for _, scope := range spec.Scopes {
		for _, action := range scope.Actions {
			wList := fmt.Sprintf("auth:%s:%s:%s", scope.Type, scope.Id, string(action))
			c.Send("SADD", wList, key)
			c.Send("SADD", keyList, wList)
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

	user, err := domain.SubjectOf(mKey)
	if err != nil {
		return keys, &domain.AuthError{Code: http.StatusForbidden}
	}

	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("auth:%s:%s:master", domain.UserType, user)
	if cVal, _ := redis.String(c.Do("GET", cKey)); cVal != mKey {
		return keys, &domain.AuthError{Code: http.StatusForbidden}
	}

	cKey = fmt.Sprintf("auth:%s:%s:keys", domain.UserType, user)
	keys.Keys, err = redis.Strings(c.Do("SMEMBERS", cKey))
	if err != nil {
		return keys, &domain.AuthError{Code: http.StatusInternalServerError}
	}

	return keys, nil
}

// FetchKeySpec retrieves key specification for given key. Note that the key
// specification can be retrieved only by providing a proper master key.
func FetchKeySpec(mKey string, key string) (domain.KeySpec, error) {
	var spec domain.KeySpec

	user, err := domain.SubjectOf(mKey)
	if err != nil {
		return spec, &domain.AuthError{Code: http.StatusForbidden}
	}

	c := cache.Connection()
	defer c.Close()

	cKey := fmt.Sprintf("auth:%s:%s:master", domain.UserType, user)
	if cVal, _ := redis.String(c.Do("GET", cKey)); cVal != mKey {
		return spec, &domain.AuthError{Code: http.StatusForbidden}
	}

	cKey = fmt.Sprintf("auth:%s:%s:keys", domain.UserType, user)
	if exists, _ := redis.Bool(c.Do("SISMEMBER", cKey, key)); !exists {
		return spec, &domain.AuthError{Code: http.StatusNotFound}
	}

	owner, err := domain.SubjectOf(key)
	if err != nil {
		return spec, &domain.AuthError{Code: http.StatusForbidden}
	}

	cKey = fmt.Sprintf("auth:keys:%s", key)
	wLists, err := redis.Strings(c.Do("SMEMBERS", cKey))
	if err != nil {
		return spec, &domain.AuthError{Code: http.StatusInternalServerError}
	}

	collect(wLists, &spec)
	spec.Owner = owner

	return spec, nil
}

func collect(wLists []string, spec *domain.KeySpec) {
	scopes := map[string]string{}

	for _, scope := range wLists {
		parts := strings.Split(scope, ":")
		resource := strings.Join(parts[1:3], ":")
		scopes[resource] = scopes[resource] + parts[3]
	}

	for k, v := range scopes {
		parts := strings.Split(k, ":")

		actions := strings.Split(v, "")
		sort.Strings(actions)

		scope := domain.Scope{
			Type:    parts[0],
			Id:      parts[1],
			Actions: strings.Join(actions, ""),
		}

		spec.Scopes = append(spec.Scopes, scope)
	}
}
