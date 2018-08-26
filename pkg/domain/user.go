package domain

import (
	"context"
	"time"
)

type User struct {
	ID       uint64
	Username string
	Email    string

	Name      *string
	Bio       *string
	URL       *string
	Twitter   *string
	Instagram *string
	Location  *string

	AvatarID *string

	Featured bool
	Admin    bool

	CreatedAt    time.Time
	LastModified time.Time
}

//go:generate moq -out user_service_runner.go . UserService

type UserService interface {
	UserByID(ctx context.Context, id uint64) (*User, error)
	UserByUsername(ctx context.Context, username string) (*User, error)
	UserByEmail(ctx context.Context, username string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	Users(ctx context.Context, limit int) (*[]User, error)

	Admins(ctx context.Context) (*[]User, error)
	IsAdmin(ctx context.Context, id uint64) (bool, error)

	Featured(ctx context.Context) (*[]User, error)
	Feature(ctx context.Context, id uint64) error
	UnFeature(ctx context.Context, id uint64) error
	IsFeatured(ctx context.Context, id uint64) (bool, error)

	SetAvatarID(ctx context.Context, id uint64, avatarID string) error
	CreateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, id uint64) error
}
