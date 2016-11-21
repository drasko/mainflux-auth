package services_test

import (
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddKey(t *testing.T) {
	var (
		username = "add-key-username"
		password = "add-key-password"
		owner    = "add-key-owner"
		access   = domain.KeySpec{owner, []domain.Scope{{"R", domain.DevType, "dev"}}}
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

	services.AddKey(oneKeyUser.MasterKey, domain.KeySpec{Owner: "fetch-keys-owner"})

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
			authErr := err.(*domain.AuthError)
			if authErr.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, authErr.Code)
			}
			continue
		}

		if len(keys.Keys) != c.total {
			t.Errorf("case %d: expected %d keys, got %d", i+1, c.total, len(keys.Keys))
		}
	}
}

func TestFetchKeySpec(t *testing.T) {
	expected := domain.KeySpec{"fetch-key-owner", []domain.Scope{{"CR", domain.DevType, "device-id"}}}

	user, _ := services.RegisterUser("fetch-key-user", "fetch-key-pass")
	key, _ := services.AddKey(user.MasterKey, expected)

	cases := []struct {
		mKey string
		key  string
		code int
	}{
		{user.MasterKey, key, http.StatusOK},
		{"bad", key, http.StatusForbidden},
		{user.MasterKey, "bad", http.StatusNotFound},
	}

	for i, c := range cases {
		actual, err := services.FetchKeySpec(c.mKey, c.key)
		if err != nil {
			authErr := err.(*domain.AuthError)
			if authErr.Code != c.code {
				t.Errorf("case %d: expected %d got %d", i+1, c.code, authErr.Code)
			}
			continue
		}

		if actual.Owner != expected.Owner {
			t.Errorf("case %d: expected owner %s got %s", i+1, expected.Owner, actual.Owner)
		}

		if len(actual.Scopes) != len(expected.Scopes) {
			t.Errorf("case %d: mismatched number of scopes, expected %d got %d", i+1, len(expected.Scopes), len(actual.Scopes))
		}

		for i, v := range actual.Scopes {
			if expected.Scopes[i] != v {
				t.Errorf("case %d: expected %V got %V", i+1, expected.Scopes[i], v)
			}
		}
	}
}
