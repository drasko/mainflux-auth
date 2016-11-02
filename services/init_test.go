package services_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/mainflux/mainflux-auth/cache"
	"github.com/mainflux/mainflux-auth/domain"
	"github.com/mainflux/mainflux-auth/services"
	dockertest "gopkg.in/ory-am/dockertest.v2"
)

const (
	username string = "services-test-user"
	password string = "services-test-pass"
)

var user domain.User

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

	result := m.Run()

	cache.Stop()
	c.KillRemove()
	os.Exit(result)
}
