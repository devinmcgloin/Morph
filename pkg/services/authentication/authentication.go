package authentication

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal-core/pkg/services/user"
	"github.com/jmoiron/sqlx"
)

//go:generate moq -out authentication_service_runner.go . AuthenticationService

type Service interface {
	CreateToken(ctx context.Context, userID uint64) (*string, error)
	VerifyToken(ctx context.Context, token string) (bool, *uint64, error)
	RefreshToken(ctx context.Context, token string) (*string, error)
	PublicKey(ctx context.Context) (string, error)
	ParseToken(ctx context.Context, token string) (*jwt.Token, error)
}

func New(db *sqlx.DB, userService user.Service, sessionLifetime time.Duration) Service {
	service := &pgAuthService{
		db:              db,
		userService:     userService,
		SessionLifetime: sessionLifetime,
	}
	service.LoadKeys()
	service.RefreshGoogleOauthKeys()
	return service
}
