package cache

import (
	"log"

	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/getsentry/raven-go"
)

const prefix = "cache"

// Get returns the data cached at the key string and throws an error otherwise.
func Get(pool *redis.Pool, key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", prefix+key))
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "redis", "module": "cache"})
		return []byte{}, err
	}

	return b, nil
}

func Invalidate(pool *redis.Pool, key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", prefix+key)
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "redis", "module": "cache"})
		return err
	}
	return nil
}

func ExpireAt(pool *redis.Pool, key string, t time.Duration) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", prefix+key, t.Seconds())
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "redis", "module": "cache"})
		return err
	}
	return nil
}

func Set(pool *redis.Pool, key string, content []byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", prefix+key, content)
	if err != nil {
		log.Println(err)
		raven.CaptureError(err, map[string]string{"type": "redis", "module": "cache"})
		return err
	}
	return nil
}
