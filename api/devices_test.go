package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestAddDeviceKey(t *testing.T) {
	var (
		username string = "dev-key-username"
		password string = "dev-key-password"
	)

	user, _ := services.RegisterUser(username, password)

	cases := []struct {
		header string
		path   string
		body   string
		code   int
	}{
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"RW","type":"device","id":"*"}]}`, http.StatusCreated},
		{user.MasterKey, user.Id, "1", http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[]}`, http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"bad"}]}`, http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"R","type":"bad","id":"*"}]}`, http.StatusBadRequest},
		{user.MasterKey, user.Id, `{"scopes":[{"actions":"R","type":"device"}]}`, http.StatusBadRequest},
		{"bad", user.Id, `{"scopes":[{"actions":"RW","type":"device","id":"*"}]}`, http.StatusForbidden},
		{user.MasterKey, "bad", `{"scopes":[{"actions":"RW","type":"device","id":"*"}]}`, http.StatusNotFound},
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

		defer res.Body.Close()
	}
}

func TestFetchDeviceKeys(t *testing.T) {
	oneKeyUser, _ := services.RegisterUser("one-dev", "one-dev")
	devId := "fetch-device-id"
	access := domain.AccessSpec{[]domain.Scope{{"R", domain.DevType, devId}}}
	services.AddDeviceKey(oneKeyUser.Id, devId, oneKeyUser.MasterKey, access)

	noKeysUser, _ := services.RegisterUser("empty-dev", "empty-dev")

	keyList := domain.KeyList{}

	cases := []struct {
		header  string
		usrPath string
		devPath string
		code    int
		total   int
	}{
		{oneKeyUser.MasterKey, oneKeyUser.Id, devId, http.StatusOK, 1},
		{oneKeyUser.MasterKey, oneKeyUser.Id, "bad", http.StatusOK, 0},
		{noKeysUser.MasterKey, noKeysUser.Id, devId, http.StatusOK, 0},
		{"bad", oneKeyUser.Id, devId, http.StatusForbidden, 0},
		{oneKeyUser.MasterKey, "bad", devId, http.StatusNotFound, 0},
	}

	for i, c := range cases {
		url := fmt.Sprintf("%s/users/%s/devices/%s/api-keys", ts.URL, c.usrPath, c.devPath)
		req, _ := http.NewRequest("GET", url, nil)
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

		if res.StatusCode == http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			if err = json.Unmarshal(body, &keyList); err != nil {
				t.Errorf("case %d: failed to unmarshal JSON", i+1)
			}

			if len(keyList.Keys) != c.total {
				t.Errorf("case %d: expected list to have %d elements, got %d", i+1, c.total, len(keyList.Keys))
			}
		}

		defer res.Body.Close()
	}
}
