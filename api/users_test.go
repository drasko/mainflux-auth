package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"username":"test","password":"test"}`, 201},
		{"malformed body", 400},
		{`{"username":"","password":"test"}`, 400},
		{`{"username":"test","password":""}`, 400},
		{`{"username":"test","password":"test"}`, 409},
	}

	url := fmt.Sprintf("%s/users", ts.URL)

	for _, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Error("cannot make a request:", err)
		}

		if res.StatusCode != c.code {
			t.Errorf("expected status %d got %d", c.code, res.StatusCode)
		}
	}
}

func TestAddUserKey(t *testing.T) {
	cases := []struct {
		header string
		path   string
		body   string
		code   int
	}{
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 201},
		{user.MasterKey, user.Id, "malformed body", 400},
		{user.MasterKey, user.Id, `{"scopes":[]}`, 400},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":""}]}`, 400},
		{"bad-key", user.Id, `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 403},
		{"", user.Id, `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 403},
		{user.MasterKey, "", `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 404},
		{user.MasterKey, "bad-id", `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 404},
	}

	for _, c := range cases {
		url := fmt.Sprintf("%s/users/%s/api-keys", ts.URL, c.path)
		b := strings.NewReader(c.body)

		req, _ := http.NewRequest("POST", url, b)
		req.Header.Set("Authorization", "Bearer "+c.header)
		req.Header.Set("Content-Type", "application/json")

		cli := &http.Client{}
		res, err := cli.Do(req)
		if err != nil {
			t.Error("cannot make a request:", err)
		}

		if res.StatusCode != c.code {
			t.Errorf("expected status %d got %d", c.code, res.StatusCode)
		}
	}
}
