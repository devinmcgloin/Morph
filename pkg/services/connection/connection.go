// Package connection : Connections between entities.
// Favorites, Followers, etc.
package connection

import (
	"image"
	"os/user"

	"github.com/fokal/fokal-core/pkg/services/stream"
)

//go:generate moq -out connection_service_runner.go . ConnectionService
type Service interface {
	Favorite(userID, imageID uint64) error
	UnFavorite(userID, imageID uint64) error
	FavoritesForUser(userID uint8) (*[]image.Image, error)
	FavoritesForImage(imageID uint8) (*[]user.User, error)

	FollowUser(userID, followID uint64) error
	UnFollowUser(userID, followID uint64) error
	UserFollowers(userID uint64) (*[]user.User, error)
	UserFollowedBy(userID uint64) (*[]user.User, error)

	FollowStream(userID, followID uint64) error
	UnFollowStream(userID, followID uint64) error
	StreamFollowers(streamID uint64) (*[]user.User, error)
	StreamsFollowedBy(userID uint64) (*[]stream.Stream, error)
}
