package redis

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
)

func IncrementCounter(ref model.Ref, counterType model.RString) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("INCR", fmt.Sprintf("%s:%s", ref.GetTag(), counterType)))
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return counter, nil
}

func GetCounter(ref model.Ref, counterType model.RString) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("GET", fmt.Sprintf("%s:%s", ref.GetTag(), counterType)))
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return counter, nil
}
