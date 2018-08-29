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
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal-core/pkg/domain"
	"github.com/jmoiron/sqlx"
)

const PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsW3uHvJvqaaMIW8wKP2E
NI3oVRghsNwUV4VN+5UH2oMAEaYaHiUfOvhXXRjPZo3q8f+v3rS4R7gfJXe8efP0
3x87DRB1uJlNNS777xDISnTLzVAOFFkLOTL9bOTJBlb69yCRhHV1NdUIPCGWntWC
WdKZBJ2zHOQUQgPpAn31imsYlvmlrLEoGNqKOPUQjwdtxEqEYpZyN84Hj5/NIhTC
F6rU8FhReQzEL27BHPfbUwTWUApmtfvCtrSc9pVM3MtlsMOf4OfoGg65kF5HJ/S8
tKRtL24z48ya+ntjbwbE3A5pEswm/Vm19wd77qbY5UILLmNf0xMQfwrkT/IcnBoD
pQIDAQAB
-----END PUBLIC KEY-----`

const KeyHash = "554b5db484856bfa16e7da70a427dc4d9989678a"

type PGAuthService struct {
	db              *sqlx.DB
	userService     domain.UserService
	SessionLifetime time.Duration
	privateKey      *rsa.PrivateKey
	publicKeys      map[string]*rsa.PublicKey
}

func New(db *sqlx.DB, userService domain.UserService, sessionLifetime time.Duration) *PGAuthService {
	service := &PGAuthService{
		db:              db,
		userService:     userService,
		SessionLifetime: sessionLifetime,
	}
	service.LoadKeys()
	service.RefreshGoogleOauthKeys()
	return service
}

func (auth *PGAuthService) CreateToken(ctx context.Context, userID uint64) (*string, error) {
	user, err := auth.userService.UserByID(ctx, userID)
	if err != nil {
		log.Error(err)
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
	token.Header["kid"] = KeyHash
	ss, err := token.SignedString(auth.privateKey)
	if err != nil {
		log.Error(err)
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
		for k := range auth.publicKeys {
			if k == kid {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("Invalid kid type.\n")

		}

		return auth.publicKeys[kid], nil
	})
	if err != nil {
		log.Error(err)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			email := claims["email"].(string)
			user, err := auth.userService.UserByEmail(ctx, email)
			if err != nil {
				return true, nil, errors.New("valid token does not match fokal user")
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
	keyBytes, err := x509.MarshalPKIXPublicKey(auth.publicKeys[KeyHash])
	if err != nil {
		log.Error(err)
		return "", err
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	})

	return string(pemBytes), nil
}

func (auth *PGAuthService) LoadKeys() {
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

	fokalPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if err != nil {
		log.Fatal(err)
	}
	parsedKeys[KeyHash] = fokalPublicKey
	auth.publicKeys = parsedKeys

	privateStr := os.Getenv("PRIVATE_KEY")
	auth.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateStr))
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
				auth.publicKeys[kid] = publicKey
			}
		}
	}()
}
