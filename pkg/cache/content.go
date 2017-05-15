package cache

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

func Get(url string) ([]byte, error) {
	conn := pool.Get()
	defer conn.Close()

	b, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return b, nil
}

func Invalidate(url string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", content)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func Cache(url string, content []byte) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", content)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
