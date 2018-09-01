package stream

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/fokal/fokal-core/pkg/services/image"
	"github.com/fokal/fokal-core/pkg/services/permission"
	"github.com/jmoiron/sqlx"
)

type PGStreamService struct {
	db          *sqlx.DB
	images      image.Service
	permissions permission.Service
}

func New(db *sqlx.DB, imageService image.Service, permissionService permission.Service) *PGStreamService {
	return &PGStreamService{db: db, images: imageService, permissions: permissionService}
}

func (stream *PGStreamService) StreamByID(ctx context.Context, id uint64) (*Stream, error) {
	var retrieved *Stream
	err := stream.db.GetContext(ctx, retrieved, "SELECT * FROM content.streams WHERE id = $1", id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return retrieved, nil
}

func (stream *PGStreamService) StreamsByCreator(ctx context.Context, userID uint64) (*[]Stream, error) {
	var streams *[]Stream
	err := stream.db.SelectContext(ctx, streams, "SELECT * FROM content.streams WHERE user_id = $1", userID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return streams, nil
}

func (stream *PGStreamService) CreateStream(ctx context.Context, creator uint64, title string) error {
	newStream := Stream{
		Creator:     creator,
		Title:       title,
		Description: nil,
		UpdatedAt:   time.Now(),
		CreatedAt:   time.Now(),
	}
	tx, err := stream.db.Beginx()
	if err != nil {
		logrus.Error(err)
		return err
	}

	rows, err := tx.QueryContext(ctx, `
	INSERT INTO content.streams(creator, title)
	VALUES($1, $2) RETURNING id;`,
		creator, title)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&newStream.ID)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}

	err = stream.permissions.AddScope(ctx, tx, creator, newStream.ID, permission.StreamClass, permission.CanEdit)
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = stream.permissions.Public(ctx, tx, newStream.ID, permission.StreamClass)
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = stream.permissions.AddScope(ctx, tx, creator, newStream.ID, permission.StreamClass, permission.CanDelete)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		if err := tx.Rollback(); err != nil {
			logrus.Error(err)
			return err
		}
		return err
	}
	return nil
}

func (stream *PGStreamService) SetDescription(ctx context.Context, id uint64, description string) error {
	_, err := stream.db.ExecContext(ctx, "UPDATE content.streams SET description = $2 WHERE id = $1", id, description)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
func (stream *PGStreamService) SetTitle(ctx context.Context, id uint64, title string) error {
	_, err := stream.db.ExecContext(ctx, "UPDATE content.streams SET title = $2 WHERE id = $1", id, title)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (stream *PGStreamService) AddImage(ctx context.Context, id, imageID uint64) error {
	_, err := stream.db.ExecContext(ctx, "INSERT INTO content.image_stream_bridge (stream_id, image_id) VALUES($1, $2)", id, imageID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (stream *PGStreamService) RemoveImage(ctx context.Context, id, imageID uint64) error {
	_, err := stream.db.ExecContext(ctx, "DELETE FROM content.image_stream_bridge WHERE stream_id = $1 AND image_id = $2", id, imageID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (stream *PGStreamService) Images(ctx context.Context, id uint64) (*[]image.Image, error) {
	var imageIDs *[]uint64
	err := stream.db.SelectContext(ctx, imageIDs, "SELECT image_id FROM content.image_stream_bridge WHERE id = $1", id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	images := []image.Image{}
	for _, id := range *imageIDs {
		img, err := stream.images.ImageByID(ctx, id)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		images = append(images, *img)
	}
	return &images, nil
}
