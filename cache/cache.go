// Package cache provides a thin wrapper on top of the 3rd party redis client
// (github.com/garyburd/redigo/redis).
package cache

import "github.com/garyburd/redigo/redis"

const maxIdle int = 10

var pool *redis.Pool

// Start starts new redis pool with allowed maximum of 10 inactive connections.
func Start(url string) {
	pool = &redis.Pool{
		MaxIdle: maxIdle,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

// Connection retrieves a redis connection from the connection pool.
func Connection() redis.Conn {
	return pool.Get()
}

// Stop terminates the redis pool.
func Stop() error {
	return pool.Close()
}
