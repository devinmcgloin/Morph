package domain

import "context"

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

type PermissionService interface {
	Public(ctx context.Context, ResourceID uint64, class ResourceClass) error
	AddScope(ctx context.Context, userID, ResouceID uint64, class ResourceClass, scope Scope) error
	ValidScope(ctx context.Context, userID, ResourceID uint64, class ResourceClass, scope Scope) (bool, error)
	RemoveScope(ctx context.Context, userID, ResourceID uint64, class ResourceClass, scope Scope) error
}
