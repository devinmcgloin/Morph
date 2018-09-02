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

	"github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal-core/pkg/log"
	"github.com/fokal/fokal-core/pkg/services/user"
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

type pgAuthService struct {
	db              *sqlx.DB
	userService     user.Service
	SessionLifetime time.Duration
	privateKey      *rsa.PrivateKey
	publicKeys      map[string]*rsa.PublicKey
}

func (auth *pgAuthService) CreateToken(ctx context.Context, userID uint64) (*string, error) {
	user, err := auth.userService.UserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return nil, err
	}
	return &ss, nil
}

func (auth *pgAuthService) ParseToken(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		interkid, ok := token.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("missing kid header in token")
		}

		kid, ok := interkid.(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("invalid kid type")
		}
		valid := false
		for k := range auth.publicKeys {
			if k == kid {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("invalid kid type")
		}

		return auth.publicKeys[kid], nil
	})
}

func (auth *pgAuthService) VerifyToken(ctx context.Context, stringToken string) (bool, *uint64, error) {
	token, err := auth.ParseToken(ctx, stringToken)
	if err != nil {
		log.WithContext(ctx).Error(err)
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			email := claims["email"].(string)
			user, err := auth.userService.UserByEmail(ctx, email)
			if err != nil {
				return false, nil, errors.New("valid token does not match fokal user")
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

func (auth *pgAuthService) RefreshToken(ctx context.Context, stringToken string) (*string, error) {
	valid, id, err := auth.VerifyToken(ctx, stringToken)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid token provided")
	}
	return auth.CreateToken(ctx, *id)
}

func (auth *pgAuthService) PublicKey(ctx context.Context) (string, error) {
	keyBytes, err := x509.MarshalPKIXPublicKey(auth.publicKeys[KeyHash])
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	})

	return string(pemBytes), nil
}

func (auth *pgAuthService) LoadKeys() {
	logrus.Println("Loading Google Auth Keys")

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		logrus.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	keys := make(map[string]string)
	err = json.Unmarshal(body, &keys)
	if err != nil {
		logrus.Fatal(err)
	}

	parsedKeys := make(map[string]*rsa.PublicKey)

	for kid, pem := range keys {
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			logrus.Fatal(err)
		}
		parsedKeys[kid] = publicKey
	}

	fokalPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if err != nil {
		logrus.Fatal(err)
	}
	parsedKeys[KeyHash] = fokalPublicKey
	auth.publicKeys = parsedKeys

	privateStr := os.Getenv("PRIVATE_KEY")
	auth.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(privateStr))
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println("Succesfully loaded Google Auth Keys")

}

func (auth *pgAuthService) RefreshGoogleOauthKeys() {
	tick := time.NewTicker(time.Minute * 10)
	go func() {
		for range tick.C {
			logrus.Println("Refreshing Google Auth Keys")
			resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
			if err != nil {
				logrus.Fatal(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logrus.Fatal(err)
			}
			resp.Body.Close()

			keys := make(map[string]string)
			err = json.Unmarshal(body, &keys)
			if err != nil {
				logrus.Fatal(err)
			}
			for kid, pem := range keys {
				publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
				if err != nil {
					logrus.Fatal(err)
				}
				auth.publicKeys[kid] = publicKey
			}
		}
	}()
}
