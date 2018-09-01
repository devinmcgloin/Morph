package permission

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Scope uint8

const (
	CanEdit Scope = iota
	CanDelete
	CanView
)

type ResourceClass uint8

const (
	UserClass ResourceClass = iota
	StreamClass
	ImageClass
)

type Permission struct {
	UserID        uint64
	ResourceID    uint64
	ResourceClass ResourceClass
	Scope         Scope
}

//go:generate moq -out permission_service_runner.go . PermissionService

type Service interface {
	ValidScope(ctx context.Context, userID, ResourceID uint64, class ResourceClass, scope Scope) (bool, error)

	Public(ctx context.Context, tx *sqlx.Tx, ResourceID uint64, class ResourceClass) error
	AddScope(ctx context.Context, tx *sqlx.Tx, userID, ResouceID uint64, class ResourceClass, scope Scope) error
	RemoveScope(ctx context.Context, tx *sqlx.Tx, userID, ResourceID uint64, class ResourceClass, scope Scope) error
}
