package redis

import "log"

func Set(key string, val interface{}) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
