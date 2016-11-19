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
