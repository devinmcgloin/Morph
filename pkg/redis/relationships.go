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
)

func LinkItems(ref, user model.Ref, relationship Relationship) error {
	var forwards RString
	var backwards RString

	if relationship == Follow {
		forwards = Followed
		backwards = FollowedBy
	} else if relationship == Favorite {
		forwards = Favorited
		backwards = FavoritedBy
	} else {
		return fmt.Errorf("Invalid relationship type %T expected (Follow|Favorite)", relationship)
	}

	if user.ItemType != Users {
		return fmt.Errorf("Invalid item type %s expected users", user.ItemType)
	}

	timestamp := time.Now().Unix()

	conn := pool.Get()
	defer conn.Close()

	if err := conn.Send("MULTI"); err != nil {
		log.Println(err)
		return err
	}

	err := conn.Send("ZADD", user.GetRString(forwards), timestamp, ref.GetTag())
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.Send("ZADD", ref.GetRString(backwards), timestamp, user.GetTag())
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

	err := conn.Send("ZADD", image.GetRString(CollectionsIn), timestamp, collection.GetTag())
	if err != nil {
		log.Println(err)
		return err
	}

	err = conn.Send("ZADD", collection.GetRString(Images), timestamp, image.GetTag())
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
