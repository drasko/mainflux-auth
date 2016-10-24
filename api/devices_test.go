package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth-server/api"
	"github.com/mainflux/mainflux-auth-server/services"
)

func TestAddDeviceKey(t *testing.T) {
	services.RegisterUser("test-dev-api", "test-dev-api")
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
		url := fmt.Sprintf("%s/users/%s/devices/%s/api-keys", ts.URL, c.path, c.path)
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
