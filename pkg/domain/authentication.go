package domain

import (
	"context"
)

//go:generate moq -out authentication_service_runner.go . AuthenticationService

type AuthenticationService interface {
	CreateToken(ctx context.Context, userID uint64) (*string, error)
	VerifyToken(ctx context.Context, token string) (bool, *uint64, error)
	RefreshToken(ctx context.Context, token string) (*string, error)
	PublicKey(ctx context.Context) (string, error)
}
