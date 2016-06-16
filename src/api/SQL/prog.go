package SQL

import "log"

func AddSrc(src ImgSource) error {
	log.Println(src)

	sql := `
      INSERT INTO sources
    Value (
      :s_id,
      :i_id,
      :s_url,
      :s_resolution,
      :s_width,
      :s_height,
      :s_size,
      :s_file_type
      );`

	_, err := db.NamedExec(sql, src)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func AddUser(user User) error {
	_, err := db.NamedExec(`
			INSERT INTO images
		VALUES (
			:u_id,
			:u_username,
			:u_email,
			:u_first_name,
			:u_last_name,
			:u_bio,
			:u_loc,
			:u_avatar_url
			);`, user)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetImageUrl(iID uint64, kind string) (ImgSource, error) {
	var origImg ImgSource

	err := db.Get(&origImg, "SELECT * FROM sources WHERE i_id = ?, s_size = ?", iID, kind)
	if err != nil {
		return ImgSource{}, err
	}
	return origImg, nil
}

func ExistsSID(SID uint64) bool {
	var count int

	query := `SELECT count(*) FROM sources WHERE s_id = ?`
	db.Get(&count, query, SID)
	if count == 0 {
		return false
	}
	return true
}

func GetImg(iID uint64) (Img, error) {

	var img Img

	err := db.Get(&img, "SELECT * FROM images WHERE i_id = ?", iID)
	if err != nil {
		log.Println(err)
		return Img{}, err
	}

	return img, nil
}

func AddImg(img Img) error {

	_, err := db.NamedExec(`
			INSERT INTO images
		VALUES (
			:i_id,
			:i_title,
			:i_desc,
			:i_aperture,
			:i_exposure_time,
			:i_focal_length,
			:i_iso,
			:i_orientation,
			:i_camera_body,
			:i_lens,
			:i_tag_1,
			:i_tag_2,
			:i_tag_3,
			:a_id,
			:i_capture_time,
			:i_publish_time,
			:i_direction,
			:u_id,
			:l_id
			);`, img)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func UpdateImg(img Img) error {

	sql := `UPDATE images
				SET
				i_title = :i_title,
				i_desc = :i_desc,
				i_aperture = :i_aperture,
				i_exposure_time = :i_exposure_time,
				i_focal_length = :i_focal_length,
				i_iso = :i_iso,
				i_camera_body = :i_camera_body,
				i_lens = :i_lens,
				i_tag_1 = :i_tag_1,
				i_tag_2 = :i_tag_2,
				i_tag_3 = :i_tag_3
				WHERE i_id = :i_id`

	_, err := db.NamedExec(sql, img)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ExistsIID(IID uint64) bool {
	var count int

	query := `SELECT count(*) FROM images WHERE i_id = ?`
	db.Get(&count, query, IID)
	if count == 0 {
		return false
	}
	return true
}
