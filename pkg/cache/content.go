package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

const prefix = "cache:"

// Get returns the data cached at the key string and throws an error otherwise.
func Get(pool *redis.Pool, key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", prefix+key))
	if err != nil {
		return []byte{}, errors.Wrap(err, "redis unable to get cached value")
	}

	return b, nil
}

func Invalidate(pool *redis.Pool, key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", prefix+key)
	if err != nil {
		return errors.Wrap(err, "unable to delete redis cached value")
	}
	return nil
}

func ExpireAt(pool *redis.Pool, key string, t time.Duration) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("EXPIRE", prefix+key, t.Seconds())
	if err != nil {
		return errors.Wrap(err, "unable to set expiration for redis cached value")
	}
	return nil
}

func Set(pool *redis.Pool, key string, content []byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", prefix+key, content)
	if err != nil {
		return errors.Wrap(err, "unable to set redis cached value")
	}
	return nil
}

func Setex(pool *redis.Pool, key string, content []byte, t time.Duration) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", prefix+key, t.Seconds(), content)
	if err != nil {
		return errors.Wrap(err, "unable to set/expiration for redis cached value")
	}
	return nil
}
