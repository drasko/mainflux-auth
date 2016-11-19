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
		username string = "user-key-username"
		password string = "user-key-password"
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

	keys := domain.KeyList{}

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
