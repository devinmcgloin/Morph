package redis

import "log"

// Set sets the given key to the interface value.
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
