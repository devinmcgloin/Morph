package content

import "log"

func AddSrc(src ImgSource) error {
	_, err := db.NamedExec(`
      INSERT INTO images (
        i_id,
        s_id,
        s_url,
        s_resolution,
        s_width,
        s_height,
        s_size,
        s_file_type,
        )
    VALUES (
      :i_id,
      :s_id,
      :s_url,
      :s_resolution,
      :s_width,
      :s_height,
      :s_size,
      :s_file_type,
      )`,
		map[string]interface{}{
			":i_id":         src.IID,
			":s_id":         src.ID,
			":s_url":        src.URL,
			":s_resolution": src.Resolution,
			":s_width":      src.Width,
			":s_height":     src.Height,
			":s_size":       src.Size,
			":s_file_type":  src.FileType,
		})
	log.Printf("DB Error Type = %s", err.Error())
	if err != nil {
		return err
	}
	return nil
}
