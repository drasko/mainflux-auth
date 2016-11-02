package api_test

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mainflux/mainflux-auth/api"
	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"

	"gopkg.in/ory-am/dockertest.v2"
)

const (
	username string = "api-test-user"
	password string = "api-test-pass"
)

var (
	ts   *httptest.Server
	user domain.User
)

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToRedis(5, time.Second, func(url string) bool {
		cache.Start(url)
		return true
	})

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// initialize globals
	user, _ = services.RegisterUser(username, password)
	ts = httptest.NewServer(api.Server())
	defer ts.Close()

	result := m.Run()

	cache.Stop()
	c.KillRemove()
	os.Exit(result)
}
