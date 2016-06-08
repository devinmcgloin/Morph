package SQL

import "log"

func AddSrc(src ImgSource) error {

	sql := `
      INSERT INTO sources *
    Values (
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

func GetImageUrl(iID uint64, kind string) (ImgSource, error) {
	var origImg ImgSource

	err := db.Get(&origImg, "SELECT * FROM sources WHERE i_id = ?, s_size = ?", iID, kind)
	if err != nil {
		return ImgSource{}, err
	}
	return origImg, nil
}
