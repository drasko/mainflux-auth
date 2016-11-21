package services_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestCheckPermissions(t *testing.T) {
	var (
		username = "access-user"
		password = "access-pass"
		device   = "device-id"
		spec     = domain.KeySpec{"owner", []domain.Scope{{"CR", domain.DevType, device}}}
	)

	user, _ := services.RegisterUser(username, password)
	key, _ := services.AddKey(user.MasterKey, spec)

	cases := []struct {
		key     string
		request domain.AccessRequest
		granted bool
	}{
		{key, domain.AccessRequest{"C", domain.DevType, device}, true},
		{key, domain.AccessRequest{"R", domain.DevType, device}, true},
		{key, domain.AccessRequest{"U", domain.DevType, device}, false},
		{"bad", domain.AccessRequest{"R", domain.DevType, device}, false},
		{key, domain.AccessRequest{"bad", domain.DevType, device}, false},
		{key, domain.AccessRequest{"C", "bad", device}, false},
		{key, domain.AccessRequest{"C", domain.DevType, "bad"}, false},
	}

	for i, c := range cases {
		err := services.CheckPermissions(c.key, c.request)

		if granted := err == nil; c.granted != granted {
			t.Errorf("case %d: expected granted %t got %t", i+1, c.granted, granted)
		}
	}
}
