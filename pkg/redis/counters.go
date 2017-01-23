package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
)

func IncrementCounter(ref model.Ref, counterType model.RString) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("INCR", ref.GetRString(counterType)))
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return counter, nil
}

func GetCounter(ref model.Ref, counterType model.RString) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	counter, err := redis.Int(conn.Do("GET", ref.GetRString(counterType)))
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return counter, nil
}
