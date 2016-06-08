package api

import "log"

func AddSrc(src ImgSource) error {

	sql := `
      INSERT INTO sources (
        i_id,
        s_id,
        s_url,
        s_resolution,
        s_width,
        s_height,
        s_size,
        s_file_type
        )
    VALUES (
      :i_id,
      :s_id,
      :s_url,
      :s_resolution,
      :s_width,
      :s_height,
      :s_size,
      :s_file_type
      )`

	_, err := db.NamedExec(sql, src)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetSrc(srcId string) error {
	return nil
}
