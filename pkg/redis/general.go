package redis

import (
	"fmt"
	"log"
	"time"

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

func CreateImage(user model.Ref, image model.Image) error {
	if user.Collection != model.Users {
		return fmt.Errorf("Invalid user type, got %s expected user", user.Collection)
	}

	conn := pool.Get()
	defer conn.Close()

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("HMSET", redis.Args{image.GetTag()}.AddFlat(image)...); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", image.GetRef().GetRString(model.CanEdit), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", image.GetRef().GetRString(model.CanView), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", image.GetRef().GetRString(model.CanDelete), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SET", image.GetRef().GetRString(model.Owner), user.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	// Adding to new images
	if err := conn.Send("LPUSH", image.GetRef().GetRSet(model.New), image.GetTag()); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("LTRIM", image.GetRef().GetRSet(model.New), 0, 501); err != nil {
		log.Println(err)
		return err
	}

	// Adding to Geo index
	if image.Location != nil {
		if err := conn.Send("GEOADD", image.GetRef().GetRSet(model.Location),
			image.Location.Coordinates[0], image.Location.Coordinates[1],
			image.GetTag()); err != nil {
			log.Println(err)
			return err
		}
	}

	_, err := conn.Do("EXEC")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func CreateUser(user model.User) error {
	conn := pool.Get()
	defer conn.Close()

	timestamp := time.Now().Unix()

	user.CreatedAt = timestamp
	user.LastModified = timestamp

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("HMSET", redis.Args{user.GetTag()}.AddFlat(user)...); err != nil {
		log.Println(err)
		return err
	}

	if err := conn.Send("SADD", "users:emails", user.Email); err != nil {
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
