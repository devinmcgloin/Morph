package permission

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
)

type PGPermission struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PGPermission {
	return &PGPermission{
		db: db,
	}
}

func (pgp *PGPermission) AddScope(ctx context.Context, tx *sqlx.Tx, userID, resourceID uint64, class ResourceClass, scope Scope) error {
	var query string
	switch scope {
	case CanEdit:
		query = "INSERT INTO permissions.can_edit (user_id, o_id, class) VALUES ($1, $2, $3)"
	case CanDelete:
		query = "INSERT INTO permissions.can_delete (user_id, o_id, class) VALUES ($1, $2, $3)"
	case CanView:
		query = "INSERT INTO permissions.can_view (user_id, o_id, class) VALUES ($1, $2, $3)"
	}
	_, err := tx.ExecContext(ctx, query, userID, resourceID, class)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (pgp *PGPermission) ValidScope(ctx context.Context, userID, resourceID uint64, class ResourceClass, scope Scope) (bool, error) {
	var query string
	switch scope {
	case CanEdit:
		query = "SELECT count(1) FROM permissions.can_edit WHERE (user_id = $1 OR user_id = -1) AND o_id = $2 AND class = $3"
	case CanDelete:
		query = "SELECT count(1) FROM permissions.can_delete WHERE (user_id = $1 OR user_id = -1) AND o_id = $2 AND class = $3"
	case CanView:
		query = "SELECT count(1) FROM permissions.can_view WHERE (user_id = $1 OR user_id = -1)AND o_id = $2 AND class = $3"
	}
	_, err := pgp.db.ExecContext(ctx, query, userID, resourceID, class)
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	return false, nil
}

func (pgp *PGPermission) RemoveScope(ctx context.Context, tx *sqlx.Tx, userID, resourceID uint64, class ResourceClass, scope Scope) error {
	var query string
	switch scope {
	case CanEdit:
		query = "DELETE FROM permissions.can_edit WHERE user_id = $1 AND o_id = $2 AND class = $3"
	case CanDelete:
		query = "DELETE FROM permissions.can_delete WHERE user_id = $1 AND o_id = $2 AND class = $3"
	case CanView:
		query = "DELETE FROM permissions.can_view WHERE user_id = $1 AND o_id = $2 AND class = $3"
	}
	_, err := tx.ExecContext(ctx, query, userID, resourceID, class)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (pgp *PGPermission) Public(ctx context.Context, tx *sqlx.Tx, resourceID uint64, class ResourceClass) error {
	_, err := tx.ExecContext(ctx, "INSERT INTO permissions.can_view (user_id, o_id, class) VALUES ($1, $2, $3)", -1, resourceID, class)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
