package services_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/mainflux/mainflux-auth/cache"
	dockertest "gopkg.in/ory-am/dockertest.v2"
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
