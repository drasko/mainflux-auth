package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddKey(t *testing.T) {
	var (
		username string = "add-key-username"
		password string = "add-key-password"
	)

	user, _ := services.RegisterUser(username, password)

	cases := []struct {
		header string
		body   string
		code   int
	}{
		{user.MasterKey, `{"owner":"test","scopes":[{"actions":"RW","type":"device","id":"1"}]}`, http.StatusCreated},
		{user.MasterKey, `{"owner":"test"}`, http.StatusCreated},
		{user.MasterKey, "malformed body", http.StatusBadRequest},
		{user.MasterKey, `{"owner":""}`, http.StatusBadRequest},
		{user.MasterKey, `{"owner":"test","scopes":[{"actions":"bad"}]}`, http.StatusBadRequest},
		{user.MasterKey, `{"owner":"test","scopes":[{"actions":"RW","type":"bad","id":"1"}]}`, http.StatusBadRequest},
		{user.MasterKey, `{"owner":"test","scopes":[{"actions":"RW","type":"device"}]}`, http.StatusBadRequest},
		{"bad", `{"owner":"test","scopes":[{"actions":"RW","type":"device","id":"1"}]}`, http.StatusForbidden},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/api-keys", ts.URL)
		b := strings.NewReader(c.body)

		req, _ := http.NewRequest("POST", url, b)
		req.Header.Set("Authorization", "Bearer "+c.header)
		req.Header.Set("Content-Type", "application/json")

		cli := &http.Client{}
		res, err := cli.Do(req)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		defer res.Body.Close()
	}
}

func TestFetchKeys(t *testing.T) {
	oneKeyUser, _ := services.RegisterUser("one-usr", "one-usr")
	noKeysUser, _ := services.RegisterUser("empty-user", "empty-usr")

	services.AddKey(oneKeyUser.MasterKey, domain.KeySpec{Owner: "owner"})

	cases := []struct {
		header string
		code   int
		total  int
	}{
		{oneKeyUser.MasterKey, http.StatusOK, 1},
		{noKeysUser.MasterKey, http.StatusOK, 0},
		{"bad", http.StatusForbidden, 0},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/api-keys", ts.URL)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+c.header)
		req.Header.Set("Content-Type", "application/json")

		cli := &http.Client{}
		res, err := cli.Do(req)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		if res.StatusCode == http.StatusOK {
			keys := domain.KeyList{}
			body, _ := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(body, &keys); err != nil {
				t.Errorf("case %d: failed to unmarshal JSON", i+1)
			}

			if len(keys.Keys) != c.total {
				t.Errorf("case %d: expected %d keys, got %d", i+1, c.total, len(keys.Keys))
			}
		}

		defer res.Body.Close()
	}
}

func TestFetchKeySpec(t *testing.T) {
	expected := domain.KeySpec{"fetch-key-owner", []domain.Scope{{"RW", "device", "device-id"}}}

	user, _ := services.RegisterUser("fetch-key-user", "fetch-key-pass")
	key, _ := services.AddKey(user.MasterKey, expected)

	cases := []struct {
		header string
		key    string
		code   int
	}{
		{user.MasterKey, key, http.StatusOK},
		{"bad", key, http.StatusForbidden},
		{user.MasterKey, "bad", http.StatusNotFound},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/api-keys/%s", ts.URL, c.key)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+c.header)
		req.Header.Set("Content-Type", "application/json")

		cli := &http.Client{}
		res, err := cli.Do(req)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		if res.StatusCode == http.StatusOK {
			actual := domain.KeySpec{}
			body, _ := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(body, &actual); err != nil {
				t.Errorf("case %d: failed to unmarshal JSON", i+1)
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

		defer res.Body.Close()
	}

}
