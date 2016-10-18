package api

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/mainflux/mainflux-auth-server/utils"

	"gopkg.in/ory-am/dockertest.v2"
)

func TestMain(m *testing.M) {
	c, err := dockertest.ConnectToRedis(5, time.Second, func(url string) bool {
		return utils.StartCache(url) == nil
	})

	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// run tests
	result := m.Run()

	// close redis
	utils.CloseCache()

	// remote used image
	c.KillRemove()

	// complete the test suite
	os.Exit(result)
}
