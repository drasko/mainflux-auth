package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddUserKey(t *testing.T) {
	var (
		username string            = "user-key-username"
		password string            = "user-key-password"
		owner    string            = "key-owner"
		access   domain.AccessSpec = domain.AccessSpec{owner, []domain.Scope{{"R", domain.DevType, "dev"}}}
	)

	user, _ := services.RegisterUser(username, password)

	cases := []struct {
		master string
		access domain.AccessSpec
		code   int
	}{
		{user.MasterKey, access, http.StatusOK},
		{user.MasterKey, domain.AccessSpec{}, http.StatusOK},
		{"bad", access, http.StatusForbidden},
		{"bad", domain.AccessSpec{}, http.StatusForbidden},
	}

	for i, c := range cases {
		key, err := services.AddKey(c.master, c.access)
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
