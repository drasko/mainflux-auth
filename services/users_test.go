package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
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

	for _, c := range cases {
		_, err := services.RegisterUser(c.username, c.password)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("expected %d got %d", c.code, auth.Code)
			}
		}
	}
}

func TestAddUserKey(t *testing.T) {
	cases := []struct {
		uid     string
		key     string
		payload domain.Payload
		code    int
	}{
		{user.Id, user.MasterKey, domain.Payload{}, 0},
		{user.Id, "invalid-key", domain.Payload{}, http.StatusForbidden},
		{"invalid-id", user.MasterKey, domain.Payload{}, http.StatusNotFound},
		{"invalid-id", "invalid-key", domain.Payload{}, http.StatusNotFound},
	}

	for _, c := range cases {
		key, err := services.AddUserKey(c.uid, c.key, c.payload)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("expected %d got %d", c.code, auth.Code)
			}
			continue
		}

		if key == "" {
			t.Errorf("expected key to be created")
		}
	}
}
