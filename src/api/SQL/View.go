package SQL

import (
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/devinmcgloin/morph/src/views/common"
)

func GetUserProfileView(UserName string) (common.UserProfileView, error) {
	var userProfileView common.UserProfileView

	query, _, err := users.Where(sq.Eq{"u_username": UserName}).ToSql()
	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	log.Println(query)
	err = db.Get(&userProfileView, query, UserName)

	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	log.Println(userProfileView)

	query, _, err = singleImg.Where(sq.Eq{"images.u_id": userProfileView.UID,
		"sources.s_size": "orig"}).ToSql()
	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	err = db.Select(&userProfileView.Images, query, userProfileView.UID, "orig")

	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	return userProfileView, nil
}

func GetFeatureSingleImgView(IID uint64) (common.SingleImgView, error) {
	var singleImgView common.SingleImgView

	query, _, err := images.Where(sq.Eq{"images.i_id": IID, "sources.s_size": "orig"}).ToSql()
	if err != nil {
		log.Println(err)
		return common.SingleImgView{}, err
	}

	err = db.Get(&singleImgView, query, IID, "orig")

	if err != nil {
		log.Println(err)
		return common.SingleImgView{}, err
	}
	return singleImgView, nil
}

func GetCollectionViewByLocation(LID uint64) (common.LocCollectionView, error) {
	var locCollectionView common.LocCollectionView

	// Fetching the location data
	query, _, err := singleImg.Where(sq.Eq{"l_id": LID}).ToSql()
	if err != nil {
		log.Println(err)
		return common.LocCollectionView{}, err
	}

	err = db.Select(&locCollectionView, query, LID)
	if err != nil {
		log.Println(err)
		return common.LocCollectionView{}, err
	}

	// Fetching images with matching location
	query, _, err = singleImg.Where(sq.Eq{"images.l_id": LID, "sources.s_size": "orig"}).ToSql()
	if err != nil {
		log.Println(err)
		return common.LocCollectionView{}, err
	}

	err = db.Select(&locCollectionView.Images, query, LID, "orig")
	if err != nil {
		log.Println(err)
		return common.LocCollectionView{}, err
	}

	return locCollectionView, nil
}

func GetAlbumCollectionView(AID uint64) (common.AlbumCollectionView, error) {
	var albumCollectionView common.AlbumCollectionView

	query := `SELECT * FROM albums
						INNER JOIN users
						ON albums.u_id=users.u_id
						WHERE albums.a_id = ?`

	err := db.Get(&albumCollectionView, query, AID)

	if err != nil {
		log.Println(err)
		return common.AlbumCollectionView{}, err
	}

	var albumImages []common.SingleImg
	query = `SELECT * FROM images
					 INNER JOIN users
					 ON images.u_id=users.u_id
					 INNER JOIN sources
					 ON images.i_id=sources.i_id
           WHERE images.a_id = ? AND
					 			 sources.s_size=?`

	err = db.Select(&albumImages, query, AID, "orig")

	if err != nil {
		log.Println(err)
		return common.AlbumCollectionView{}, err
	}

	albumCollectionView.Images = albumImages

	return albumCollectionView, nil

}

func GetCollectionViewByTag(tag string) (common.TagCollectionView, error) {

	var tagCollectionView common.TagCollectionView

	// query, _, err := singleImg.
	// 	Where(sq.Eq{"sources.s_size": "tag"}).
	// 	Where(sq.Or{"i_tag_1": tag,
	// 		"i_tag_2": tag,
	// 		"i_tag_3": tag}).
	// 	ToSql()
	// if err != nil {
	// 	log.Println(err)
	// 	return common.TagCollectionView{}, err
	// }
	//
	// // `SELECT * FROM images
	// // 					INNER JOIN users
	// // 					ON images.u_id=users.u_id
	// // 					INNER JOIN sources
	// // 					ON images.i_id=sources.i_id
	// // 					WHERE (images.i_tag_1 = ?
	// // 					OR images.i_tag_2 = ?
	// // 					OR images.i_tag_3 = ?) AND
	// // 					sources.s_size=?`
	//
	// err = db.Select(&tagCollectionView.Images, query, tag, tag, tag, "orig")
	//
	// if err != nil {
	// 	log.Println(err)
	// 	return common.TagCollectionView{}, err
	// }
	//
	// tagCollectionView.Tag = tag

	return tagCollectionView, nil
}

func GetNumMostRecentImg(limit uint64, size string) (common.CollectionView, error) {
	var imgCollectionView common.CollectionView

	var images []common.SingleImg
	query, _, err := singleImg.Where(sq.Eq{"sources.s_size": "orig"}).OrderBy("images.i_publish_time").
		Limit(limit).ToSql()
	if err != nil {
		log.Println(err)
		return common.CollectionView{}, err
	}

	err = db.Select(&images, query, "orig")

	if err != nil {
		log.Println(err)
		return common.CollectionView{}, err
	}

	imgCollectionView.Images = images
	return imgCollectionView, nil
}
