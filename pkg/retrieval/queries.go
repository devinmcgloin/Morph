package retrieval

import (
	"log"

	"fmt"

	"database/sql"
	"errors"
	"net/http"

	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/stats"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// GetUser returns the fields of a user row into a User struct, including image references.
func GetUser(state *handler.State, u int64) (model.User, error) {
	user := model.User{}
	err := state.DB.Get(&user, "SELECT * FROM content.users WHERE id = $1", u)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	images := []string{}
	err = state.DB.Select(&images, `SELECT shortcode FROM content.images WHERE user_id = $1`, u)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	imageLinks := make([]string, len(images))
	for i, v := range images {
		imageLinks[i] = model.Ref{Collection: model.Images, Shortcode: v}.ToURL(state.Port, state.Local)
	}
	user.ImageLinks = &imageLinks

	favorites := []string{}
	err = state.DB.Select(&favorites, `
	SELECT images.shortcode
	FROM content.images AS images
		JOIN content.user_favorites AS favs ON favs.image_id = images.id
	WHERE favs.user_id = $1`, u)
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	favoriteLinks := make([]string, len(favorites))
	for i, v := range favorites {
		favoriteLinks[i] = model.Ref{Collection: model.Images, Shortcode: v}.ToURL(state.Port, state.Local)
	}

	user.FavoriteLinks = &favoriteLinks
	user.Avatars = imageSources(user.Username, "avatar")

	user.Permalink = model.Ref{Collection: model.Users, Shortcode: user.Username}.ToURL(state.Port, state.Local)

	return user, nil
}

// GetUsers TODO rewrite this to make a single call to the database.
func GetUsers(state *handler.State, userIds []int64) ([]model.User, error) {
	users := []model.User{}
	for _, userId := range userIds {
		usr, err := GetUser(state, userId)
		if err != nil {
			log.Println(err)
			return []model.User{}, err
		}
		users = append(users, usr)
	}
	return users, nil
}

// GetImages TODO retwrite to make a single call to the database
func GetImages(state *handler.State, imageIDS []int64) ([]model.Image, error) {
	images := []model.Image{}
	for _, imageId := range imageIDS {
		img, err := GetImage(state, imageId)
		if err != nil {
			return []model.Image{}, err
		}
		images = append(images, img)
	}
	return images, nil
}

// GetImage takes an image ID and returns a image row into a Image struct including metadata
// and user data.
func GetImage(state *handler.State, i int64) (model.Image, error) {
	img := model.Image{}
	q := `
	SELECT id, shortcode, publish_time, last_modified, user_id, featured FROM content.images AS images
	WHERE images.id = %[1]d;

	-- metadata

	SELECT aperture, exposure_time, focal_length, iso, make, model,
	lens_make, lens_model, pixel_yd, pixel_xd, capture_time, loc, dir
	FROM content.image_metadata AS meta
	JOIN content.image_geo AS geo ON geo.image_id = meta.image_id
	WHERE meta.image_id = %[1]d;

	--landmarks


	SELECT landmark.description, landmark.location, bridge.score FROM content.image_landmark_bridge AS bridge
	JOIN content.landmarks AS landmark ON bridge.landmark_id = landmark.id
	WHERE bridge.image_id = %[1]d;

	--labels

	SELECT description, score FROM content.image_label_bridge
	JOIN content.labels ON content.image_label_bridge.label_id = content.labels.id
	WHERE content.image_label_bridge.image_id = %[1]d;

	--tags

	SELECT description FROM content.image_tags AS tags
	JOIN content.image_tag_bridge AS bridge ON tags.id = bridge.tag_id
	WHERE bridge.image_id = %[1]d;

	--colors

	SELECT red,green,blue, hue,saturation,val, shade, color, pixel_fraction, score FROM content.colors AS colors
	JOIN content.image_color_bridge AS bridge ON colors.id = bridge.color_id
	WHERE bridge.image_id = %[1]d;

	--stats
	SELECT count(*) FROM content.user_favorites
	WHERE image_id = %[1]d;

	SELECT COALESCE(sum(total),0) FROM content.image_stats
	WHERE image_id = %[1]d AND stat_type = 'view';

	SELECT COALESCE(sum(total),0) FROM content.image_stats
	WHERE image_id = %[1]d AND stat_type = 'download';
	`

	rows, err := state.DB.Queryx(fmt.Sprintf(q, i))
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&img.Id, &img.Shortcode, &img.PublishTime, &img.LastModified, &img.UserId,
			&img.Featured)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return img, handler.StatusError{Code: http.StatusNotFound, Err: errors.New("Image not found.")}
			default:
				return img, err
			}
		}
	}

	img.Metadata, err = imageMetadata(rows)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}

	img.Landmarks, err = imageLandmarks(rows)
	if err != nil {
		log.Println(err)

		return model.Image{}, err
	}

	img.Labels, err = imageLabels(rows)
	if err != nil {
		log.Println(err)
		return model.Image{}, err
	}
	img.Tags, err = imageTags(rows)
	if err != nil {
		log.Println(err)

		return model.Image{}, err
	}
	img.Colors, err = imageColors(rows)
	if err != nil {
		log.Println(err)

		return model.Image{}, err
	}

	img.Stats, err = imageStats(rows)
	if err != nil {
		log.Println(err)

		return model.Image{}, err
	}

	usr, err := GetUser(state, img.UserId)
	if err != nil {
		return model.Image{}, err
	}
	img.User = &usr
	img.Source = imageSources(img.Shortcode, "content")

	img.Permalink = model.Ref{Collection: model.Images, Shortcode: img.Shortcode}.ToURL(state.Port, state.Local)
	stats.AddStat(state.DB, i, "view")
	return img, nil
}

