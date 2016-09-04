package mongo

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

func ValidateConnection() {
	if dbname == "" {
		log.Fatal("MONGODB_NAME must be set")
	}

	if dbURI == "" {
		log.Fatal("MONGODB_URI must be set")
	}
}

var (
	dbname = os.Getenv("MONGODB_NAME")
	dbURI  = os.Getenv("MONGODB_URI")
)

type MgoStore struct {
	session *mgo.Session
}

func ConnectStore() *MgoStore {
	log.Printf("Connect Store URI: %s DBNAME: %s", dbURI, dbname)
	session, err := mgo.Dial(dbURI)
	if err != nil {
		log.Fatal(err)
	}

	err = session.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &MgoStore{session}
}

func (ds *MgoStore) getSession() *mgo.Session {
	return ds.session.Copy()
}
