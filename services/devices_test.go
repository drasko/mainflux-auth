package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddDeviceKey(t *testing.T) {
	var (
		username string            = "dev-key-username"
		password string            = "dev-key-password"
		access   domain.AccessSpec = domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, "dev"}}}
	)

	user, _ := services.RegisterUser(username, password)

	cases := []struct {
		id     string
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
		key, err := services.AddDeviceKey(c.id, c.id, c.key, c.access)
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
