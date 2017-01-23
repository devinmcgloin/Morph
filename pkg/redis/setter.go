package redis

import (
	"log"

	"github.com/sprioc/composer/pkg/model"
)

// Set sets the given key to the interface value.
func SetKey(key string, val interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func SetHash(key string, val interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", key, val)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ExecuteModifications(user model.Ref, changes [][]string) error {
	conn := pool.Get()
	defer conn.Close()

	for _, change := range changes {
		switch change[0] {
		case "SET":
			conn.Send("SET", change[1], change[2])
		case "DEL":
			conn.Send("DEL", change[1])
		case "ADD":
			conn.Send("SADD", change[1], change[2])
		case "RM":
			conn.Send("SREM", change[1], change[2])
		}
	}
	return conn.Flush()

}
