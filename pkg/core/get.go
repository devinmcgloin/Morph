package core

import (
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/mongo"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

func GetUser(ref model.Ref) (model.User, rsp.Response) {

	if ref.Collection == model.Users {
		return model.User{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	var user = model.User{}

	err := mongo.Get(ref, &user)
	if err != nil {
		return model.User{}, rsp.Response{Message: "User not found",
			Code: http.StatusNotFound}
	}

	return user, rsp.Response{Code: http.StatusOK}
}

func GetImage(ref model.Ref) (model.Image, rsp.Response) {
	if ref.Collection != model.Images {
		return model.Image{}, rsp.Response{Message: "Ref is of the wrong collection type",
			Code: http.StatusBadRequest}
	}

	var image model.Image

	err := mongo.Get(ref, &image)
	if err != nil {
		return model.Image{}, rsp.Response{Message: "Image not found",
			Code: http.StatusNotFound}
	}

	SetRedisImageValues(&image)

	return image, rsp.Response{Code: http.StatusOK}
}

func GetCollection(ref model.Ref) (model.Collection, rsp.Response) {
	if ref.Collection != model.Collections {
		return model.Collection{}, rsp.Response{Message: "Ref is of the wrong type",
			Code: http.StatusBadRequest}
	}

	var col model.Collection

	err := mongo.Get(ref, &col)
	if err != nil {
		return model.Collection{}, rsp.Response{Message: "Collection not found",
			Code: http.StatusNotFound}
	}

	return col, rsp.Response{Code: http.StatusOK}
}

func GetCollectionImages(ref model.Ref) ([]model.Image, rsp.Response) {
	if ref.Collection != model.Collections {
		return []model.Image{}, rsp.Response{Message: "Ref is of the wrong type",
			Code: http.StatusBadRequest}
	}

	var images []model.Image

	err := mongo.GetAll("images", bson.M{"collections": ref}, &images)
	if err != nil {
		return []model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
	}

	if len(images) == 0 {
		return []model.Image{}, rsp.Response{Code: http.StatusNotFound,
			Message: "Collection does not exist or has not uploaded any images."}
	}

	return images, rsp.Response{Code: http.StatusOK}
}

func GetUserImages(ref model.Ref) ([]model.Image, rsp.Response) {
	if ref.Collection != model.Users {
		return []model.Image{}, rsp.Response{Message: "Ref is of the wrong type",
			Code: http.StatusBadRequest}
	}

	imageRefs, err := redis.GetSortedSet(ref.GetRString(model.Images), 0, -1)
	if err != nil {
		return []model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
	}
	var images []model.Image
	for _, ref := range imageRefs {
		img, resp := GetImage(ref)
		if !resp.Ok() {
			return []model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
		}
		images = append(images, img)
	}

	if len(images) == 0 {
		return []model.Image{}, rsp.Response{Code: http.StatusNotFound,
			Message: "User does not exist or has not uploaded any images."}
	}

	return images, rsp.Response{Code: http.StatusOK}
}

func GetFeaturedImages() ([]model.Image, rsp.Response) {
	var images []model.Image

	err := mongo.GetAll("images", bson.M{"featured": true}, &images)
	if err != nil {
		return []model.Image{}, rsp.Response{Code: http.StatusInternalServerError}
	}

	if len(images) == 0 {
		return []model.Image{}, rsp.Response{Code: http.StatusNoContent,
			Message: "No featured images exist at this time."}
	}

	return images, rsp.Response{Code: http.StatusOK}
}

func SetRedisImageValues(image *model.Image) error {
	ref := model.Ref{Collection: model.Images, ShortCode: image.ShortCode}

	request := make(map[string]redis.RedisType)

	request[ref.GetRString(model.Downloads)] = redis.Int
	request[ref.GetRString(model.Views)] = redis.Int
	request[ref.GetRString(model.Purchases)] = redis.Int
	request[ref.GetRString(model.Owner)] = redis.Ref
	request[ref.GetRString(model.FavoritedBy)] = redis.RefOrdSet
	request[ref.GetRString(model.CollectionsIn)] = redis.RefOrdSet
	request[ref.GetRSet(model.Favorited)] = redis.MemberSet

	values, err := redis.GetItems(request)
	if err != nil {
		return err
	}

	image.Downloads, _ = values[ref.GetRString(model.Downloads)].(int)
	image.Views, _ = values[ref.GetRString(model.Views)].(int)
	image.Purchases, _ = values[ref.GetRString(model.Purchases)].(int)
	image.Featured, _ = values[ref.GetRString(model.Favorited)].(bool)

	str, ok := values[ref.GetRString(model.Owner)].(model.Ref)
	if ok {
		image.OwnerLink = refs.GetURL(str)
	}

	strs, ok := values[ref.GetRString(model.FavoritedBy)].([]model.Ref)
	if ok {
		image.FavoritedByLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.CollectionsIn)].([]model.Ref)
	if ok {
		image.CollectionInLinks = refs.GetURLs(strs)
	}

	return nil
}

func SetRedisCollectionValues(col *model.Collection) error {
	ref := model.Ref{Collection: model.Collections, ShortCode: col.ShortCode}

	request := make(map[string]redis.RedisType)

	request[ref.GetRString(model.Views)] = redis.Int
	request[ref.GetRString(model.ViewType)] = redis.String
	request[ref.GetRString(model.Owner)] = redis.Ref
	request[ref.GetRString(model.FavoritedBy)] = redis.RefOrdSet
	request[ref.GetRString(model.FollowedBy)] = redis.RefOrdSet
	request[ref.GetRString(model.Images)] = redis.RefOrdSet
	request[ref.GetRString(model.Featured)] = redis.MemberSet

	values, err := redis.GetItems(request)
	if err != nil {
		return err
	}

	col.Views, _ = values[ref.GetRString(model.Views)].(int)
	col.Featured, _ = values[ref.GetRSet(model.Featured)].(bool)
	col.ViewType, _ = values[ref.GetRString(model.ViewType)].(string)

	str, ok := values[ref.GetRString(model.Owner)].(model.Ref)
	if ok {
		col.OwnerLink = refs.GetURL(str)
	}

	strs, ok := values[ref.GetRString(model.FavoritedBy)].([]model.Ref)
	if ok {
		col.FavoritedByLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.FollowedBy)].([]model.Ref)
	if ok {
		col.FollowedByLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.Images)].([]model.Ref)
	if ok {
		col.ImageLinks = refs.GetURLs(strs)
	}

	return nil
}

