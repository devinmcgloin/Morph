package content

func (img Img) GetLocation() Location {
	return Location{}
}

func (img Img) GetImageUrl(kind string) (ImgSource, error) {
	var origImg ImgSource

	imageID := img.ID

	err := db.Get(&origImg, "SELECT * FROM sources WHERE i_id = ?, s_size = ?", imageID, kind)
	if err != nil {
		return ImgSource{}, err
	}
	return origImg, nil
}

func GetImg(iID string) (Img, error) {

	var img Img

	err := db.Get(&img, "SELECT * FROM images WHERE i_id = ?", iID)
	if err != nil {
		return Img{}, err
	}

	return img, nil
}

func AddImg(img Img) error {

	_, err := db.NamedExec(`
			INSERT INTO images (
				i_id,
				i_title,
				i_desc,
				i_aperture,
				i_exposure_time,
				i_focal_length,
				i_iso,
				i_orientation,
				i_camera_body,
				i_lens,
				i_tag_1,
				i_tag_2,
				i_tag_3,
				i_album,
				i_capture_time,
				i_publish_time,
				i_lat,
				i_lon,
				i_direction,
				i_loc,
				i_user)
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
			:i_album,
			:i_capture_time,
			:i_publish_time,
			:i_lat,
			:i_lon,
			:i_direction,
			:i_loc,
			:i_user)`,
		map[string]interface{}{
			":i_id":            img.ID,
			":i_title":         img.Title,
			":i_desc":          img.Desc,
			":i_aperture":      img.Aperture,
			":i_exposure_time": img.ExposureTime,
			":i_focal_length":  img.FocalLength,
			":i_iso":           img.ISO,
			":i_orientation":   img.Orientation,
			":i_camera_body":   img.CameraBody,
			":i_lens":          img.Lens,
			":i_tag_1":         img.TagOne,
			":i_tag_2":         img.TagTwo,
			":i_tag_3":         img.TagThree,
			":i_album":         img.Album,
			":i_capture_time":  img.CaptureTime,
			":i_publish_time":  img.PublishTime,
			":i_direction":     img.ImgDirection,
			":l_loc":           img.Location,
			":u_id":            img.User,
		})
	if err != nil {
		return err
	}
	return nil
}

func GetAlbum(albumTag string) (ImgCollection, error) {
	var collectionPage ImgCollection

	var images []Img

	err := db.Select(&images, "SELECT * FROM images WHERE i_album = ?", albumTag)
	if err != nil {
		return ImgCollection{}, err
	}

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = albumTag
	return collectionPage, nil
}

func GetAllImgs() (ImgCollection, error) {
	var collectionPage ImgCollection

	var images []Img

	err := db.Select(&images, "SELECT * FROM images ORDER BY i_publish_time DESC")
	if err != nil {
		return ImgCollection{}, err
	}

	collectionPage.Images = images
	collectionPage.NumImg = len(collectionPage.Images)
	collectionPage.Title = "Morph"
	return collectionPage, nil
}
