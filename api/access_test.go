package api_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
)

func TestCheckCredentials(t *testing.T) {
	var (
		username string            = "access-user"
		password string            = "access-pass"
		devId    string            = "test-device"
		chanId   string            = "test-chan"
		spec     domain.AccessSpec = domain.AccessSpec{
			[]domain.Scope{
				{"R", "channel", chanId},
				{"RW", "device", devId},
			},
		}
	)

	user, _ := services.RegisterUser(username, password)
	key, _ := services.AddUserKey(user.Id, user.MasterKey, spec)

	cases := []struct {
		body string
		code int
	}{
		{fmt.Sprintf(`{"action":"R","type":"user","id":"%s","owner":"%s","key":"%s"}`, user.Id, user.Id, user.MasterKey), http.StatusOK},
		{`malformed body`, http.StatusBadRequest},
		{fmt.Sprintf(`{"action":"bad","type":"user","id":"%s","owner":"%s","key":"%s"}`, user.Id, user.Id, user.MasterKey), http.StatusBadRequest},
		{fmt.Sprintf(`{"action":"R","type":"bad","id":"%s","owner":"%s","key":"%s"}`, user.Id, user.Id, user.MasterKey), http.StatusBadRequest},

		{fmt.Sprintf(`{"action":"R","type":"user","key":"%s"}`, user.MasterKey), http.StatusBadRequest},
		{fmt.Sprintf(`{"action":"R","type":"user","id":"%s"}`, user.Id), http.StatusBadRequest},

		{fmt.Sprintf(`{"action":"R","type":"user","id":"%s","owner":"%s","key":"%s"}`, user.Id, user.Id, key), http.StatusForbidden},
		{fmt.Sprintf(`{"action":"R","type":"device","id":"%s","owner":"%s","key":"%s"}`, devId, user.Id, key), http.StatusOK},
		{fmt.Sprintf(`{"action":"W","type":"device","id":"%s","owner":"%s","key":"%s"}`, devId, user.Id, key), http.StatusOK},
		{fmt.Sprintf(`{"action":"R","type":"channel","id":"%s","owner":"%s","key":"%s"}`, chanId, devId, key), http.StatusOK},
		{fmt.Sprintf(`{"action":"R","type":"channel","id":"%s","key":"%s"}`, chanId, key), http.StatusBadRequest},
		{fmt.Sprintf(`{"action":"X","type":"channel","id":"%s","owner":"%s","key":"%s"}`, chanId, devId, key), http.StatusForbidden},
	}

	url := fmt.Sprintf("%s/access-checks", ts.URL)

	for i, c := range cases {
		b := strings.NewReader(c.body)

		res, err := http.Post(url, "application/json", b)
		if err != nil {
			t.Errorf("case %d: %s", i+1, err.Error())
		}

		if res.StatusCode != c.code {
			t.Errorf("case %d: expected status %d got %d", i+1, c.code, res.StatusCode)
		}
	}
}
