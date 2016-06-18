package store

import (
	"log"

	"github.com/devinmcgloin/morph/src/model"
	"gopkg.in/mgo.v2/bson"
)

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

func (ds *MgoStore) GetImageByShortCode(imageShortCode string) (model.Image, error) {
	session := ds.session.Copy()
	defer session.Close()

	var image model.Image

	c := session.DB("morph").C("images")

	err := c.Find(bson.M{"shortcode": imageShortCode}).One(&image)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	return image, nil

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

func (ds *MgoStore) GetFeatureSingleImgView(imageShortCode string) (model.SingleImgView, error) {
	session := ds.session.Copy()
	defer session.Close()

	var singleImgView model.SingleImgView

	log.Println(imageShortCode)

	c := session.DB("morph").C("images")
	err := c.Find(bson.M{"shortcode": imageShortCode}).One(&singleImgView.Image)
	if err != nil {
		log.Println(err)
		return model.SingleImgView{}, err
	}

	c = session.DB("morph").C("users")
	err = c.FindId(singleImgView.User.ID).One(&singleImgView.User)
	if err != nil {
		log.Println(err)
		return model.SingleImgView{}, err
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

	var imgCollectionView model.CollectionView
	c := session.DB("morph").C("images")
	err := c.Find(nil).Sort("publish_time").Limit(limit).All(&imgCollectionView.Images)
	if err != nil {
		log.Println(err)
		return model.CollectionView{}, err
	}

	return imgCollectionView, nil
}
