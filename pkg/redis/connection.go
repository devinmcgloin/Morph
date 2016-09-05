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
	pool = newPool(redisURL, redisPass)
}

var redisURL string
var redisPass string
var pool *redis.Pool

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxActive:   15,
		Wait:        true,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("PING"); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
