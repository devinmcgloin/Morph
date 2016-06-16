package SQL

import sq "github.com/Masterminds/squirrel"

var images = sq.Select("images.i_id", "i_title", "i_desc", "i_aperture",
	"i_exposure_time", "i_focal_length", "i_iso", "i_tag_1", "i_tag_2", "i_tag_3",
	"i_capture_time", "i_publish_time", "i_camera_body", "i_lens", "l_lat", "l_lon", "l_desc", "s_url",
	"images.u_id", "u_username", "u_first_name", "u_last_name", "u_bio",
	"u_avatar_url").From("images").Join("sources ON images.i_id = sources.i_id").
	Join("users ON users.u_id = images.u_id").Join("locations ON images.l_id = locations.l_id")

var singleImg = sq.Select("images.i_id", "i_title", "i_desc",
	"s_url").From("images").Join("sources ON images.i_id = sources.i_id")

var users = sq.Select("*").From("users").Join("locations ON users.l_id = locations.l_id")
