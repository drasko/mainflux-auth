package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth-server/api"
)

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		in  string
		out int
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
		b := strings.NewReader(c.in)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Error("cannot make a request:", err)
		}

		if res.StatusCode != c.out {
			t.Errorf("expected status %d got %d", c.out, res.StatusCode)
		}
	}
}
