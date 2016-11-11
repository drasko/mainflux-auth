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

func TestFetchDeviceKeys(t *testing.T) {
	oneKeyUser, _ := services.RegisterUser("one-dev", "one-dev")
	devId := "fetch-device-id"
	access := domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, devId}}}
	services.AddDeviceKey(oneKeyUser.Id, devId, oneKeyUser.MasterKey, access)

	noKeysUser, _ := services.RegisterUser("empty-dev", "empty-dev")

	cases := []struct {
		userId string
		devId  string
		key    string
		code   int
		total  int
	}{
		{oneKeyUser.Id, devId, oneKeyUser.MasterKey, 0, 1},
		{noKeysUser.Id, devId, noKeysUser.MasterKey, 0, 0},
		{oneKeyUser.Id, devId, "bad", http.StatusForbidden, 0},
		{"bad", devId, oneKeyUser.MasterKey, http.StatusNotFound, 0},
		{oneKeyUser.Id, "bad", oneKeyUser.MasterKey, http.StatusNotFound, 0},
	}

	for i, c := range cases {
		list, err := services.FetchDeviceKeys(c.userId, c.devId, c.key)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
			continue
		}

		if len(list.Keys) != c.total {
			t.Errorf("case %d: expected %d items got %d", i+1, c.total, len(list.Keys))
		}
	}
}
