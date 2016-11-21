package api_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestCheckAccess(t *testing.T) {
	var (
		username = "access-user"
		password = "access-pass"
		device   = "test-device"
		spec     = domain.KeySpec{"owner", []domain.Scope{{"CR", domain.DevType, device}}}
	)

	user, _ := services.RegisterUser(username, password)
	key, _ := services.AddKey(user.MasterKey, spec)

	cases := []struct {
		token    string
		resource string
		action   string
		code     int
	}{
		{key, fmt.Sprintf("%s/devices/%s", ts.URL, device), "GET", http.StatusOK},
		{key, fmt.Sprintf("%s/devices/%s", ts.URL, device), "POST", http.StatusOK},
		{key, fmt.Sprintf("%s/devices/%s", ts.URL, device), "PUT", http.StatusForbidden},
		{key, fmt.Sprintf("%s/devices/%s", ts.URL, device), "DELETE", http.StatusForbidden},
		{"", fmt.Sprintf("%s/devices/%s", ts.URL, device), "GET", http.StatusForbidden},
	}

	url := fmt.Sprintf("%s/access-checks", ts.URL)

	for i, c := range cases {
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Set("Authorization", "Bearer "+c.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Resource", c.resource)
		req.Header.Set("X-Action", c.action)

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
