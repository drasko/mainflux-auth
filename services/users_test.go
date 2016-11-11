package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestRegisterUser(t *testing.T) {
	var (
		username string = "register-username"
		password string = "register-password"
	)

	cases := []struct {
		username string
		password string
		code     int
	}{
		{username, password, 0},
		{username, password, http.StatusConflict},
		{username, "", http.StatusBadRequest},
		{"", password, http.StatusBadRequest},
	}

	for i, c := range cases {
		_, err := services.RegisterUser(c.username, c.password)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
		}
	}
}

func TestLogin(t *testing.T) {
	var (
		username string = "login-username"
		password string = "login-password"
	)

	services.RegisterUser(username, password)

	cases := []struct {
		username string
		password string
		code     int
	}{
		{username, password, 0},
		{username, "", http.StatusBadRequest},
		{"", password, http.StatusBadRequest},
		{"bad", password, http.StatusForbidden},
		{username, "bad", http.StatusForbidden},
	}

	for i, c := range cases {
		_, err := services.Login(c.username, c.password)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
		}
	}
}

func TestAddUserKey(t *testing.T) {
	var (
		username string            = "user-key-username"
		password string            = "user-key-password"
		access   domain.AccessSpec = domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, "dev"}}}
	)

	user, _ := services.RegisterUser(username, password)

	cases := []struct {
		uid    string
		key    string
		access domain.AccessSpec
		code   int
	}{
		{user.Id, user.MasterKey, access, 0},
		{"bad", user.MasterKey, access, http.StatusNotFound},
		{user.Id, "bad", access, http.StatusForbidden},
		{user.Id, user.MasterKey, domain.AccessSpec{}, http.StatusBadRequest},
		{"bad", "bad", domain.AccessSpec{}, http.StatusNotFound},
	}

	for i, c := range cases {
		key, err := services.AddUserKey(c.uid, c.key, c.access)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
			continue
		}

		if key == "" {
			t.Errorf("case %d: expected key to be created", i+1)
		}
	}
}

func TestFetchUserKeys(t *testing.T) {
	oneKeyUser, _ := services.RegisterUser("one", "one")
	access := domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, "dev"}}}
	services.AddUserKey(oneKeyUser.Id, oneKeyUser.MasterKey, access)
	noKeysUser, _ := services.RegisterUser("empty", "empty")

	cases := []struct {
		uid   string
		key   string
		code  int
		total int
	}{
		{oneKeyUser.Id, oneKeyUser.MasterKey, 0, 1},
		{noKeysUser.Id, noKeysUser.MasterKey, 0, 0},
		{oneKeyUser.Id, "bad", http.StatusForbidden, 0},
		{"bad", oneKeyUser.MasterKey, http.StatusNotFound, 0},
	}

	for i, c := range cases {
		list, err := services.FetchUserKeys(c.uid, c.key)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
			continue
		}

		if len(list.Keys) != c.total {
			t.Errorf("case %d: expected %d items got %d", i+1, c.total, len(list.Keys))
		}
	}
}
