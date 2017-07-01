package social

import (
	"log"

	"net/http"

	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func FavoriteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	imageRef, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.FavoriteImage(usrRef, imageRef)
}

func UnFavoriteImage(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["IID"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	imageRef, resp := core.GetImageRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.UnFavoriteImage(usrRef, imageRef)
}

func Follow(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["username"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	userRef, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.FollowUser(usrRef, userRef)
}

func UnFollow(w http.ResponseWriter, r *http.Request) rsp.Response {
	var usrRef model.Ref
	vars := mux.Vars(r)

	id := vars["username"]

	usr, ok := context.GetOk(r, "auth")
	if ok {
		usrRef = usr.(model.Ref)
	} else {
		return rsp.Response{
			Message: "Unauthorized Request, must be logged in to modify an image",
			Code:    http.StatusUnauthorized,
		}
	}

	userRef, resp := core.GetUserRef(id)
	if !resp.Ok() {
		return resp
	}

	return core.UnFollowUser(usrRef, userRef)
}

func Favorite(uID, iID int64) error {
	stmt, err := db.Preparex("INSERT INTO content.user_favorites (user_id, image_id) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(stmt, uID, iID)
}

func UnFavorite(uID, iID int64) error {
	stmt, err := db.Preparex("DELETE FROM content.user_favorites WHERE user_id = $1 AND image_id = $2")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(stmt, uID, iID)
}

func Follow(idA, idB int64) error {
	stmt, err := db.Preparex("INSERT INTO content.user_follows (user_id, user_follows) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(stmt, idA, idB)
}

func UnFollow(idA, idB int64) error {
	stmt, err := db.Preparex("DELETE FROM content.user_follows WHERE user_id = $1 AND user_follows = $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(stmt, idA, idB)
}

func AddTag(iID, tagID int64) error {
	stmt, err := db.Preparex("INSERT INTO content.image_tag_bridge (image_id, tag_id) VALUES ($1, $2)")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(stmt, iID, tagID)
}

func RemoveTag(iID, tagID int64) error {
	stmt, err := db.Preparex("DELETE FROM content.image_tag_bridge WHERE image_id = $1 AND tag_id = $2")
	if err != nil {
		log.Println(err)
		return err
	}
	return modify(stmt, iID, tagID)
}

func Feature(iID int64) error {
	_, err := db.Exec(`
	UPDATE content.images
		SET featured = TRUE
	WHERE id = $1`, iID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UnFeature(iID int64) error {
	_, err := db.Exec(`
	UPDATE content.images
		SET featured = FALSE
	WHERE id = $1`, iID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func modify(stmt *sqlx.Stmt, idA int64, idB int64) error {
	_, err := stmt.Exec(idA, idB)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
