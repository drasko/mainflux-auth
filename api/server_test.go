package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	ts := httptest.NewServer(Server())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/status")
	if err != nil {
		t.Error("cannot make a request:", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("expected status 200 got %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal("cannot read response:", err)
	}

	expected := `{"status": "OK"}`
	if expected != string(body) {
		t.Errorf("expected response %s got %s", expected, string(body))
	}
}
