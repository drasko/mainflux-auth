package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
)

func TestAddDeviceKey(t *testing.T) {
	services.RegisterUser("test-dev", "test-dev")
	uid, masterKey, err := fetchCredentials()
	if err != nil {
		t.Errorf("failed to retrieve created user data")
	}

	cases := []struct {
		id      string
		key     string
		payload domain.Payload
		code    int
	}{
		{uid, masterKey, domain.Payload{}, 0},
		{uid, "invalid-key", domain.Payload{}, http.StatusForbidden},
		{"invalid-id", masterKey, domain.Payload{}, http.StatusNotFound},
		{"invalid-id", "invalid-key", domain.Payload{}, http.StatusNotFound},
	}

	for _, c := range cases {
		key, err := services.AddDeviceKey(c.id, c.id, c.key, c.payload)
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
