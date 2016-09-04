package redis

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/refs"
)

func GetLogin(ref model.Ref) (map[string]string, error) {
	conn := pool.Get()
	defer conn.Close()

	m, err := redis.StringMap(conn.Do("HGET", ref.GetTag()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return m, nil
}

func GetObjectId(ref model.Ref) (bson.ObjectId, error) {
	conn := pool.Get()
	defer conn.Close()

	objectIdHex, err := redis.String(conn.Do("GET",
		ref.GetTag()))
	if err != nil {
		log.Println(err)
		return bson.ObjectId(""), err
	}
	return bson.ObjectIdHex(objectIdHex), nil
}

func GetOwner(ref model.Ref) (model.Ref, error) {
	if ref.Collection == "users" {
		return model.Ref{}, fmt.Errorf("Invalid ref type, got %s expected (images|collections)", ref.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	userTag, err := redis.String(conn.Do("GET",
		ref.GetTag()))
	if err != nil {
		log.Println(err)
		return model.Ref{}, err
	}
	return refs.GetRedisRef(userTag), nil
}

func SetViewType(ref model.Ref, viewType string) error {
	if ref.Collection != model.Collections {
		return fmt.Errorf("Invalid ref type, got %s expected collections", ref.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SET", ref.GetTag(), viewType))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateImage(user, image model.Ref, objectID bson.ObjectId) error {
	if user.Collection != model.Users {
		return fmt.Errorf("Invalid user type, got %s expected user", user.Collection)
	}
	if image.Collection != model.Images {
		return fmt.Errorf("Invalid image type, got %s expected iamges", image.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	timestamp := time.Now().Unix()

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}
	if err := conn.Send("ZADD", user.GetRString(model.Collections), timestamp, image.GetTag()); err != nil {
		log.Println(err)
		return err
	}
	if err := conn.Send("SET", image.GetRString(model.Owner), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}
	if err := conn.Send("SET", image.GetTag(), objectID.Hex()); err != nil {
		log.Println(err)
		return err
	}
	if err := conn.Send("SADD", image.GetRString(model.CanEdit), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", image.GetRString(model.CanView), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", image.GetRString(model.CanDelete), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateCollection(user, collection model.Ref, objectID bson.ObjectId) error {
	if user.Collection != model.Users {
		return fmt.Errorf("Invalid user type, got %s expected user", user.Collection)
	}
	if collection.Collection != model.Collections {
		return fmt.Errorf("Invalid collection type, got %s expected collections", collection.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	timestamp := time.Now().Unix()

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("ZADD", user.GetRString(model.Collections), timestamp, collection.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SET", collection.GetRString(model.Owner), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SET", collection.GetTag(), objectID.Hex()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", collection.GetRString(model.CanEdit), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", collection.GetRString(model.CanView), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", collection.GetRString(model.CanDelete), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func CreateUser(user model.Ref, objectID bson.ObjectId, email, password, salt string) error {
	if user.Collection != model.Users {
		return fmt.Errorf("Invalid user type, got %s expected user", user.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	timestamp := time.Now().Unix()

	var m map[string]string
	m["objectID"] = objectID.Hex()
	m["email"] = email
	m["password"] = password
	m["salt"] = salt
	m["created_at"] = strconv.FormatInt(timestamp, 10)
	m["last_modified"] = strconv.FormatInt(timestamp, 10)

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	for key, val := range m {
		if err := conn.Send("HMSET", user.GetTag(), key, val); err != nil {
			log.Println(err)
			return err
		}
	}

	if err := conn.Send("SADD", "users:emails", email); err != nil {
		log.Println(err)
		return err
	}

	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
