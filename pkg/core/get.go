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

	SetRedisValues(&image)

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

func SetRedisValues(image *model.Image) error {
	ref := model.Ref{Collection: model.Images, ShortCode: image.ShortCode}

	downloads, err := redis.GetInt(ref.GetRString(model.Downloads))
	if err != nil {
		return err
	}
	image.Downloads = downloads

	views, err := redis.GetInt(ref.GetRString(model.Views))
	if err != nil {
		return err
	}
	image.Views = views

	purchases, err := redis.GetInt(ref.GetRString(model.Purchases))
	if err != nil {
		return err
	}
	image.Purchases = purchases

	owner, err := redis.GetRef(ref.GetRString(model.Owner))
	if err != nil {
		return err
	}
	image.OwnerLink = refs.GetURL(owner)

	favoritedBy, err := redis.GetSortedSet(ref.GetRString(model.FavoritedBy), 0, -1)
	if err != nil {
		return err
	}
	image.FavoritedByLinks = refs.GetURLs(favoritedBy)

	collectionsIn, err := redis.GetSortedSet(ref.GetRString(model.CollectionsIn), 0, -1)
	if err != nil {
		return err
	}
	image.CollectionInLinks = refs.GetURLs(collectionsIn)
	return nil
}
