package SQL

import (
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/devinmcgloin/morph/src/views/common"
)

func GetUserProfileView(UserName string) (common.UserProfileView, error) {
	var userProfileView common.UserProfileView

	query, _, err := users.Where(sq.Eq{"u_id": UserName}).ToSql()
	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	log.Println(query)
	err = db.Get(&userProfileView, query)

	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	log.Println(userProfileView)

	var userProfileImages []common.SingleImg
	query, _, err = singleImg.Where(sq.Eq{"images.u_id": userProfileView.UID,
		"sources.s_size": "orig"}).ToSql()
	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	err = db.Select(&userProfileImages, query)

	if err != nil {
		log.Println(err)
		return common.UserProfileView{}, err
	}

	userProfileView.Images = userProfileImages

	return userProfileView, nil
}

func GetFeatureSingleImgView(IID uint64) (common.SingleImgView, error) {
	var singleImgView common.SingleImgView

	query, _, err := images.Where(sq.Eq{"images.i_id": IID, "sources.s_size": "orig"}).ToSql()
	if err != nil {
		log.Println(err)
		return common.SingleImgView{}, err
	}
	// `SELECT * FROM images
	// 					INNER JOIN users
	// 					ON images.u_id=users.u_id
	// 					INNER JOIN sources
	// 					ON images.i_id=sources.i_id
	// 					WHERE images.i_id = ? AND
	// 								sources.s_size = ?`

	err = db.Get(&singleImgView, query, IID, "orig")

	if err != nil {
		log.Println(err)
		return common.SingleImgView{}, err
	}
	return singleImgView, nil
}

func GetCollectionViewByLocation(LID uint64) (common.TagCollectionView, error) {
	// var locCollectionView CollectionView
	//
	// query := `SELECT * FROM images
	// 					INNER JOIN users
	// 					ON images.u_id=users.u_id
	// 					INNER JOIN sources
	// 					ON images.i_id=sources.i_id
	// 					WHERE images.l_id = ? AND
	// 								sources.s_size=?`
	//
	// err := db.Select(&locCollectionView.Images, query, LID, "orig")
	//
	// if err != nil {
	// 	log.Println(err)
	// 	return CollectionView{}, err
	// }
	//
	// locCollectionView.Query = string(LID)
	// locCollectionView.Type = "loc"
	//
	// return locCollectionView, nil
	return common.TagCollectionView{}, nil
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
	// `SELECT * FROM images
	// 				 INNER JOIN users
	// 				 ON images.u_id=users.u_id
	// 				 INNER JOIN sources
	// 				 ON images.i_id=sources.i_id
	// 				 WHERE sources.s_size=?
	// 				 ORDER BY images.i_publish_time
	// 				 LIMIT ?`

	err = db.Select(&images, query, "orig")

	if err != nil {
		log.Println(err)
		return common.CollectionView{}, err
	}

	imgCollectionView.Images = images
	return imgCollectionView, nil
}
