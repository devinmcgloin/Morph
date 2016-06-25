package store

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/devinmcgloin/sprioc/pkg/model"

	"gopkg.in/mgo.v2/bson"
)

func get(ds *MgoStore, ID model.DBRef) (interface{}, error) {
	session := ds.getSession()
	defer session.Close()

	var document interface{}

	c := session.DB(ID.Database).C(ID.Collection)

	err := c.Find(bson.M{"shortcode": ID.Shortcode}).One(&document)
	if err != nil {
		return nil, errors.New("Not found")
	}

	return document, nil
}

func create(ds *MgoStore, collection string, document interface{}) error {
	session := ds.getSession()
	defer session.Close()

	c := session.DB(dbname).C(collection)

	err := c.Insert(document)
	if err != nil {
		log.Println(err)
		return errors.New("Unable to add document to DB")
	}
	return nil
}

func delete(ds *MgoStore, ID model.DBRef) error {
	session := ds.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	return c.Remove(bson.M{"shortcode": ID.Shortcode})

}

func modify(ds *MgoStore, ID model.DBRef, changes bson.M) error {
	err := verifyModificationOps(changes)
	if err != nil {
		log.Println(err)
		return err
	}
	session := ds.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	return c.Update(bson.M{"shortcode": ID.Shortcode}, changes)
}

// Checks to see if the operations used are allowed operations.
func verifyModificationOps(ops bson.M) error {
	var validOps = []string{"$inc", "$set", "$unset", "$push", "$pull"}

	for k := range ops {
		valid := in(k, validOps)
		if !valid {
			return fmt.Errorf("Operation %s is not valid", k)
		}
	}
	return nil
}

func in(contentType string, opts []string) bool {
	for _, opt := range opts {
		if strings.Compare(contentType, opt) == 0 {
			return true
		}
	}
	return false
}

func link(ds *MgoStore, actor model.DBRef, recipient model.DBRef, relation string, link bool) error {
	forwards, backwards, err := linkType(relation)
	if err != nil {
		return fmt.Errorf("Invalid relation: %s", relation)
	}

	if strings.Compare(actor.Collection, "users") != 0 {
		return fmt.Errorf("Invalid actor: %s", actor.Collection)
	}

	var op string
	if link {
		op = "$push"
	} else {
		op = "$pull"
	}

	err = modify(ds, actor, bson.M{op: bson.M{forwards: recipient}})
	if err != nil {
		return err
	}

	err = modify(ds, recipient, bson.M{op: bson.M{backwards: actor}})
	if err != nil {
		return err
	}

	return nil
}

func modifyRef(ds *MgoStore, storeID model.DBRef, RefID model.DBRef, add bool) error {
	var op string
	if add {
		op = "$push"
	} else {
		op = "$pull"
	}

	correctStore := in(storeID.Collection, []string{"albums", "collections"})
	if !correctStore {
		return fmt.Errorf("Invalid store type %s", storeID.Collection)
	}

	err := modify(ds, storeID, bson.M{op: bson.M{RefID.Collection: RefID}})
	if err != nil {
		return err
	}
	return nil
}

func linkType(relation string) (string, string, error) {
	var forwards string
	var backwards string

	if strings.Compare("follow", relation) == 0 {
		forwards = "follows"
		backwards = "followers"
	} else if strings.Compare("favorite", relation) == 0 {
		forwards = "favorites"
		backwards = "favoriters"
	} else {
		return "", "", fmt.Errorf("Invalid link type: %s", relation)
	}
	return forwards, backwards, nil
}
