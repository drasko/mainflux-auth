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
		hasErr bool
	}{
		// master key request access to user
		{domain.AccessRequest{Action: "R", Resource: "user", Id: user.Id, Key: user.MasterKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "user", Id: user.Id, Key: user.MasterKey}, false},
		{domain.AccessRequest{Action: "X", Resource: "user", Id: user.Id, Key: user.MasterKey}, false},

		// master key request access to device
		{domain.AccessRequest{Action: "R", Resource: "device", Id: devId, Key: user.MasterKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "device", Id: devId, Key: user.MasterKey}, false},
		{domain.AccessRequest{Action: "X", Resource: "device", Id: devId, Key: user.MasterKey}, false},

		// master key request access to channel
		{domain.AccessRequest{Action: "R", Resource: "channel", Id: chanId, Device: devId, Key: user.MasterKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "channel", Id: chanId, Device: devId, Key: user.MasterKey}, false},
		{domain.AccessRequest{Action: "X", Resource: "channel", Id: chanId, Device: devId, Key: user.MasterKey}, false},

		// user key request access to user
		{domain.AccessRequest{Action: "R", Resource: "user", Id: user.Id, Key: usrKey}, true},

		// user key request access to device
		{domain.AccessRequest{Action: "R", Resource: "device", Id: devId, Key: usrKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "device", Id: devId, Key: usrKey}, false},
		{domain.AccessRequest{Action: "X", Resource: "device", Id: devId, Key: usrKey}, true},

		// user key request access to channel
		{domain.AccessRequest{Action: "R", Resource: "channel", Id: chanId, Device: devId, Key: usrKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "channel", Id: chanId, Device: devId, Key: usrKey}, true},

		// device key request access to user
		{domain.AccessRequest{Action: "R", Resource: "user", Id: user.Id, Key: devKey}, true},

		// device key request access to device
		{domain.AccessRequest{Action: "R", Resource: "device", Id: devId, Key: devKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "device", Id: devId, Key: devKey}, false},
		{domain.AccessRequest{Action: "X", Resource: "device", Id: devId, Key: devKey}, true},

		// device key request access to channel
		{domain.AccessRequest{Action: "R", Resource: "channel", Id: chanId, Device: devId, Key: devKey}, false},
		{domain.AccessRequest{Action: "W", Resource: "channel", Id: chanId, Device: devId, Key: devKey}, true},

		// missing or invalid request data
		{domain.AccessRequest{Action: "illegal-action", Resource: "device", Id: devId, Key: devKey}, true},
		{domain.AccessRequest{Action: "R", Resource: "illegal-resource", Id: devId, Key: devKey}, true},
		{domain.AccessRequest{Action: "R", Resource: "device", Id: "unknown-id", Key: devKey}, true},
		{domain.AccessRequest{Action: "R", Resource: "device", Id: devId, Key: "invalid-key"}, true},
		{domain.AccessRequest{Action: "R", Resource: "channel", Id: devId, Key: devKey}, true},
	}

	for i, c := range cases {
		err := services.CheckPermissions(c.AccessRequest)

		if c.hasErr && err == nil {
			t.Errorf("case %d: expected an error to be returned", i+1)
		}

		if !c.hasErr && err != nil {
			t.Errorf("case %d: didn't expect an error to be returned", i+1)
		}
	}
}
