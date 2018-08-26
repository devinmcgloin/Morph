package domain

import (
	"context"

	jwt "github.com/dgrijalva/jwt-go"
)

//go:generate moq -out authentication_service_runner.go . AuthenticationService

type AuthenticationService interface {
	CreateToken(ctx context.Context, userID uint64) (*jwt.Token, error)
	VerifyToken(ctx context.Context, token *jwt.Token) (bool, error)
	RefreshToken(ctx context.Context, token *jwt.Token) (*jwt.Token, error)
	PublicKey(ctx context.Context) (string, error)
}
