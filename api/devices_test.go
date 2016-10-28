package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestAddDeviceKey(t *testing.T) {
	cases := []struct {
		header string
		path   string
		body   string
		code   int
	}{
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"RW","resource":"device","id":"*"}]}`, 201},
		{user.MasterKey, user.Id, "1", 400},
		{user.MasterKey, user.Id, `{"scopes":[]}`, 400},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"bad"}]}`, 400},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"R","resource":"bad","id":"*"}]}`, 400},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"R","resource":"device"}]}`, 400},
		{"bad", user.Id, `{"scopes":[{"actions":"RW","resource":"device","id":"*"}]}`, 403},
		{user.MasterKey, "bad", `{"scopes":[{"actions":"RW","resource":"device","id":"*"}]}`, 404},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/users/%s/devices/%s/api-keys", ts.URL, c.path, c.path)
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
	}
}
