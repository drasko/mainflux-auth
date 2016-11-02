package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
)

// CheckPermissions checks whether or not a provided key has sufficient
// permissions to perform specified action on a specified resource.
func CheckPermissions(req *domain.AccessRequest) error {
	if valid := req.Validate(); !valid {
		return &domain.ServiceError{Code: http.StatusBadRequest}
	}

	content, err := domain.ContentOf(req.Key)
	if err != nil {
		return err
	}

	c := cache.Connection()
	defer c.Close()

	ac := &accessChecker{db: c, req: req, data: content}
	ac.check()
	return ac.err
}

type accessChecker struct {
	db   redis.Conn
	req  *domain.AccessRequest
	data domain.KeyData
	err  error
}

func (ac *accessChecker) check() {
	ac.verifyOwnership()
	if ac.err == nil {
		ac.checkMasterKey()
		ac.checkUserKeys()
		ac.checkDeviceKeys()
	}
}

func (ac *accessChecker) verifyOwnership() {
	userId := ac.data.Subject
	if (ac.req.Type == domain.UserType || ac.req.Type == domain.DevType) && userId != ac.req.Owner {
		ac.err = &domain.ServiceError{Code: http.StatusForbidden}
		return
	}
}

func (ac *accessChecker) checkMasterKey() {
	userId := ac.data.Subject
	cKey := fmt.Sprintf("users:%s", userId)

	mKey, _ := redis.String(ac.db.Do("HGET", cKey, "masterKey"))
	if mKey != ac.req.Key {
		ac.err = &domain.ServiceError{Code: http.StatusForbidden}
		return
	}

	if ac.req.Type == domain.UserType && ac.req.Id == userId {
		return
	}

	if ac.req.Type != domain.UserType {
		devId := ac.req.Id
		if ac.req.Type == domain.ChanType {
			devId = ac.req.Owner
		}

		cKey := fmt.Sprintf("users:%s:devices:%s:keys", userId, devId)
		if exists, _ := redis.Bool(ac.db.Do("EXISTS", cKey)); exists {
			return
		}
	}

	ac.err = &domain.ServiceError{Code: http.StatusForbidden}
}

func (ac *accessChecker) checkUserKeys() {
	if ac.err == nil {
		return
	}

	userId := ac.data.Subject
	cKey := fmt.Sprintf("users:%s:keys", userId)

	if exists, _ := redis.Bool(ac.db.Do("SISMEMBER", cKey, ac.req.Key)); exists {
		ac.checkScopes()
		return
	}

	ac.err = &domain.ServiceError{Code: http.StatusForbidden}
}

func (ac *accessChecker) checkDeviceKeys() {
	if ac.err == nil {
		return
	}

	userId := ac.data.Subject

	devId := ac.req.Id
	if ac.req.Type == domain.ChanType {
		devId = ac.req.Owner
	}

	cKey := fmt.Sprintf("users:%s:devices:%s:keys", userId, devId)
	if exists, _ := redis.Bool(ac.db.Do("SISMEMBER", cKey, ac.req.Key)); exists {
		ac.checkScopes()
		return
	}

	ac.err = &domain.ServiceError{Code: http.StatusForbidden}
}

func (ac *accessChecker) checkScopes() {
	for _, s := range ac.data.Scopes {
		if s.Type == ac.req.Type && s.Id == ac.req.Id && strings.Contains(s.Actions, ac.req.Action) {
			ac.err = nil
			return
		}
	}

	ac.err = &domain.ServiceError{Code: http.StatusForbidden}
}
