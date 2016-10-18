package utils

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

const poolSize int = 10

var pool *redis.Pool

func StartCache(url string) error {
	if pool != nil {
		return nil
	}

	pool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", url)
		if err != nil {
			return nil, err
		}

		return c, err
	}, poolSize)

	if pool == nil {
		return errors.New("can't create redis pool")
	}

	return nil
}

func CacheConnection() (redis.Conn, error) {
	if pool == nil {
		return nil, errors.New("cache not initialized")
	}

	return pool.Get(), nil
}

func CloseCache() error {
	return pool.Close()
}
