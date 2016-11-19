package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

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
