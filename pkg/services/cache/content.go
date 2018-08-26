package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

type RedisCache struct {
	rd            *redis.Pool
	prefix        string
	defaultExpiry time.Duration
}

func (rc *RedisCache) qualifiedID(key string) string {
	return rc.prefix + key
}

// Get returns the data cached at the key string and throws an error otherwise.
func (rc *RedisCache) Get(key string) ([]byte, error) {
	conn := rc.rd.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", rc.qualifiedID(key)))
	if err != nil {
		return []byte{}, errors.Wrap(err, "redis unable to get cached value")
	}

	return b, nil
}

func (rc *RedisCache) Invalidate(key string) error {
	conn := rc.rd.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", rc.qualifiedID(key))
	if err != nil {
		return errors.Wrap(err, "unable to delete redis cached value")
	}
	return nil
}

func (rc *RedisCache) Set(key string, content []byte) error {
	conn := rc.rd.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", rc.qualifiedID(key), content, rc.defaultExpiry)
	if err != nil {
		return errors.Wrap(err, "unable to set redis cached value")
	}
	return nil
}

func (rc *RedisCache) SetWithExpiry(key string, content []byte, t time.Duration) error {
	conn := rc.rd.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", rc.qualifiedID(key), t.Seconds(), content)
	if err != nil {
		return errors.Wrap(err, "unable to set/expiration for redis cached value")
	}
	return nil
}
