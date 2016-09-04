package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
)

func GetInt(key string) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	i, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return i, nil
}

func GetBool(key string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	b, err := redis.Bool(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return b, nil
}

func GetString(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	str, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return str, nil
}

func GetRef(key string) (model.Ref, error) {
	conn := pool.Get()
	defer conn.Close()

	str, err := redis.String(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return model.Ref{}, err
	}
	return refs.GetRedisRef(str), nil
}

func GetSortedSet(key string, start, stop int) ([]model.Ref, error) {
	conn := pool.Get()
	defer conn.Close()

	redisStrings, err := redis.Strings(conn.Do("ZRANGE", key, start, stop))
	if err != nil {
		log.Println(err)
		return []model.Ref{}, err
	}
	return refs.GetRedisRefs(redisStrings), nil
}

func GetSet(key string) ([]model.Ref, error) {
	conn := pool.Get()
	defer conn.Close()

	redisStrings, err := redis.Strings(conn.Do("GET", key))
	if err != nil {
		log.Println(err)
		return []model.Ref{}, err
	}
	return refs.GetRedisRefs(redisStrings), nil
}
