package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestRegisterUser(t *testing.T) {
	username := "test"
	password := "test"

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

func TestAddUserKey(t *testing.T) {
	access := domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, "test-id"}}}

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
