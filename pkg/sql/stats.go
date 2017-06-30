package sql

import "log"

func AddStat(imageId int64, t string) error {
	_, err := db.Exec(`INSERT INTO content.image_stats (image_id, type)
			VALUES ($1, $2)`, imageId, t)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
