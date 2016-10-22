package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth-server/api"
	"github.com/mainflux/mainflux-auth-server/cache"
)

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"username":"test","password":"test"}`, 201},
		{"1", 400},
		{`{"username":"","password":"test"}`, 400},
		{`{"username":"test","password":""}`, 400},
		{`{"username":"test","password":"test"}`, 409},
	}

	ts := httptest.NewServer(api.Server())
	defer ts.Close()

	url := ts.URL + "/users"

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
	uid, masterKey, err := fetchCredentials()
	if err != nil {
		t.Error("failed to retrieve created user data")
	}

	cases := []struct {
		header string
		path   string
		body   string
		code   int
	}{
		{masterKey, uid, `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 201},
		{masterKey, uid, "1", 400},
		{masterKey, uid, `{"scopes":[]}`, 400},
		{masterKey, uid, `{"scopes":[{"actions":""}]}`, 400},
		{"bad-key", uid, `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 403},
		{"", uid, `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 403},
		{masterKey, "", `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 404},
		{masterKey, "bad-id", `{"scopes":[{"actions":"RW","resource":"*","id":"*"}]}`, 404},
	}

	ts := httptest.NewServer(api.Server())
	defer ts.Close()

	for _, c := range cases {
		url := ts.URL + "/users/" + c.path + "/api-keys"
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

//
// NOTE: only one user will be created during the test container lifecycle
//
func fetchCredentials() (string, string, error) {
	c := cache.Connection()
	vals, err := c.Do("KEYS", "users:*")
	if err != nil {
		return "", "", err
	}

	cKeys, err := redis.Strings(vals, err)
	if err != nil {
		return "", "", err
	}

	cKey := cKeys[0]
	id := strings.Split(cKey, ":")[1]

	cVal, err := c.Do("HGET", cKey, "masterKey")
	if err != nil {
		return "", "", err
	}

	token, err := redis.String(cVal, err)
	if err != nil {
		return "", "", err
	}

	return id, token, nil
}
