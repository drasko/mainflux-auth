package services

import (
	"fmt"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

func AddKey(mKey string, spec domain.AccessSpec) (string, error) {
	c := cache.Connection()
	defer c.Close()

	id, err := domain.SubjectOf(mKey)
	if err != nil {
		return "", err
	}

	cKey := fmt.Sprintf("auth:user:%s:master", id)
	if k, _ := redis.String(c.Do("GET", cKey)); k != mKey {
		return "", &domain.ServiceError{Code: http.StatusForbidden}
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
		return "", &domain.ServiceError{Code: http.StatusInternalServerError}
	}

	return key, nil
}
