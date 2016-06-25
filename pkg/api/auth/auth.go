package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"encoding/pem"

	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/sprioc/sprioc-core/pkg/api/store"
	"github.com/sprioc/sprioc-core/pkg/env"
	"github.com/sprioc/sprioc-core/pkg/generator"
	"github.com/sprioc/sprioc-core/pkg/model"
)

var mongo = store.ConnectStore()

var privKey *rsa.PrivateKey
var pubKey *rsa.PublicKey
var dbase = env.Getenv("MONGODB_NAME", "morph")

var sessionLifetime = time.Minute * 10
var refreshAt = time.Minute * 1

func init() {
	var err error

	privKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = privKey.Validate()
	if err != nil {
		fmt.Println("Validation failed.", err)
	}
	pubKey = &privKey.PublicKey
}

func ValidateCredentialsByUserName(username string, password string) (bool, model.User, error) {
	user, err := mongo.GetByUserName(username)
	if err != nil {
		return false, model.User{}, errors.New("Invalid Credentials")
	}
	return validUser(user, password)
}

func validUser(user model.User, password string) (bool, model.User, error) {
	salt := user.Salt
	hasher := sha1.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(user.Pass, shaString) == 0 {
		return true, user, nil
	}
	return false, model.User{}, errors.New("Invalid Credentials")
}

func GetSaltPass(password string) (string, string, error) {
	salt, err := generator.GenerateSecureString(64)
	if err != nil {
		return "", "", errors.New("Unable to create user")
	}
	hasher := sha1.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	saltedPass := hex.EncodeToString(sha)
	return saltedPass, salt, nil
}

func CheckUser(r *http.Request) (model.User, error) {
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.
		ExtractToken(r)

	if err != nil {
		return model.User{}, errors.New("Bearer Header not present")
	}

	token := strings.Split(tokenStrings, " ")
	if len(tokenStrings) < 2 {
		return model.User{}, errors.New("Check to make sure your header follows 'Bearer <token>'")
	}

	userRef, err := VerifyJWT(token[1])
	if err != nil {
		return model.User{}, err
	}
	log.Println(userRef)

	user, err := mongo.GetUser(userRef)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func CreateJWT(u model.User) (string, error) {

	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:    "sprioc-core",
		Subject:   u.ShortCode,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(privKey)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return ss, nil
}

func VerifyJWT(tokenString string) (model.DBRef, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
	})

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			id := claims["sub"].(string)
			return model.DBRef{Database: dbase, Collection: "users", Shortcode: id}, nil
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return model.DBRef{}, errors.New("Token is Malformed")
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return model.DBRef{}, errors.New("Token is inactive")
		}
	}

	return model.DBRef{}, errors.New("Token is Invalid")
}

func GetPublicKey() []byte {

	pub, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		fmt.Println("Failed to get der format for PublicKey.", err)
		return []byte{}
	}

	pubBLK := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pub,
	}
	return pem.EncodeToMemory(&pubBLK)
}
