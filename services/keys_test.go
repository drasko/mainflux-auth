package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddKey(t *testing.T) {
	var (
		username string         = "user-key-username"
		password string         = "user-key-password"
		owner    string         = "key-owner"
		access   domain.KeySpec = domain.KeySpec{owner, []domain.Scope{{"R", domain.DevType, "dev"}}}
	)

	user, _ := services.RegisterUser(username, password)

	cases := []struct {
		master string
		access domain.KeySpec
		code   int
	}{
		{user.MasterKey, access, http.StatusOK},
		{user.MasterKey, domain.KeySpec{}, http.StatusOK},
		{"bad", access, http.StatusForbidden},
		{"bad", domain.KeySpec{}, http.StatusForbidden},
	}

	for i, c := range cases {
		key, err := services.AddKey(c.master, c.access)
		if err != nil {
			auth := err.(*domain.AuthError)
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

func TestFetchKeys(t *testing.T) {
	oneKeyUser, _ := services.RegisterUser("one-usr", "one-usr")
	noKeysUser, _ := services.RegisterUser("empty-usr", "empty-usr")

	services.AddKey(oneKeyUser.MasterKey, domain.KeySpec{Owner: "owner"})

	cases := []struct {
		mKey  string
		code  int
		total int
	}{
		{oneKeyUser.MasterKey, http.StatusOK, 1},
		{noKeysUser.MasterKey, http.StatusOK, 0},
		{"bad", http.StatusForbidden, 0},
	}

	for i, c := range cases {
		keys, err := services.FetchKeys(c.mKey)
		if err != nil {
			auth := err.(*domain.AuthError)
			if auth.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, auth.Code)
			}
			continue
		}

		if len(keys.Keys) != c.total {
			t.Errorf("case %d: expected %d keys, got %d", i+1, c.total, len(keys.Keys))
		}
	}
}
