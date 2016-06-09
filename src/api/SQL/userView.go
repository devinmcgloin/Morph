package SQL

import "log"

func GetUserProfileView(UID uint64) (UserProfileView, error) {
	var userProfileView UserProfileView

	query := `SELECT * FROM users
						INNER JOIN locations
						ON users.l_id=locations.l_id
            WHERE u_id = ?`

	err := db.Get(&userProfileView, query, UID)

	if err != nil {
		log.Println(err)
		return UserProfileView{}, err
	}

	log.Println(userProfileView)

	var userProfileImages []SingleImgView
	query = `SELECT * FROM images
					 INNER JOIN sources
					 ON images.i_id=sources.i_id
					 WHERE images.u_id = ? AND sources.s_size=?`

	err = db.Select(&userProfileImages, query, UID, "orig")

	if err != nil {
		log.Println(err)
		return UserProfileView{}, err
	}

	userProfileView.Images = userProfileImages

	return userProfileView, nil
}
