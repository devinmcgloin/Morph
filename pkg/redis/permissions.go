package redis

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
)

func Permissions(item model.Ref, permission model.RString, userRef model.Ref) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	isAdmin, err := IsAdmin(userShortCode)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if isAdmin {
		return isAdmin, nil
	}

	containsWildcard, err := redis.Bool(conn.Do("SISMEMBER", item.GetRString(permission),
		"*"))
	if err != nil {
		log.Println(err)
		return false, err
	}
	if containsWildcard {
		return containsWildcard, nil
	}

	isMember, err := redis.Bool(conn.Do("SISMEMBER", item.GetRString(permission), userRef.GetTag()))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return isMember, nil
}

func AddPermissions(item model.Ref, permission model.RString, userRef model.Ref) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", fmt.Sprintf("%s:%s", item.GetTag(), permission),
		fmt.Sprintf("%s", userRef.GetTag()))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func IsAdmin(ref model.Ref) (bool, error) {
	conn := pool.Get()
	defer conn.Close()

	isAdmin, err := redis.Bool(conn.Do("GET", fmt.Sprintf("%s:%s", ref.GetTag(), model.Admin)))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return isAdmin, nil
}