func imageLandmarks(rows *sqlx.Rows) ([]model.Landmark, error) {
	landmarks := []model.Landmark{}
	var err error
	if !rows.NextResultSet() {
		log.Println(rows.Err())
		return landmarks, rows.Err()
	}
	for rows.Next() {
		landmark := model.Landmark{}
		err = rows.Scan(&landmark.Description, &landmark.Location, &landmark.Score)
		if err != nil {
			log.Println(err)
		}

		landmarks = append(landmarks, landmark)
	}
	return landmarks, nil
}

func imageSources(shortcode, location string) model.ImageSource {
	var prefix = "https://images.fok.al/" + location + "/"
	var resourceBaseURL = prefix + shortcode
	return model.ImageSource{
		Raw:    resourceBaseURL,
		Large:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
		Medium: resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
		Small:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
		Thumb:  resourceBaseURL + "?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
	}
}

func imageLabels(rows *sqlx.Rows) ([]model.Label, error) {
	labels := []model.Label{}
	var err error
	if !rows.NextResultSet() {
		return labels, rows.Err()
	}

	for rows.Next() {
		label := model.Label{}
		err = rows.Scan(&label.Description, &label.Score)
		if err != nil {
			return labels, err
		}
		labels = append(labels, label)
	}

	return labels, nil
}

