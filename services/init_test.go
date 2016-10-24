package services_test

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mainflux/mainflux-auth-server/cache"
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
