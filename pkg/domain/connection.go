// Package domain : Connections between entities.
// Favorites, Followers, etc.
package domain

//go:generate moq -out connection_service_runner.go . ConnectionService
type ConnectionService interface {
	Favorite(userID, imageID uint64) error
	UnFavorite(userID, imageID uint64) error
	FavoritesForUser(userID uint8) (*[]Image, error)
	FavoritesForImage(imageID uint8) (*[]User, error)

	FollowUser(userID, followID uint64) error
	UnFollowUser(userID, followID uint64) error
	UserFollowers(userID uint64) (*[]User, error)
	UserFollowedBy(userID uint64) (*[]User, error)

	FollowStream(userID, followID uint64) error
	UnFollowStream(userID, followID uint64) error
	StreamFollowers(streamID uint64) (*[]User, error)
	StreamsFollowedBy(userID uint64) (*[]Stream, error)
}
