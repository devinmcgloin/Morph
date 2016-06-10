package SQL

import "log"

func GetFeatureSingleImgView(IID uint64) (SingleImgView, error) {
	var singleImgView SingleImgView

	query := `SELECT * FROM images
						INNER JOIN users
						ON images.u_id=users.u_id
						INNER JOIN sources
						ON images.i_id=sources.i_id
						WHERE images.i_id = ? AND
									sources.s_size = ?`

	err := db.Get(&singleImgView, query, IID, "orig")

	if err != nil {
		log.Println(err)
		return SingleImgView{}, err
	}
	return singleImgView, nil
}
