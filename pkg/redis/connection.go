package redis

import (
	"log"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

func init() {
	redisURL = os.Getenv("REDIS_URL")
	redisPass = os.Getenv("REDIS_PASS")

	if redisURL == "" {
		log.Fatal("REDIS_URL not set")
	}

	if redisPass == "" {
		log.Fatal("REDIS_PASS not set")
	}
	newPool()
}

var redisURL string
var redisPass string
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
			if _, err = c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
