package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/devinmcgloin/sprioc/src/api/store"
	"github.com/devinmcgloin/sprioc/src/generator"
	"github.com/devinmcgloin/sprioc/src/model"
	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
)

var mongo = store.ConnectStore()

var hmacSecret = os.Getenv("HMAC_SECRET")
var dbase = os.Getenv("MONGODB_NAME")

var sessionLifetime = time.Minute * 10
var refreshAt = time.Minute * 1

func ValidateCredentialsByUserName(username string, password string) (bool, error) {
	user, err := mongo.GetByUserName(model.UserName(username))
	if err != nil {
		return false, errors.New("Invalid Credentials")
	}
	return validUser(user, password)
}

func validUser(user model.User, password string) (bool, error) {
	salt := user.Salt
	hasher := sha1.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(user.Pass, shaString) == 0 {
		return true, nil
	}
	return false, errors.New("Invalid Credentials")
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
	tokenString, err := jwtreq.HeaderExtractor{"Bearer"}.ExtractToken(r)
	if err != nil {
		return model.User{}, errors.New("Bearer Header not present")
	}

	userRef, err := VerifyJWT(tokenString)
	if err != nil {
		return model.User{}, errors.New("Malformed Header")
	}

	user, err := mongo.GetUser(userRef)
	if err != nil {
		return model.User{}, errors.New("Not Found")
	}

	return user, nil
}

func CreateJWT(u model.User) (string, error) {

	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:    "sprioc-core",
		Subject:   u.ID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	ss, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}
	return ss, nil
}

func VerifyJWT(tokenString string) (mgo.DBRef, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			id := claims["sub"]
			return mgo.DBRef{Database: dbase, Collection: "users", Id: id}, nil
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return mgo.DBRef{}, errors.New("Token is Malformed")
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return mgo.DBRef{}, errors.New("Token is inactive")
		}
	}

	return mgo.DBRef{}, errors.New("Token is Malformed")
}
