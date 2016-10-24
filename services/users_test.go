package services_test

import (
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	dockertest "gopkg.in/ory-am/dockertest.v2"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth-server/cache"
	"github.com/mainflux/mainflux-auth-server/domain"
	"github.com/mainflux/mainflux-auth-server/services"
)

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToRedis(5, time.Second, func(url string) bool {
		cache.Start(url)
		return true
	})

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	result := m.Run()

	cache.Stop()
	c.KillRemove()
	os.Exit(result)
}

func TestRegisterUser(t *testing.T) {
	cases := []struct {
		username string
		password string
		code     int
	}{
		{"test", "test", 0},
		{"test", "test", http.StatusConflict},
		{"test", "", http.StatusBadRequest},
		{"", "test", http.StatusBadRequest},
	}

	for _, c := range cases {
		_, err := services.RegisterUser(c.username, c.password)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("expected %d got %d", c.code, auth.Code)
			}
		}
	}
}

func TestAddUserKey(t *testing.T) {
	uid, masterKey, err := fetchCredentials()
	if err != nil {
		t.Errorf("failed to retrieve created user data")
	}

	cases := []struct {
		uid     string
		key     string
		payload domain.Payload
		code    int
	}{
		{uid, masterKey, domain.Payload{}, 0},
		{uid, "invalid-key", domain.Payload{}, http.StatusForbidden},
		{"invalid-id", masterKey, domain.Payload{}, http.StatusNotFound},
		{"invalid-id", "invalid-key", domain.Payload{}, http.StatusNotFound},
	}

	for _, c := range cases {
		key, err := services.AddUserKey(c.uid, c.key, c.payload)
		if err != nil {
			auth := err.(*domain.ServiceError)
			if auth.Code != c.code {
				t.Errorf("expected %d got %d", c.code, auth.Code)
			}
			continue
		}

		if key == "" {
			t.Errorf("expected key to be created")
		}
	}
}

//
// NOTE: only one user will be created during the test container lifecycle
//
func fetchCredentials() (string, string, error) {
	c := cache.Connection()
	vals, err := c.Do("KEYS", "users:*")
	if err != nil {
		return "", "", err
	}

	cKeys, err := redis.Strings(vals, err)
	if err != nil {
		return "", "", err
	}

	cKey := cKeys[0]
	id := strings.Split(cKey, ":")[1]

	cVal, err := c.Do("HGET", cKey, "masterKey")
	if err != nil {
		return "", "", err
	}

	token, err := redis.String(cVal, err)
	if err != nil {
		return "", "", err
	}

	return id, token, nil
}
