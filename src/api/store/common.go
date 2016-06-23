package store

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func get(ds *MgoStore, ID mgo.DBRef) (interface{}, error) {
	session := ds.getSession()
	defer session.Close()

	var document interface{}

	c := session.DB(dbname).C(ID.Collection)

	err := c.FindId(ID.Id).One(&document)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Not found")
	}

	log.Println(document)

	resolveRefs(document)

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

func resolveRefs(document interface{}) (interface{}, error) {
	m := document.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case mgo.DBRef:
			fmt.Println(k, "is a DBRef", vv)
		case []mgo.DBRef:
			fmt.Println(k, "is an array of DBRef:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	return nil, nil
}

func delete(ds *MgoStore, ID mgo.DBRef) error {
	session := ds.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	return c.RemoveId(ID.Id)

}

func modify(ds *MgoStore, ID mgo.DBRef, changes bson.M) error {
	session := ds.getSession()
	defer session.Close()

	c := session.DB(dbname).C(ID.Collection)

	return c.UpdateId(ID.Id, changes)
}

// Checks to see if the operations used are allowed operations.
func verifyModificationOps(ops bson.M) error {
	var validOps = []string{"$inc", "$set", "$unset", "$push", "$pull"}

	for _, v := range ops {
		valid := in(v.(string), validOps)
		if !valid {
			return errors.New(fmt.Sprintf("Operation %s is not valid", v))
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

func link(ds *MgoStore, actor mgo.DBRef, recipient mgo.DBRef, relation string, link bool) error {
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

func modifyRef(ds *MgoStore, storeID mgo.DBRef, RefID mgo.DBRef, add bool) error {
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
