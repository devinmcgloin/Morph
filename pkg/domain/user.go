package domain

import (
	"time"
)

type User struct {
	ID       uint64
	Username string
	Email    string

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
	UserByID(id uint64) (*User, error)
	UserByUsername(username string) (*User, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
	Users() ([]*User, error)
	Admins() ([]*User, error)

	Featured() ([]*User, error)
	Feature(id uint64, user uint64) error
	UnFeature(id uint64, user uint64) error

	SetAvatarID(id uint64, avatarID string) error
	IsAdmin(id uint64) (bool, error)
	CreateUser(u *User) error
	DeleteUser(id uint64) error
}
