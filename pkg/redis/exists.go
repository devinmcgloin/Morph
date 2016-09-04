package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
)

func Exists(ref model.Ref) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", ref.GetTag()))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}

func ExistsEmail(email string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("SISMEMBER", "users:emails", email))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}
