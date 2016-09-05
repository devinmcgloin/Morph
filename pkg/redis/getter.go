package redis

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
)

type RedisType int

const (
	Int RedisType = iota
	Bool
	String
	StringSet
	StringOrdSet
	RefSet
	RefOrdSet
	Ref
)

func GetItems(m map[string]RedisType) (map[string]interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	values := make(map[string]interface{})
	var iterOrder []string

	for key, t := range m {
		iterOrder = append(iterOrder, key)
		switch t {
		case Int:
			fallthrough
		case Bool:
			fallthrough
		case String:
			fallthrough
		case Ref:
			conn.Send("GET", key)
		case StringSet:
			fallthrough
		case RefSet:
			conn.Send("SMEMBERS", key)
		case StringOrdSet:
			fallthrough
		case RefOrdSet:
			conn.Send("ZRANGE", key, 0, -1)
		}
	}

	err := conn.Flush()
	if err != nil {
		log.Println(err)
		return values, err
	}

	for _, key := range iterOrder {
		t := m[key]
		value, err := conn.Receive()
		if err != nil {
			log.Println(err)
			return values, err
		}

		if value == nil {
			continue
		}

		switch t {
		case Int:
			values[key], err = redis.Int(value, nil)
			if err != nil {
				log.Println(err)
				return make(map[string]interface{}), err
			}
		case Bool:
			values[key], err = redis.Bool(value, nil)
			if err != nil {
				log.Println(err)
				return make(map[string]interface{}), err
			}
		case String:
			values[key], err = redis.String(value, nil)
			if err != nil {
				log.Println(err)
				return make(map[string]interface{}), err
			}
		case StringSet:
			fallthrough
		case StringOrdSet:
			values[key], err = redis.Strings(value, nil)
			if err != nil {
				log.Println(err)
				return make(map[string]interface{}), err
			}
		case Ref:
			str, err := redis.String(value, nil)
			if err != nil {
				log.Println(err)
				return make(map[string]interface{}), err
			}
			values[key] = refs.GetRedisRef(str)
		case RefSet:
			fallthrough
		case RefOrdSet:
			strings, err := redis.Strings(value, nil)
			if err != nil {
				log.Println(err)
				return make(map[string]interface{}), err
			}
			values[key] = refs.GetRedisRefs(strings)
		}
	}

	return values, nil
}

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
