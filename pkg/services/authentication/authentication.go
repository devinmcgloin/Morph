package authentication

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal-core/pkg/domain"
	"github.com/fokal/fokal-core/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type PGAuthService struct {
	db              *sqlx.DB
	userService     domain.UserService
	SessionLifetime time.Duration
	PrivateKey      *rsa.PrivateKey
	PublicKeys      map[string]*rsa.PublicKey
	KeyHash         string
}

func (auth *PGAuthService) CreateToken(ctx context.Context, userID uint64) (*string, error) {
	user, err := auth.userService.UserByID(ctx, userID)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	claims := &jwt.MapClaims{
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(auth.SessionLifetime).Unix(),
		"iss":   "fokal",
		"sub":   user.Username,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = auth.KeyHash
	ss, err := token.SignedString(auth.PrivateKey)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}
	return &ss, nil
}

func (auth *PGAuthService) VerifyToken(ctx context.Context, stringToken string) (bool, *uint64, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		interkid, ok := token.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("Missing kid header in token.\n")
		}

		kid, ok := interkid.(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("Invalid kid type.\n")
		}

		valid := false
		for k := range auth.PublicKeys {
			if k == kid {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("Invalid kid type.\n")

		}

		return auth.PublicKeys[kid], nil
	})
	if err != nil {
		logger.Error(ctx, err)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			email := claims["email"].(string)
			user, err := auth.userService.UserByEmail(ctx, email)
			if err != nil {
				return false, nil, errors.New("token is malformed")
			}
			return true, &user.ID, nil
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return false, nil, errors.New("token is malformed")
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return false, nil, errors.New("token is inactive")
		}
	}

	return false, nil, errors.New("token is invalid")
}

func (auth *PGAuthService) RefreshToken(ctx context.Context, stringToken string) (*string, error) {
	valid, id, err := auth.VerifyToken(ctx, stringToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid token provided")
	}
	return auth.CreateToken(ctx, *id)
}

func (auth *PGAuthService) PublicKey(ctx context.Context) (string, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(auth.PublicKeys[auth.KeyHash])
	if err != nil {
		logger.Error(ctx, err)
		return "", err
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	})

	return string(pemBytes), nil
}

func (auth *PGAuthService) LoadKeys(publicKey string) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	keys := make(map[string]string)
	err = json.Unmarshal(body, &keys)
	if err != nil {
		log.Fatal(err)
	}

	parsedKeys := make(map[string]*rsa.PublicKey)

	for kid, pem := range keys {
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			log.Fatal(err)
		}
		parsedKeys[kid] = publicKey
	}

	fokalPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		log.Fatal(err)
	}
	parsedKeys[auth.KeyHash] = fokalPublicKey
	auth.PublicKeys = parsedKeys

	privateStr := os.Getenv("PRIVATE_KEY")
	auth.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateStr))
	if err != nil {
		log.Fatal(err)
	}

}

func (auth *PGAuthService) RefreshGoogleOauthKeys() {
	tick := time.NewTicker(time.Minute * 10)
	go func() {
		for range tick.C {
			log.Println("Refreshing Google Auth Keys")
			resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
			if err != nil {
				log.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()

			keys := make(map[string]string)
			err = json.Unmarshal(body, &keys)
			if err != nil {
				log.Fatal(err)
			}
			for kid, pem := range keys {
				publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
				if err != nil {
					log.Fatal(err)
				}
				auth.PublicKeys[kid] = publicKey
			}
		}
	}()
}
