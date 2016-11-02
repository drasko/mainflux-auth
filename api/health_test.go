package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mainflux/mainflux-auth/api"
)

func TestHealthCheck(t *testing.T) {
	cases := []struct {
		body string
		code int
	}{
		{`{"status": "OK"}`, 200},
	}

	ts := httptest.NewServer(api.Server())
	defer ts.Close()

	url := ts.URL + "/status"

	for i, c := range cases {
		res, err := http.Get(url)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatalf("case %d: %s", i+1, err.Error())
		}

		if c.body != string(body) {
			t.Errorf("case %d: expected response %s got %s", i+1, c.body, string(body))
		}
	}
}
