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
		// master key request access to user
		{domain.AccessRequest{"R", "user", user.Id, "", user.MasterKey}, true},
		{domain.AccessRequest{"W", "user", user.Id, "", user.MasterKey}, true},
		{domain.AccessRequest{"X", "user", user.Id, "", user.MasterKey}, true},

		// master key request access to device
		{domain.AccessRequest{"R", "device", devId, "", user.MasterKey}, true},
		{domain.AccessRequest{"W", "device", devId, "", user.MasterKey}, true},
		{domain.AccessRequest{"X", "device", devId, "", user.MasterKey}, true},

		// master key request access to channel
		{domain.AccessRequest{"R", "channel", chanId, devId, user.MasterKey}, true},
		{domain.AccessRequest{"W", "channel", chanId, devId, user.MasterKey}, true},
		{domain.AccessRequest{"X", "channel", chanId, devId, user.MasterKey}, true},

		// illegal master key access (resources not owned)
		{domain.AccessRequest{"R", "user", "bad", "", user.MasterKey}, false},
		{domain.AccessRequest{"R", "device", "bad", "", user.MasterKey}, false},
		{domain.AccessRequest{"R", "channel", "bad", "bad", user.MasterKey}, false},

		// user key request access to user
		{domain.AccessRequest{"R", "user", user.Id, "", usrKey}, false},

		// user key request access to device
		{domain.AccessRequest{"R", "device", devId, "", usrKey}, true},
		{domain.AccessRequest{"W", "device", devId, "", usrKey}, true},
		{domain.AccessRequest{"X", "device", devId, "", usrKey}, false},

		// user key request access to channel
		{domain.AccessRequest{"R", "channel", chanId, devId, usrKey}, true},
		{domain.AccessRequest{"W", "channel", chanId, devId, usrKey}, false},

		// device key request access to user
		{domain.AccessRequest{"R", "user", user.Id, "", devKey}, false},

		// device key request access to device
		{domain.AccessRequest{"R", "device", devId, "", devKey}, true},
		{domain.AccessRequest{"W", "device", devId, "", devKey}, true},
		{domain.AccessRequest{"X", "device", devId, "", devKey}, false},

		// device key request access to channel
		{domain.AccessRequest{"R", "channel", chanId, devId, devKey}, true},
		{domain.AccessRequest{"W", "channel", chanId, devId, devKey}, false},

		// missing or invalid request data
		{domain.AccessRequest{"illegal-action", "device", devId, "", devKey}, false},
		{domain.AccessRequest{"R", "illegal-resource", devId, "", devKey}, false},
		{domain.AccessRequest{"R", "device", "unknown-id", "", devKey}, false},
		{domain.AccessRequest{"R", "device", devId, "", "invalid-key"}, false},
		{domain.AccessRequest{"R", "channel", devId, "", devKey}, false},
	}

	for i, c := range cases {
		err := services.CheckPermissions(&c.AccessRequest)

		granted := err == nil

		if c.granted == !granted {
			t.Errorf("case %d: expected granted %t got %t", i+1, c.granted, granted)
		}
	}
}
