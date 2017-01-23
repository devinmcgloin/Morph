package redis

import (
	"errors"
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
	MemberSet
)

func getItems(m map[string]RedisType, ref model.Ref) (map[string]interface{}, error) {
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
			conn.Send("SISMEMBER", key, ref.GetTag())
		case StringOrdSet:
			fallthrough
		case RefOrdSet:
			conn.Send("ZRANGE", key, 0, -1)
		case MemberSet:
			conn.Send("SISMEMBER", key, ref.GetTag())
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
		case MemberSet:
			fallthrough
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

func GetUser(u model.Ref, priv bool) (model.User, error) {
	if u.Collection != model.Users {
		return model.User{}, errors.New("Invalid Reference Type")
	}
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", u.GetTag()))
	if err != nil {
		return model.User{}, err
	}
	var user = model.User{}
	user.ShortCode = u
	if err := redis.ScanStruct(values, &user); err != nil {
		return model.User{}, err
	}

	if priv {
		setPrivateRedisUserValues(&user)
	}
	setRedisUserValues(&user)

	return user, nil
}

func GetImage(i model.Ref) (model.Image, error) {
	if i.Collection != model.Images {
		return model.Image{}, errors.New("Invalid Reference Type")
	}
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Values(conn.Do("HGETALL", i.GetTag()))
	if err != nil {
		return model.Image{}, err
	}
	var image = model.Image{}
	image.ShortCode = i
	if err := redis.ScanStruct(values, &image); err != nil {
		return model.Image{}, err
	}
	setRedisImageValues(&image)
	IncrementCounter(i, model.Views)
	return image, nil
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

func CheckMembership(key, item string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	b, err := redis.Bool(conn.Do("SISMEMBER", key, item))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return b, nil
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

func setRedisImageValues(image *model.Image) error {
	ref := image.GetRef()

	request := make(map[string]RedisType)

	request[ref.GetRString(model.Downloads)] = Int
	request[ref.GetRString(model.Views)] = Int
	request[ref.GetRString(model.Owner)] = String

	values, err := getItems(request, image.GetRef())
	if err != nil {
		return err
	}

	image.Downloads, _ = values[ref.GetRString(model.Downloads)].(int)
	image.Views, _ = values[ref.GetRString(model.Views)].(int)

	str, ok := values[ref.GetRString(model.Owner)].(string)
	if ok {
		image.Owner = refs.GetRedisRef(str)
	}
	log.Print(str)

	return nil
}

func setRedisUserValues(user *model.User) error {
	ref := user.GetRef()

	request := make(map[string]RedisType)

	request[ref.GetRString(model.Images)] = RefOrdSet
	request[ref.GetRSet(model.Featured)] = MemberSet
	request[ref.GetRSet(model.Admin)] = MemberSet
	request[ref.GetRString(model.Views)] = Int

	values, err := getItems(request, user.GetRef())
	if err != nil {
		return err
	}

	user.Views, _ = values[ref.GetRString(model.Views)].(int)
	user.Admin, _ = values[ref.GetRSet(model.Admin)].(bool)
	user.Featured, _ = values[ref.GetRSet(model.Featured)].(bool)

	strs, ok := values[ref.GetRString(model.Images)].([]string)
	if ok {
		user.Images = refs.GetRedisRefs(strs)
	}

	return nil
}

func setPrivateRedisUserValues(user *model.User) error {
	return nil
}
