package SQL

import "log"

func GetUserProfileView(UID uint64) (UserProfileView, error) {
	var userProfileView UserProfileView

	query := `SELECT * FROM users
            WHERE u_id = ?`

	err := db.Get(&userProfileView, query, UID)

	if err != nil {
		log.Println(err)
		return UserProfileView{}, err
	}

	var userProfileImages []Img
	query = `SELECT * FROM images
            WHERE u_id = ?`

	err = db.Select(&userProfileImages, query, UID)

	if err != nil {
		log.Println(err)
		return UserProfileView{}, err
	}

	userProfileView.Images = userProfileImages

	return userProfileView, nil
}
