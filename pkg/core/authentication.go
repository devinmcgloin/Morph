package core

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/sprioc/conductor/pkg/generator"
	"github.com/sprioc/conductor/pkg/model"
	"github.com/sprioc/conductor/pkg/rsp"
)

var hmacSecret = []byte(os.Getenv("HMAC_SECRET"))
var dbase = os.Getenv("MONGODB_NAME")

var sessionLifetime = time.Minute * 10
var refreshAt = time.Minute * 1

func ValidateCredentialsByUserName(username string, password string) (bool, model.User, rsp.Response) {
	user, resp := GetUser(model.DBRef{Database: dbase, Collection: "users", Shortcode: username})
	if !resp.Ok() {
		return false, model.User{}, rsp.Response{Message: "Invalid Credentials."}
	}
	return validUser(user, password)
}

func validUser(user model.User, password string) (bool, model.User, rsp.Response) {
	salt := user.Salt
	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(user.Pass, shaString) == 0 {
		return true, user, rsp.Response{Code: http.StatusOK}
	}
	return false, model.User{}, rsp.Response{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
}

func GetSaltPass(password string) (string, string, rsp.Response) {
	salt, err := generator.GenerateSecureString(64)
	if err != nil {
		return "", "", rsp.Response{Message: "Unable to create user", Code: http.StatusInternalServerError}
	}
	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	saltedPass := hex.EncodeToString(sha)
	return saltedPass, salt, rsp.Response{Code: http.StatusOK}
}

func CheckUser(r *http.Request) (model.User, rsp.Response) {
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.
		ExtractToken(r)

	if err != nil {
		return model.User{}, rsp.Response{Message: "Bearer Header not present", Code: http.StatusUnauthorized}
	}

	token := strings.Replace(tokenStrings, "Bearer ", "", 1)

	userRef, resp := VerifyJWT(token)
	if !resp.Ok() {
		log.Println("CheckUser")
		return model.User{}, resp
	}

	user, resp := GetUser(userRef)
	if !resp.Ok() {
		return model.User{}, resp
	}

	return user, rsp.Response{Code: http.StatusOK}
}

func CreateJWT(u model.User) (string, rsp.Response) {

	claims := &jwt.StandardClaims{
		//IssuedAt:  time.Now().Unix(),
		//ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:  "conductor",
		Subject: u.ShortCode,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Println(err)
		return "", rsp.Response{Code: http.StatusInternalServerError, Message: "Unable to create token."}
	}

	return ss, rsp.Response{Code: http.StatusOK}
}

func VerifyJWT(tokenString string) (model.DBRef, rsp.Response) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		log.Println(err)
		return model.DBRef{}, rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			id := claims["sub"].(string)
			return model.DBRef{Database: dbase, Collection: "users", Shortcode: id},
				rsp.Response{Code: http.StatusOK}
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return model.DBRef{}, rsp.Response{Message: "Token is Malformed", Code: http.StatusBadRequest}
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return model.DBRef{}, rsp.Response{Message: "Token is inactive", Code: http.StatusBadRequest}
		}
	}

	return model.DBRef{}, rsp.Response{Message: "Token is invalid", Code: http.StatusBadRequest}
}

// RefreshToken updates a token with a new expiration time. After 3 days it expires.
// TODO need to implement refreshtoken
func RefreshToken(tokenString string) string {
	return ""
}
