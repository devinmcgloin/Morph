package store

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/devinmcgloin/morph/src/env"
	"github.com/devinmcgloin/morph/src/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type MgoStore struct {
	session *mgo.Session
}

func NewStore() MgoStore {
	session, err := mgo.Dial(env.Getenv("MONGO_URL", "mongodb://localhost:27017/morph"))
	if err != nil {
		log.Fatal(err)
	}

	err = session.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB Session Started")
	return MgoStore{session}
}

// AddSrc takes a source struct and adds it to the db in the proper image
func (ds *MgoStore) AddSrc(imageID bson.ObjectId, src model.ImgSource) error {
	session := ds.session.Copy()
	defer session.Close()

	c := session.DB("morph").C("images")

	err := c.UpdateId(imageID, bson.M{"$push": bson.M{"sources": src}})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// AddUser takes the given user and adds it to the Database.
func (ds *MgoStore) AddUser(user model.User) error {
	session := ds.session.Copy()
	defer session.Close()

	user.ID = bson.NewObjectId()

	c := session.DB("morph").C("users")

	err := c.Insert(user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetImage takes the given image ID and gets the image object.
func (ds *MgoStore) GetImageByID(imageID bson.ObjectId) (model.Image, error) {
	session := ds.session.Copy()
	defer session.Close()

	var image model.Image

	c := session.DB("morph").C("images")

	err := c.FindId(imageID).One(&image)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	return image, nil

}

func (ds *MgoStore) GetImageByTitle(imageShortTitle string) (model.Image, error) {
	session := ds.session.Copy()
	defer session.Close()

	var image model.Image

	c := session.DB("morph").C("images")

	err := c.Find(bson.M{"short_title": imageShortTitle}).One(&image)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	return image, nil

}

func (ds *MgoStore) AddImg(image model.Image) error {
	session := ds.session.Copy()
	defer session.Close()

	log.Println(image)

	image.ID = bson.NewObjectId()

	c := session.DB("morph").C("images")

	err := c.Insert(image)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UpdateImage matches the DB image with the one passed in.
func (ds *MgoStore) UpdateImage(img model.Image) error {
	return errors.New("Not Implemented")
}

func (ds *MgoStore) GetUserProfileView(UserName string) (model.UserProfileView, error) {
	session := ds.session.Copy()
	defer session.Close()

	var userProfileView model.UserProfileView

	c := session.DB("morph").C("users")
	err := c.Find(bson.M{"username": UserName}).One(&userProfileView.User)
	if err != nil {
		log.Println(err)
		return model.UserProfileView{}, err
	}

	c = session.DB("morph").C("images")
	err = c.Find(bson.M{"user_id": userProfileView.User.ID}).All(&userProfileView.Images)
	if err != nil {
		log.Println(err)
		return model.UserProfileView{}, err
	}

	return userProfileView, nil
}

func (ds *MgoStore) GetFeatureSingleImgView(imageShortTitle string) (model.SingleImgView, error) {
	session := ds.session.Copy()
	defer session.Close()

	var singleImgView model.SingleImgView

	c := session.DB("morph").C("images")
	err := c.Find(bson.M{"short_title": imageShortTitle}).One(&singleImgView.Image)
	if err != nil {
		log.Println(err)
		return model.SingleImgView{}, nil
	}

	c = session.DB("morph").C("users")
	err = c.FindId(singleImgView.User.ID).One(&singleImgView.User)
	if err != nil {
		log.Println(err)
		return model.SingleImgView{}, nil
	}

	return singleImgView, nil
}

func (ds *MgoStore) GetCollectionViewByLocations(locationShortText ...string) (model.LocCollectionView, error) {
	session := ds.session.Copy()
	defer session.Close()

	var locCollectionView model.LocCollectionView
	_ = session.DB("morph").C("images")
	// Can use MongoDB $in tag here. Need to normalize location schema though.
	//err := c.Find(bson.M{})

	return locCollectionView, nil
}

func (ds *MgoStore) GetAlbumCollectionView(albumTitle string, userName string) (model.AlbumCollectionView, error) {
	session := ds.session.Copy()
	defer session.Close()

	var albumCollectionView model.AlbumCollectionView

	c := session.DB("morph").C("albums")
	err := c.Find(bson.M{"Title": albumTitle, "username": userName}).One(&albumCollectionView.Album)
	if err != nil {
		log.Println(err)
		return model.AlbumCollectionView{}, nil
	}

	c = session.DB("morph").C("images")
	err = c.Find(bson.M{"$in": bson.M{"_id": albumCollectionView.Album.Images}}).All(&albumCollectionView.Images)
	if err != nil {
		log.Println(err)
		return model.AlbumCollectionView{}, nil
	}

	c = session.DB("morph").C("users")
	err = c.Find(bson.M{"_id": albumCollectionView.Album.UserID}).One(&albumCollectionView.User)
	if err != nil {
		log.Println(err)
		return model.AlbumCollectionView{}, nil
	}

	return albumCollectionView, nil

}

func (ds *MgoStore) GetCollectionViewByTags(tags ...string) (model.TagCollectionView, error) {
	session := ds.session.Copy()
	defer session.Close()

	var tagCollectionView model.TagCollectionView
	c := session.DB("morph").C("images")
	err := c.Find(bson.M{"$in": bson.M{"_id": tags}}).All(&tagCollectionView.Images)
	if err != nil {
		log.Println(err)
		return model.TagCollectionView{}, nil
	}
	tagCollectionView.Tags = tags

	return tagCollectionView, nil
}

func (ds *MgoStore) GetNumMostRecentImg(limit int) (model.CollectionView, error) {
	session := ds.session.Copy()
	defer session.Close()

	log.Println("enetering num recent")

	var imgCollectionView model.CollectionView
	c := session.DB("morph").C("images")
	err := c.Find(nil).Sort("publish_time").Limit(limit).All(&imgCollectionView.Images)
	if err != nil {
		log.Println(err)
		return model.CollectionView{}, err
	}

	log.Println(imgCollectionView)

	return imgCollectionView, nil
}

func (ds *MgoStore) ExistsUserName(userName string) bool {
	return false
}

func (ds *MgoStore) ExistsAlbumTitle(userName string, albumTitle string) bool {
	return false
}

func (ds *MgoStore) ExistsImageShortTitle(userName string) bool {
	return false
}

func (ds *MgoStore) GetShortTitle() string {
	randTitle := randSeq(12)
	if ds.ExistsImageShortTitle(randTitle) {
		randTitle = randSeq(12)
	}
	return randTitle
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
