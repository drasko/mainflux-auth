package services_test

import (
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestHasSufficientPermissions(t *testing.T) {
	var (
		devId  string            = "test-device"
		chanId string            = "test-channel"
		access domain.AccessSpec = domain.AccessSpec{[]domain.Scope{{"RW", domain.DevType, devId}, {"R", domain.ChanType, chanId}}}
	)

	devKey, _ := services.AddDeviceKey(user.Id, devId, user.MasterKey, access)
	usrKey, _ := services.AddUserKey(user.Id, user.MasterKey, access)

	cases := []struct {
		domain.AccessRequest
		granted bool
	}{
		{domain.AccessRequest{"R", domain.UserType, user.Id, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"W", domain.UserType, user.Id, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"X", domain.UserType, user.Id, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"R", domain.DevType, devId, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"W", domain.DevType, devId, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"X", domain.DevType, devId, user.Id, user.MasterKey}, true},
		{domain.AccessRequest{"R", domain.ChanType, chanId, devId, user.MasterKey}, true},
		{domain.AccessRequest{"W", domain.ChanType, chanId, devId, user.MasterKey}, true},
		{domain.AccessRequest{"X", domain.ChanType, chanId, devId, user.MasterKey}, true},

		// illegal master key access (resources not owned)
		{domain.AccessRequest{"R", domain.UserType, "bad", user.Id, user.MasterKey}, false},
		{domain.AccessRequest{"R", domain.DevType, "bad", user.Id, user.MasterKey}, false},
		{domain.AccessRequest{"R", domain.ChanType, "bad", "bad", user.MasterKey}, false},

		{domain.AccessRequest{"R", domain.UserType, user.Id, user.Id, usrKey}, false},
		{domain.AccessRequest{"R", domain.DevType, devId, user.Id, usrKey}, true},
		{domain.AccessRequest{"W", domain.DevType, devId, user.Id, usrKey}, true},
		{domain.AccessRequest{"X", domain.DevType, devId, user.Id, usrKey}, false},
		{domain.AccessRequest{"R", domain.ChanType, chanId, devId, usrKey}, true},
		{domain.AccessRequest{"W", domain.ChanType, chanId, devId, usrKey}, false},

		// device key requests
		{domain.AccessRequest{"R", domain.UserType, user.Id, user.Id, devKey}, false},
		{domain.AccessRequest{"R", domain.DevType, devId, user.Id, devKey}, true},
		{domain.AccessRequest{"W", domain.DevType, devId, user.Id, devKey}, true},
		{domain.AccessRequest{"X", domain.DevType, devId, user.Id, devKey}, false},
		{domain.AccessRequest{"R", domain.ChanType, chanId, devId, devKey}, true},
		{domain.AccessRequest{"W", domain.ChanType, chanId, devId, devKey}, false},

		// missing or invalid request data
		{domain.AccessRequest{"illegal-action", domain.DevType, devId, user.Id, devKey}, false},
		{domain.AccessRequest{"R", "illegal-resource", devId, user.Id, devKey}, false},
		{domain.AccessRequest{"R", domain.DevType, "unknown-id", user.Id, devKey}, false},
		{domain.AccessRequest{"R", domain.DevType, devId, user.Id, "invalid-key"}, false},
		{domain.AccessRequest{"R", domain.ChanType, devId, devId, devKey}, false},
	}

	for i, c := range cases {
		err := services.CheckPermissions(&c.AccessRequest)
		granted := err == nil

		if c.granted == !granted {
			t.Errorf("case %d: expected granted %t got %t", i+1, c.granted, granted)
		}
	}
}
