package services_test

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	dockertest "gopkg.in/ory-am/dockertest.v2"

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
			auth := err.(*domain.AuthError)

			if auth.Code != c.code {
				t.Errorf("expected %d got %d", c.code, auth.Code)
			}
		}
	}
}