func imageTags(rows *sqlx.Rows) ([]string, error) {
	tags := []string{}
	var err error
	if !rows.NextResultSet() {
		return tags, rows.Err()
	}

	for rows.Next() {
		tag := ""
		err = rows.Scan(&tag)
		if err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func imageStats(rows *sqlx.Rows) (model.ImageStats, error) {
	stat := model.ImageStats{}
	var err error
	if !rows.NextResultSet() {
		return stat, rows.Err()
	}

	for rows.Next() {
		err = rows.Scan(&stat.Favorites)
		if err != nil {
			return stat, err
		}
	}
	if !rows.NextResultSet() {
		return stat, rows.Err()
	}

	for rows.Next() {
		err = rows.Scan(&stat.Views)
		if err != nil {
			return stat, err
		}
	}
	if !rows.NextResultSet() {
		return stat, rows.Err()
	}

	for rows.Next() {
		err = rows.Scan(&stat.Downloads)
		if err != nil {
			return stat, err
		}
	}
	return stat, nil
}

func imageColors(rows *sqlx.Rows) ([]model.Color, error) {
	colors := []model.Color{}
	var err error
	if !rows.NextResultSet() {
		return colors, rows.Err()
	}

	for rows.Next() {
		color := model.Color{}
		err = rows.Scan(&color.SRGB.R, &color.SRGB.G, &color.SRGB.B,
			&color.HSV.H, &color.HSV.S, &color.HSV.V, &color.Shade, &color.ColorName,
			&color.PixelFraction, &color.Score)
		if err != nil {
			log.Println(err)
			return colors, err
		}
		color.Hex = color.SRGB.Hex()
		colors = append(colors, color)
	}
	return colors, nil
}

func imageMetadata(rows *sqlx.Rows) (model.ImageMetadata, error) {
	meta := model.ImageMetadata{}
	if !rows.NextResultSet() {
		return meta, rows.Err()
	}

	for rows.Next() {
		err := rows.Scan(&meta.Aperture, &meta.ExposureTime, &meta.FocalLength, &meta.ISO, &meta.Make, &meta.Model,
			&meta.LensMake, &meta.LensModel, &meta.PixelYDimension, &meta.PixelXDimension, &meta.CaptureTime, &meta.Location, &meta.ImageDirection)
		if err != nil {
			return meta, err
		}
	}
	return meta, nil
}

func GetUserFavorites(state *handler.State, userId int64) ([]model.Image, error) {
	images := []int64{}

	err := state.DB.Select(&images, `
			SELECT favs.image_id
			FROM content.user_favorites AS favs
			INNER JOIN permissions.can_view AS view ON view.o_id = favs.image_id
				WHERE view.user_id = -1 AND favs.user_id = $1
			ORDER BY favs.created_at DESC`, userId)

	if err != nil {
		log.Println(err)
		return []model.Image{}, err
	}

	return GetImages(state, images)
}

func GetUserImages(state *handler.State, userId int64) ([]model.Image, error) {
	images := []int64{}

	err := state.DB.Select(&images, `
			SELECT images.id
			FROM content.images AS images
			INNER JOIN permissions.can_view AS view ON view.o_id = images.id
				WHERE view.user_id = -1 AND images.user_id = $1
			ORDER BY images.publish_time DESC`, userId)

	if err != nil {
		log.Println(err)
		return []model.Image{}, err
	}

	return GetImages(state, images)
}

func GetImageRef(db *sqlx.DB, i string) (model.Ref, error) {
	ref := model.Ref{Collection: model.Images, Shortcode: i}
	err := db.Get(&ref.Id, "SELECT id FROM content.images WHERE shortcode = $1", i)
	if err != nil {
		log.Printf("Error Retrieving: %v %v\n", ref, err)
		return model.Ref{}, err
	}
	return ref, nil
}

func GetUserRef(db *sqlx.DB, u string) (model.Ref, error) {
	ref := model.Ref{Collection: model.Users, Shortcode: u}
	err := db.Get(&ref.Id, "SELECT id FROM content.users WHERE username = $1", u)
	if err != nil {
		log.Printf("Error Retrieving: %v %v %v\n", u, ref, err)
		return model.Ref{}, err
	}
	return ref, nil
}

func GetUserRefByEmail(db *sqlx.DB, email string) (model.Ref, error) {
	ref := model.Ref{Collection: model.Users}
	err := db.Get(&ref, "SELECT id, username AS shortcode FROM content.users WHERE email = $1", email)
	if err != nil {
		log.Printf("Error Retrieving: %v %v %v\n", email, ref, err)
		return model.Ref{}, err
	}
	return ref, nil
}

func TaggedImages(state *handler.State, tagText string) ([]model.Image, error) {
	ids := []int64{}

	err := state.DB.Select(&ids, `
		SELECT images.id
		FROM content.image_tag_bridge AS bridge
		  JOIN content.image_tags AS tags ON bridge.tag_id = tags.id
		  JOIN content.images AS images ON bridge.image_id = images.id
		WHERE UPPER(tags.description) = UPPER($1)
		ORDER BY ranking(1, views + favorites, featured :: INT + 3) DESC;
	`, tagText)
	if err != nil {
		return []model.Image{}, err
	}
	return GetImages(state, ids)
}

func Trending(state *handler.State, limit int) ([]model.Image, error) {
	ids := []int64{}

	err := state.DB.Select(&ids, `
	SELECT id FROM content.images
	ORDER BY ranking(id, views + favorites , featured::int + 3) DESC
	LIMIT $1
	`, limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return GetImages(state, ids)

}

func FeaturedImages(state *handler.State, limit int) ([]model.Image, error) {
	imgs := []int64{}
	var stmt *sqlx.Stmt
	var err error
	stmt, err = state.DB.Preparex(`
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE view.user_id = -1 AND images.featured = TRUE
		ORDER BY publish_time DESC
		LIMIT $1
		`)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	err = stmt.Select(&imgs,
		limit)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("%+v", err)
		}
		return []model.Image{}, err
	}
	return GetImages(state, imgs)
}

func RecentImages(state *handler.State, limit int) ([]model.Image, error) {
	imageIds := []int64{}
	var err error
	err = state.DB.Select(&imageIds, `
		SELECT images.id
		FROM content.images AS images
		INNER JOIN permissions.can_view AS view ON view.o_id = images.id
		WHERE view.user_id = -1
		ORDER BY publish_time DESC
		LIMIT $1
		`, limit)

	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			log.Printf("GetRecentImages limit: %d %+v", limit, err)
		}
		return []model.Image{}, err
	}
	return GetImages(state, imageIds)
}
