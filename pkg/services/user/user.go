package user

import (
	"context"
	"time"

	"github.com/fokal/fokal-core/pkg/request"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	Name      *string `json:"name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	URL       *string `json:"url,omitempty"`
	Twitter   *string `json:"twitter,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
	Location  *string `json:"location,omitempty"`

	AvatarID *string `db:"avatar_id" json:"avatar_id,omitempty"`

	Featured bool `json:"featured"`
	Admin    bool `json:"admin"`

	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
}

//go:generate moq -out user_service_runner.go . UserService

type Service interface {
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
	PatchUser(ctx context.Context, u uint64, changes request.PatchUser) error
	DeleteUser(ctx context.Context, id uint64) error
}
