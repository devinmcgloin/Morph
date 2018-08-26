package domain

import jwt "github.com/dgrijalva/jwt-go"

//go:generate moq -out authentication_service_runner.go . AuthenticationService

type AuthenticationService interface {
	CreateToken(userID uint64) (*jwt.Token, error)
	VerifyToken(token *jwt.Token) (bool, error)
	RefreshToken(token *jwt.Token) (*jwt.Token, error)
	PublicKey() (string, error)
}