func SetRedisUserValues(user *model.User) error {
	ref := model.Ref{Collection: model.Users, ShortCode: user.ShortCode}

	request := make(map[string]redis.RedisType)

	request[ref.GetRString(model.Images)] = redis.RefOrdSet
	request[ref.GetRString(model.Collections)] = redis.RefOrdSet
	request[ref.GetRString(model.Followed)] = redis.RefOrdSet
	request[ref.GetRString(model.Favorited)] = redis.RefOrdSet

	request[ref.GetRSet(model.Featured)] = redis.MemberSet
	request[ref.GetRSet(model.Admin)] = redis.MemberSet
	request[ref.GetRString(model.Views)] = redis.Int

	values, err := redis.GetItems(request)
	if err != nil {
		return err
	}

	user.Views, _ = values[ref.GetRString(model.Views)].(int)
	user.Admin, _ = values[ref.GetRString(model.Views)].(bool)
	user.Featured, _ = values[ref.GetRString(model.Featured)].(bool)

	strs, ok := values[ref.GetRString(model.Images)].([]model.Ref)
	if ok {
		user.ImageLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.Collections)].([]model.Ref)
	if ok {
		user.CollectionLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.Favorited)].([]model.Ref)
	if ok {
		user.FavoritedLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.Followed)].([]model.Ref)
	if ok {
		user.FollowedLinks = refs.GetURLs(strs)
	}

	return nil
}

func SetPrivateRedisUserValues(user *model.User) error {
	ref := model.Ref{Collection: model.Users, ShortCode: user.ShortCode}

	request := make(map[string]redis.RedisType)

	request[ref.GetRString(model.FollowedBy)] = redis.RefOrdSet
	request[ref.GetRString(model.FavoritedBy)] = redis.RefOrdSet
	request[ref.GetRString(model.Seen)] = redis.RefOrdSet
	request[ref.GetRString(model.Purchased)] = redis.RefOrdSet

	values, err := redis.GetItems(request)
	if err != nil {
		return err
	}

	strs, ok := values[ref.GetRString(model.FavoritedBy)].([]model.Ref)
	if ok {
		user.FavoritedByLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.FollowedBy)].([]model.Ref)
	if ok {
		user.FollowedByLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.Seen)].([]model.Ref)
	if ok {
		user.SeenLinks = refs.GetURLs(strs)
	}

	strs, ok = values[ref.GetRString(model.Purchased)].([]model.Ref)
	if ok {
		user.Purchased = refs.GetURLs(strs)
	}

	return nil
}
