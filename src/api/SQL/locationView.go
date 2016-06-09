package SQL

import "log"

func GetCollectionViewByLocation(LID uint64) (CollectionView, error) {
	var locCollectionView CollectionView

	query := `SELECT * FROM images
						INNER JOIN users
						ON images.u_id=users.u_id
						INNER JOIN sources
						ON images.i_id=sources.i_id
						WHERE images.l_id = ? AND
									sources.s_size=?`

	err := db.Select(&locCollectionView.Images, query, LID, "orig")

	if err != nil {
		log.Println(err)
		return CollectionView{}, err
	}

	locCollectionView.Query = string(LID)
	locCollectionView.Type = "loc"

	return locCollectionView, nil
}
