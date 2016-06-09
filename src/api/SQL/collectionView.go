package SQL

import "log"

func GetAlbumCollectionView(AID uint64) (AlbumCollectionView, error) {
	var albumCollectionView AlbumCollectionView

	query := `SELECT * FROM albums
						INNER JOIN users
						ON albums.u_id=users.u_id
						WHERE albums.a_id = ?`

	err := db.Get(&albumCollectionView, query, AID)

	if err != nil {
		log.Println(err)
		return AlbumCollectionView{}, err
	}

	var albumImages []SingleImgView
	query = `SELECT * FROM images
					 INNER JOIN users
					 ON images.u_id=users.u_id
           WHERE a_id = ?`

	err = db.Select(&albumImages, query, AID)

	if err != nil {
		log.Println(err)
		return AlbumCollectionView{}, err
	}

	albumCollectionView.Images = albumImages

	return albumCollectionView, nil

}

func GetCollectionViewByTag(tag string) (CollectionView, error) {

	var tagCollectionView CollectionView

	query := `SELECT * FROM images
						INNER JOIN users
						ON images.u_id=users.u_id
						WHERE images.i_tag_1 = ?
						OR images.i_tag_2 = ?
						OR images.i_tag_3 = ?`

	err := db.Select(&tagCollectionView.Images, query, tag, tag, tag)

	if err != nil {
		log.Println(err)
		return CollectionView{}, err
	}

	return tagCollectionView, nil
}

func GetCollectionViewByLocation(LID uint64) (CollectionView, error) {
	var locCollectionView CollectionView

	query := `SELECT * FROM images
						INNER JOIN users
						ON images.u_id=users.u_id
						WHERE images.LID = ?`

	err := db.Select(&locCollectionView.Images, query, LID)

	if err != nil {
		log.Println(err)
		return CollectionView{}, err
	}

	return locCollectionView, nil
}
