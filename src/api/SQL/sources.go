package SQL

import (
	"log"

	"github.com/devinmcgloin/morph/src/api"
)

func AddSrc(src api.ImgSource) error {

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

func GetImageUrl(img api.Img, kind string) (api.ImgSource, error) {
	var origImg api.ImgSource

	imageID := img.IID

	err := db.Get(&origImg, "SELECT * FROM sources WHERE i_id = ?, s_size = ?", imageID, kind)
	if err != nil {
		return api.ImgSource{}, err
	}
	return origImg, nil
}
