package cache

import (
	"time"

	"log"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func Configure(server, password string) {
	pool = &redis.Pool{
		MaxActive:   15,
		Wait:        true,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			_, err = c.Do("AUTH", password)
			if err != nil {
				c.Close()
				log.Println(err)
				return nil, err
			}
			if _, err := c.Do("PING"); err != nil {
				c.Close()
				log.Println(err)
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
