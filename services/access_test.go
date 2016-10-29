package services_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
)

func TestHasSufficientPermissions(t *testing.T) {
	var (
		devId  string            = "test-device"
		chanId string            = "test-channel"
		access domain.AccessSpec = domain.AccessSpec{[]domain.Scope{{"RW", "device", devId}, {"R", "channel", chanId}}}
	)

	devKey, _ := services.AddDeviceKey(user.Id, devId, user.MasterKey, access)
	usrKey, _ := services.AddUserKey(user.Id, user.MasterKey, access)

	cases := []struct {
		domain.AccessRequest
		granted bool
	}{
		{domain.AccessRequest{"R", "user", user.Id, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"W", "user", user.Id, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"X", "user", user.Id, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"R", "device", devId, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"W", "device", devId, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"X", "device", devId, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"R", "channel", chanId, devId, user.MasterKey}, true},
		{domain.AccessRequest{"W", "channel", chanId, devId, user.MasterKey}, true},
		{domain.AccessRequest{"X", "channel", chanId, devId, user.MasterKey}, true},

		// illegal master key access (resources not owned)
		{domain.AccessRequest{"R", "user", "bad", user.Id, user.MasterKey}, false},
		{domain.AccessRequest{"R", "device", "bad", user.Id, user.MasterKey}, false},
		{domain.AccessRequest{"R", "channel", "bad", "bad", user.MasterKey}, false},

		{domain.AccessRequest{"R", "user", user.Id, user.Id, usrKey}, false},
		{domain.AccessRequest{"R", "device", devId, user.Id, usrKey}, true},
		{domain.AccessRequest{"W", "device", devId, user.Id, usrKey}, true},
		{domain.AccessRequest{"X", "device", devId, user.Id, usrKey}, false},
		{domain.AccessRequest{"R", "channel", chanId, devId, usrKey}, true},
		{domain.AccessRequest{"W", "channel", chanId, devId, usrKey}, false},

		// device key requests
		{domain.AccessRequest{"R", "user", user.Id, user.Id, devKey}, false},
		{domain.AccessRequest{"R", "device", devId, user.Id, devKey}, true},
		{domain.AccessRequest{"W", "device", devId, user.Id, devKey}, true},
		{domain.AccessRequest{"X", "device", devId, user.Id, devKey}, false},
		{domain.AccessRequest{"R", "channel", chanId, devId, devKey}, true},
		{domain.AccessRequest{"W", "channel", chanId, devId, devKey}, false},

		// missing or invalid request data
		{domain.AccessRequest{"illegal-action", "device", devId, user.Id, devKey}, false},
		{domain.AccessRequest{"R", "illegal-resource", devId, user.Id, devKey}, false},
		{domain.AccessRequest{"R", "device", "unknown-id", user.Id, devKey}, false},
		{domain.AccessRequest{"R", "device", devId, user.Id, "invalid-key"}, false},
		{domain.AccessRequest{"R", "channel", devId, devId, devKey}, false},
	}

	for i, c := range cases {
		err := services.CheckPermissions(&c.AccessRequest)
		granted := err == nil

		if c.granted == !granted {
			t.Errorf("case %d: expected granted %t got %t", i+1, c.granted, granted)
		}
	}
}
