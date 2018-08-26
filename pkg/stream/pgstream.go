package stream

import (
	"context"
	"time"

	"github.com/fokal/fokal-core/pkg/domain"
	"github.com/fokal/fokal-core/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type PGStreamService struct {
	db *sqlx.DB
}

func (stream *PGStreamService) StreamByID(ctx context.Context, id uint64) (*domain.Stream, error) {
	var stream *domain.Stream
	err := stream.db.GetContext(ctx, "SELECT * FROM content.streams WHERE id = $1", id)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return stream, nil
}

func (stream *PGStreamService) StreamsByCreator(ctx context.Context, userID uint64) (*[]Stream, error) {
	var streams *[]domain.Stream
	err := stream.db.SelectContext(ctx, "SELECT * FROM content.streams WHERE user_id = $1", userID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return stream, nil
}

func (stream *PGStreamService) CreateStream(ctx context.Context, creator uint64, title string) error {
	newStream := domain.Stream{
		Creator:     creator,
		Title:       title,
		Description: nil,
		UpdatedAt:   time.Now,
		CreatedAt:   time.Now,
	}
	tx, err := store.db.Beginx()
	if err != nil {
		logger.Error(err)
		return err
	}

	rows, err := tx.QueryContext(ctx, `
	INSERT INTO content.streams(creator, title)
	VALUES($1, $2) RETURNING id;`,
		creator, title)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&newStream.ID)
		if err != nil {
			logger.Error(ctx, err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(ctx, err)
		if err := tx.Rollback(); err != nil {
			logger.Error(ctx, err)
			return err
		}
		return err
	}

	err = store.permissions.AddScope(userID, newStream.ID, domain.StreamClass, domain.CanEdit)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	err = store.permissions.Public(userID, newStream.ID, domain.StreamClass)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	err = store.permissions.AddScope(userID, newStream.ID, domain.StreamClass, domain.CanDelete)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}

	return nil
}

func (stream *pgstreamservice) SetDescription(ctx context.context, id uint64, description string) error {
	err := stream.db.ExecContext(ctx, "UPDATE content.streams SET description = $2 WHERE id = $1", id, description)
	if err != nil {
		logger.error(ctx, err)
		return err
	}

	return nil
}
func (stream *pgstreamservice) SetTitle(ctx context.context, id uint64, title string) error {
	err := stream.db.ExecContext(ctx, "UPDATE content.streams SET title = $2 WHERE id = $1", id, title)
	if err != nil {
		logger.error(ctx, err)
		return err
	}

	return nil
}
