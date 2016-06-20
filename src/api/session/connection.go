package session

import (
	"time"

	"github.com/devinmcgloin/morph/src/env"
	"github.com/garyburd/redigo/redis"
)

var redisURL = env.Getenv("REDIS_URL", "redis://localhost:6379")
var redisPass = env.Getenv("REDIS_PASS", "root")
var pool *redis.Pool

func newPool() {
	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", redisPass); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func init() {
	newPool()
}
