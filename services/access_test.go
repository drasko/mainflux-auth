package services_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestCheckPermissions(t *testing.T) {
	var (
		username string         = "access-user"
		password string         = "access-pass"
		device   string         = "device-id"
		spec     domain.KeySpec = domain.KeySpec{"owner", []domain.Scope{{"RW", "device", device}}}
	)

	user, _ := services.RegisterUser(username, password)
	key, _ := services.AddKey(user.MasterKey, spec)

	cases := []struct {
		key     string
		request domain.AccessRequest
		granted bool
	}{
		{key, domain.AccessRequest{"R", "device", device}, true},
		{key, domain.AccessRequest{"W", "device", device}, true},
		{key, domain.AccessRequest{"X", "device", device}, false},
		{"bad", domain.AccessRequest{"R", "device", device}, false},
		{key, domain.AccessRequest{"X", "device", "bad"}, false},
	}

	for i, c := range cases {
		err := services.CheckPermissions(c.key, c.request)

		if granted := err == nil; c.granted != granted {
			t.Errorf("case %d: expected granted %t got %t", i+1, c.granted, granted)
		}
	}
}
