package session

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/sprioc/conductor/pkg/generator"
	"gopkg.in/mgo.v2/bson"
)

// This REDIS follows the follwoing schema

// sessionID is a random generated ID

// recent:<sessionID> links the the string of the last page the session visited.

// session:<sessionID> links to the ObjectID of the user logged in. If not
// found, the user is not logged in.

var sessionLifetime = int((time.Minute * 10).Seconds())
var refreshAt = int((time.Minute * 1).Seconds())

func NewSessionID() string {

	conn := pool.Get()
	defer conn.Close()

	randTitle := generator.RandString(12)
	if ExistsSessionID(randTitle) {
		randTitle = generator.RandString(12)
	}

	conn.Do("LPUSH", "recent:"+randTitle, "/")
	conn.Do("PEXPIRE", "recent:"+randTitle, time.Hour*2)
	return randTitle

}

func ExistsSessionID(sessionID string) bool {
	conn := pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "recent:"+sessionID))
	if err != nil {
		log.Println(err)
		return false
	}
	return exists
}

func GetUserId(sessionId string) (bson.ObjectId, error) {
	conn := pool.Get()
	defer conn.Close()

	hexUserID, err := redis.String(conn.Do("GET", "session:"+sessionId))
	if err != nil {
		log.Println(err)
		return bson.ObjectId(""), err
	}

	userID := bson.ObjectIdHex(hexUserID)

	if !userID.Valid() {
		log.Println(userID)
	}
	return userID, nil
}

func SetUserID(sessionID string, userID bson.ObjectId) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", "session:"+sessionID, userID.Hex())
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = conn.Do("EXPIRE", "session:"+sessionID, sessionLifetime)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func RefreshUserID(sessionID string) error {
	conn := pool.Get()
	defer conn.Close()

	time, err := redis.Int(conn.Do("TTL", "session:"+sessionID))
	if err != nil {
		log.Println(err)
		return err
	}
	if time < refreshAt {
		return nil
	}

	_, err = conn.Do("EXPIRE", "session:"+sessionID, sessionLifetime)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteSessionID(sessionID string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", "session:"+sessionID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
