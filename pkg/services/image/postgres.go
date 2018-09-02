package image

import (
	"context"
	"fmt"

	"github.com/fatih/structs"
	"github.com/fokal/fokal-core/pkg/log"
	"github.com/jmoiron/sqlx"
)

type PGImageService struct {
	db *sqlx.DB
}

func (pg *PGImageService) ImageByID(ctx context.Context, id uint64) (*Image, error) {
	image := new(Image)
	err := pg.db.GetContext(ctx, image, "SELECT * FROM content.images WHERE id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	return image, nil
}

func (pg *PGImageService) ImageByShortcode(ctx context.Context, shortcode string) (*Image, error) {
	image := new(Image)
	err := pg.db.GetContext(ctx, image, "SELECT * FROM content.images WHERE shortcode = $1", shortcode)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	return image, nil
}
func (pg *PGImageService) ExistsByShortcode(ctx context.Context, shortcode string) (bool, error) {
	var exists bool
	err := pg.db.GetContext(ctx, &exists, "SELECT count(1) FROM content.images WHERE shortcode = $1", shortcode)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return false, err
	}
	return exists, nil
}

// func (pg *PGImageService) CreateImage(ctx context.Context, i *Image) error {}
func (pg *PGImageService) DeleteImage(ctx context.Context, id uint64) error {
	log.WithContext(ctx).WithField("image-id", id).Warn("performing distructive action: deleting image")
	_, err := pg.db.ExecContext(ctx, "DELETE FROM content.images WHERE id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

// func (pg *PGImageService) RandomImage(ctx context.Context) (*Image, error)                       {

// }
// func (pg *PGImageService) RandomImageForUser(ctx context.Context, userID uint64) (*Image, error) {

// }

func (pg *PGImageService) Feature(ctx context.Context, id uint64, user uint64) error {
	log.WithContext(ctx).WithField("image-id", id).Debug("featuring image")
	_, err := pg.db.ExecContext(ctx, "UPDATE content.images SET featured = TRUE WHERE id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}
func (pg *PGImageService) UnFeature(ctx context.Context, id uint64, user uint64) error {
	log.WithContext(ctx).WithField("image-id", id).Debug("unfeaturing image")
	_, err := pg.db.ExecContext(ctx, "UPDATE content.images SET featured = FALSE WHERE id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (pg *PGImageService) ImagesForUser(ctx context.Context, id uint64) (*[]Image, error) {
	imageIDs := []uint64{}
	err := pg.db.SelectContext(ctx, &imageIDs, "SELECT * FROM content.images WHERE user_id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}

	images := make([]Image, len(imageIDs))
	for i, imageID := range imageIDs {
		img, err := pg.ImageByID(ctx, imageID)
		if err != nil {
			log.WithContext(ctx).Error(err)
			return nil, err
		}
		images[i] = *img
	}

	return &images, nil
}

func (pg *PGImageService) RecordStat(ctx context.Context, id uint64, stat StatType) error {
	_, err := pg.db.ExecContext(ctx, `IF (SELECT count(*)
      FROM CONTENT.image_stats
      WHERE date = CURRENT_DATE
            AND stat_type = $2
            AND image_id = $1) > 0
  THEN
    UPDATE content.image_stats
    SET total = total + 1
    WHERE stat_type = $2 AND image_id = $1 AND date = current_date;
  ELSE
    INSERT INTO content.image_stats VALUES ($1, current_date, $2, 1);
	END IF;`, id, stat)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return err
	}
	return nil
}

func (pg *PGImageService) ImageStats(ctx context.Context, id uint64) (*ImageStats, error) {
	stats := new(ImageStats)

	err := pg.db.Get(&stats.Favorites, "SELECT count(*) FROM content.user_favorites WHERE image_id = $1;", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	err = pg.db.Get(&stats.Views, "SELECT COALESCE(sum(total),0) FROM content.image_stats WHERE image_id = $1 AND stat_type = 'view'", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	err = pg.db.Get(&stats.Downloads, "SELECT COALESCE(sum(total),0) FROM content.image_stats WHERE image_id = $1 AND stat_type = 'downloads'", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	return stats, nil
}

func (pg *PGImageService) ImageMetadata(ctx context.Context, id uint64) (*ImageMetadata, error) {
	meta := new(ImageMetadata)
	err := pg.db.GetContext(ctx, meta, "SELECT * FROM content.image_metadata WHERE image_id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	return meta, nil
}
func (pg *PGImageService) SetImageMetadata(ctx context.Context, id uint64, metadata ImageMetadata) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		log.WithContext(ctx).Println(err)
		return err
	}

	for key, val := range structs.Map(metadata) {
		log.WithContext(ctx).Debugf("UPDATE content.image_metadata SET %s = '%s' WHERE id = %d", key, val, id)
		stmt := fmt.Sprintf("UPDATE content.image_metadata SET %s = $1 WHERE id = $2", key)
		_, err = tx.Exec(stmt, val, id)
		if err != nil {
			log.WithContext(ctx).Println(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.WithContext(ctx).Println(err)
		return err
	}

	return nil
}

func (pg *PGImageService) ImageLocation(ctx context.Context, id uint64) (*Location, error) {
	loc := new(Location)
	err := pg.db.GetContext(ctx, loc, "SELECT * FROM content.image_geo WHERE image_id = $1", id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	return loc, nil
}

func (pg *PGImageService) SetImageLocation(ctx context.Context, id uint64, location Location) error {
	tx, err := pg.db.Beginx()
	if err != nil {
		log.WithContext(ctx).Println(err)
		return err
	}

	for key, val := range structs.Map(location) {
		log.WithContext(ctx).Debugf("UPDATE content.image_geo SET %s = '%s' WHERE id = %d", key, val, id)
		stmt := fmt.Sprintf("UPDATE content.image_geo SET %s = $1 WHERE id = $2", key)
		_, err = tx.Exec(stmt, val, id)
		if err != nil {
			log.WithContext(ctx).Println(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.WithContext(ctx).Println(err)
		return err
	}

	return nil
}

func (pg *PGImageService) ImageColors(ctx context.Context, id uint64) (*[]Color, error) {
	colors := []Color{}
	rows, err := pg.db.QueryContext(ctx,
		`SELECT red,green,blue, hue,saturation,val, shade, color, pixel_fraction, score FROM content.colors AS colors
		JOIN content.image_color_bridge AS bridge ON colors.id = bridge.color_id
		WHERE bridge.image_id = $1`, id)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return nil, err
	}
	if !rows.NextResultSet() {
		return &colors, rows.Err()
	}

	for rows.Next() {
		color := Color{}
		err = rows.Scan(&color.SRGB.R, &color.SRGB.G, &color.SRGB.B,
			&color.HSV.H, &color.HSV.S, &color.HSV.V, &color.Shade, &color.ColorName,
			&color.PixelFraction, &color.Score)
		if err != nil {
			log.WithContext(ctx).Error(err)
			return &colors, err
		}
		color.Hex = "#" + color.SRGB.Hex()
		colors = append(colors, color)
	}
	return &colors, nil
}

func (pg *PGImageService) SetImageColors(ctx context.Context, id uint64, colors []Color) error {
	var colorID int64
	for _, color := range colors {
		err := pg.db.GetContext(ctx, &colorID, "SELECT id FROM content.colors "+
			"WHERE red = $1 AND green = $2 AND blue = $3", color.SRGB.R, color.SRGB.G, color.SRGB.B)
		if err != nil {
			l, a, b := color.SRGB.CIELAB()
			lab := fmt.Sprintf("(%f, %f, %f)", l, a, b)

			err = pg.db.GetContext(ctx, &colorID, `
			INSERT INTO content.colors (red, green, blue, hue, saturation, val, shade, color, cielab)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::cube) RETURNING id;`, color.SRGB.R, color.SRGB.G, color.SRGB.B,
				color.HSV.H, color.HSV.S, color.HSV.V, color.Shade, color.ColorName, lab)
			if err != nil {
				log.WithContext(ctx).Error(err)
				return err
			}
		}
		_, err = pg.db.ExecContext(ctx, `
			INSERT INTO content.image_color_bridge(image_id, color_id, pixel_fraction, score)
			VALUES ($1, $2, $3, $4)`, id, colorID, color.PixelFraction, color.Score)
		if err != nil {
			log.WithContext(ctx).Error(err)
			return err
		}
	}
	return nil
}

// func (pg *PGImageService) ImageLabels(ctx context.Context, id uint64) (*[]Label, error) {

// }
