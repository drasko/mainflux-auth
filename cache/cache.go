package cache

import "github.com/garyburd/redigo/redis"

const maxIdle int = 10

var pool *redis.Pool

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

func Connection() redis.Conn {
	return pool.Get()
}

func Stop() error {
	return pool.Close()
}
