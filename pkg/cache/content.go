package cache

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

// Get returns the data cached at the key string and throws an error otherwise.
func Get(key string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return b, nil
}

func Invalidate(key string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func Cache(key string, content []byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, content)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
