package permissionservice

import (
	"context"

	"github.com/fokal/fokal-core/pkg/logger"

	"github.com/fokal/fokal-core/pkg/domain"
	"github.com/jmoiron/sqlx"
)

type PGPermission struct {
	db *sqlx.DB
}

func (pgp *PGPermission) AddScope(ctx context.Context, userID, resourceID uint64, class domain.ResourceClass, scope domain.Scope) error {
	var query string
	switch scope {
	case domain.CanEdit:
		query = "INSERT INTO permissions.can_edit (user_id, o_id, type) VALUES ($1, $2, $3)"
	case domain.CanDelete:
		query = "INSERT INTO permissions.can_delete (user_id, o_id, type) VALUES ($1, $2, $3)"
	case domain.CanView:
		query = "INSERT INTO permissions.can_view (user_id, o_id, type) VALUES ($1, $2, $3)"
	}
	_, err := pgp.db.ExecContext(ctx, query, userID, resourceID, class)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	return nil
}

func (pgp *PGPermission) ValidScope(ctx context.Context, userID, resourceID uint64, class domain.ResourceClass, scope domain.Scope) (bool, error) {
	var query string
	switch scope {
	case domain.CanEdit:
		query = "SELECT count(1) FROM permissions.can_edit WHERE (user_id = $1 OR user_id = -1) AND o_id = $2 AND type = $3"
	case domain.CanDelete:
		query = "SELECT count(1) FROM permissions.can_delete WHERE (user_id = $1 OR user_id = -1) AND o_id = $2 AND type = $3"
	case domain.CanView:
		query = "SELECT count(1) FROM permissions.can_view WHERE (user_id = $1 OR user_id = -1)AND o_id = $2 AND type = $3"
	}
	_, err := pgp.db.ExecContext(ctx, query, userID, resourceID, class)
	if err != nil {
		logger.Error(ctx, err)
		return false, err
	}
	return false, nil
}

func (pgp *PGPermission) RemoveScope(ctx context.Context, userID, resourceID uint64, class domain.ResourceClass, scope domain.Scope) error {
	var query string
	switch scope {
	case domain.CanEdit:
		query = "DELETE FROM permissions.can_edit WHERE user_id = $1 AND o_id = $2 AND type = $3"
	case domain.CanDelete:
		query = "DELETE FROM permissions.can_delete WHERE user_id = $1 AND o_id = $2 AND type = $3"
	case domain.CanView:
		query = "DELETE FROM permissions.can_view WHERE user_id = $1 AND o_id = $2 AND type = $3"
	}
	_, err := pgp.db.ExecContext(ctx, query, userID, resourceID, class)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	return nil
}

func (pgp *PGPermission) Public(ctx context.Context, resourceID uint64, class domain.ResourceClass) error {
	_, err := pgp.db.ExecContext(ctx, "INSERT INTO permissions.can_view (user_id, o_id, type) VALUES ($1, $2, $3)", -1, resourceID, class)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	return nil
}
