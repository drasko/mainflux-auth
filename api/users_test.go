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

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"username":"test","password":"test"}`, http.StatusCreated},
		{"malformed body", http.StatusBadRequest},
		{`{"username":"","password":"test"}`, http.StatusBadRequest},
		{`{"username":"test","password":""}`, http.StatusBadRequest},
		{`{"username":"test","password":"test"}`, http.StatusConflict},
	}

	url := fmt.Sprintf("%s/users", ts.URL)

	for i, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}
	}
}

func TestLoginUser(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"username":"test","password":"test"}`, http.StatusCreated},
		{"malformed body", http.StatusBadRequest},
		{`{"username":"","password":""}`, http.StatusBadRequest},
		{`{"username":"","password":"test"}`, http.StatusBadRequest},
		{`{"username":"test","password":""}`, http.StatusBadRequest},
		{`{"username":"bad","password":"test"}`, http.StatusForbidden},
		{`{"username":"test","password":"bad"}`, http.StatusForbidden},
	}

	url := fmt.Sprintf("%s/sessions", ts.URL)

	for i, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		defer res.Body.Close()
	}
}

func TestAddUserKey(t *testing.T) {
	cases := []struct {
		header string
		path   string
		body   string
		code   int
	}{
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"RW","type":"device","id":"*"}]}`, http.StatusCreated},
		{user.MasterKey, user.Id, "malformed body", http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[]}`, http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"bad"}]}`, http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"RW","type":"bad","id":"*"}]}`, http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"RW","type":"device"}]}`, http.StatusBadRequest},
		{"bad", user.Id, `{"scopes":[{"actions":"RW","type":"device","id":"*"}]}`, http.StatusForbidden},
		{user.MasterKey, "bad", `{"scopes":[{"actions":"RW","type":"device","id":"*"}]}`, http.StatusNotFound},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/users/%s/api-keys", ts.URL, c.path)
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

func TestReadUserKeys(t *testing.T) {
	// create test objects
	oneKeyUser, _ := services.RegisterUser("one", "one")
	access := domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, "dev"}}}
	services.AddUserKey(oneKeyUser.Id, oneKeyUser.MasterKey, access)
	noKeysUser, _ := services.RegisterUser("empty", "empty")

	keyList := domain.KeyList{}

	cases := []struct {
		header string
		path   string
		code   int
		total  int
	}{
		{oneKeyUser.MasterKey, oneKeyUser.Id, http.StatusOK, 1},
		{noKeysUser.MasterKey, noKeysUser.Id, http.StatusOK, 0},
		{"bad", user.Id, http.StatusForbidden, 0},
		{user.MasterKey, "bad", http.StatusNotFound, 0},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/users/%s/api-keys", ts.URL, c.path)
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
			if err = json.Unmarshal(body, &keyList); err != nil {
				t.Errorf("case %d: failed to unmarshal JSON", i+1)
			}

			if len(keyList.Keys) != c.total {
				t.Errorf("case %d: expected list to have %d elements, got %d", i+1, c.total, len(keyList.Keys))
			}
		}

		defer res.Body.Close()
	}
}
