package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestCheckAccess(t *testing.T) {
	var (
		username string         = "access-user"
		password string         = "access-pass"
		device   string         = "test-device"
		spec     domain.KeySpec = domain.KeySpec{"owner", []domain.Scope{{"RW", "device", device}}}
	)

	user, _ := services.RegisterUser(username, password)
	key, _ := services.AddKey(user.MasterKey, spec)

	cases := []struct {
		header string
		body   string
		code   int
	}{
		{key, fmt.Sprintf(`{"action":"R","type":"device","id":"%s"}`, device), http.StatusOK},
		{key, fmt.Sprintf(`{"action":"W","type":"device","id":"%s"}`, device), http.StatusOK},
		{key, `malformed body`, http.StatusBadRequest},
		{key, fmt.Sprintf(`{"type":"device","id":"%s"}`, device), http.StatusBadRequest},
		{key, fmt.Sprintf(`{"action":"bad","type":"device","id":"%s"}`, device), http.StatusBadRequest},
		{key, fmt.Sprintf(`{"action":"R","type":"bad","id":"%s"}`, device), http.StatusBadRequest},
		{key, fmt.Sprintf(`{"action":"R","id":"%s"}`, device), http.StatusBadRequest},
		{key, fmt.Sprintf(`{"action":"X","type":"device","id":"%s"}`, device), http.StatusForbidden},
		{key, `{"action":"R","type":"device","id":"bad"}`, http.StatusForbidden},
	}

	url := fmt.Sprintf("%s/access-checks", ts.URL)

	for i, c := range cases {
		req, _ := http.NewRequest("POST", url, strings.NewReader(c.body))
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
