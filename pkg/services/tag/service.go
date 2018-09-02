package tag

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/fokal/fokal-core/pkg/services/image"

	"github.com/jmoiron/sqlx"
)

type TagStore struct {
	db    *sqlx.DB
	image image.Service
}

func New(db *sqlx.DB) *TagStore {
	return &TagStore{
		db: db,
	}
}

func (store *TagStore) TagByID(ctx context.Context, id uint64) (*Tag, error) {
	var tag *Tag
	err := store.db.GetContext(ctx, tag, "SELECT id, description FROM content.image_tags WHERE id = $1", id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return tag, nil
}

func (store *TagStore) TagByDescription(ctx context.Context, desc string) (*Tag, error) {
	var tag *Tag
	err := store.db.GetContext(ctx, tag, "SELECT id, description FROM content.image_tags WHERE description = $1", desc)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return tag, nil
}

func (store *TagStore) CreateTag(ctx context.Context, desc string) (*Tag, error) {
	tag := &Tag{
		Description: desc,
	}

	rows, err := store.db.QueryContext(ctx, `
	INSERT INTO content.image_tags(description)
	VALUES($1) RETURNING id;`,
		desc)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&tag.ID)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
	}

	return tag, nil
}

func (store *TagStore) ImagesForTag(ctx context.Context, id uint64) (*[]image.Image, error) {
	var imageIDs []uint64
	err := store.db.SelectContext(ctx, &imageIDs, "SELECT image_id FROM content.image_tag_bridge WHERE tag_id = $1", id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var images []image.Image
	for _, id := range imageIDs {
		image, err := store.image.ImageByID(ctx, id)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		images = append(images, *image)
	}
	return &images, nil
}

func (store *TagStore) TagsForImage(ctx context.Context, id uint64) (*[]Tag, error) {
	var tagIDs []uint64
	err := store.db.SelectContext(ctx, &tagIDs, "SELECT tag_id FROM content.image_tag_bridge WHERE image_id = $1", id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var tags []Tag
	for _, id := range tagIDs {
		tag, err := store.TagByID(ctx, id)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return &tags, nil
}

func (store *TagStore) TagImage(ctx context.Context, id uint64, imageID uint64) error {
	_, err := store.db.ExecContext(ctx, "INSERT INTO content.image_tag_bridge (tag_id, image_id) VALUES ($1, $2);", id, imageID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (store *TagStore) UnTagImage(ctx context.Context, id uint64, imageID uint64) error {
	_, err := store.db.ExecContext(ctx, "DELETE FROM content.image_tag_bridge WHERE tag_id = $1 AND image_id = $2;", id, imageID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
