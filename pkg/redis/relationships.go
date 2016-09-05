package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/sprioc/composer/pkg/model"
)

type Relationship int

const (
	Follow Relationship = iota
	Favorite
	Collection
)

func LinkItems(user model.Ref, relationship Relationship, ref model.Ref, unlink bool) error {
	var forwards model.RString
	var backwards model.RString
	var method string
	if unlink {
		method = "ZREM"
	} else {
		method = "ZADD"
	}

	if relationship == Follow {
		forwards = model.Followed
		backwards = model.FollowedBy
	} else if relationship == Favorite {
		forwards = model.Favorited
		backwards = model.FavoritedBy
	} else if relationship == Collection {
		forwards = model.Collections
		backwards = model.CollectionsIn
	} else {
		return fmt.Errorf("Invalid relationship type %T expected (Follow|Favorite|Collection)", relationship)
	}

	if user.Collection != model.Users {
		return fmt.Errorf("Invalid item type %s expected users", user.Collection)
	}

	timestamp := time.Now().Unix()

	conn := pool.Get()
	defer conn.Close()

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	err := conn.Send(method, user.GetRString(forwards), timestamp, ref.GetTag())
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.Send(method, ref.GetRString(backwards), timestamp, user.GetTag())
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.Send("EXEC")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func AddToCollection(image, collection model.Ref) error {
	conn := pool.Get()
	defer conn.Close()
	timestamp := time.Now().Unix()

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	err := conn.Send("ZADD", image.GetRString(model.CollectionsIn), timestamp, collection.GetTag())
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.Send("ZADD", collection.GetRString(model.Images), timestamp, image.GetTag())
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.Send("EXEC")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}
