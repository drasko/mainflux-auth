package api_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mainflux/mainflux-auth-server/api"
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

	for _, c := range cases {
		res, err := http.Get(url)
		if err != nil {
			t.Error("cannot make a request:", err)
		}

		if res.StatusCode != c.code {
			t.Errorf("expected status %d got %d", c.code, res.StatusCode)
		}

		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal("cannot read response:", err)
		}

		if c.body != string(body) {
			t.Errorf("expected response %s got %s", c.body, string(body))
		}
	}
}
