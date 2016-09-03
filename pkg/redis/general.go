package redis

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/composer/pkg/model"
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
	splitTag := strings.Split(userTag, ":")
	return model.Ref{
		ShortCode:  splitTag[1],
		Collection: model.Users,
	}, nil
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

func AddImage(user, image model.Ref, objectId bson.ObjectId) error {
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
	if err := conn.Send("SET", image.GetTag(), objectId.Hex()); err != nil {
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

func AddCollection(user, collection model.Ref, objectId bson.ObjectId) error {
	if user.Collection != model.Users {
		return fmt.Errorf("Invalid user type, got %s expected user", user.Collection)
	}
	if collection.Collection != model.Collections {
		return fmt.Errorf("Invalid collection type, got %s expected collections", collection.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	timestamp := time.Now().Unix()

	if err := c.Send("MULTI"); err != nil {
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

	if err := conn.Send("SET", collection.GetTag(), objectId.Hex()); err != nil {
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
