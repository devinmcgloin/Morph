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

	"github.com/devinmcgloin/fokal/pkg/generator"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/rsp"
	"github.com/devinmcgloin/fokal/pkg/sql"
	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
)

var hmacSecret = []byte(os.Getenv("HMAC_SECRET"))
var dbase = os.Getenv("MONGODB_NAME")

const sessionLifetime = time.Minute * 10
const refreshAt = time.Minute * 1

func ValidateCredentialsByUserName(username string, password string) (bool, rsp.Response) {
	user, err := sql.GetLogin(username)
	if err != nil {
		return false, rsp.Response{Message: "Invalid Credentials.", Code: http.StatusUnauthorized}
	}
	return validUser(user, password)
}

func validUser(user map[string]interface{}, password string) (bool, rsp.Response) {
	salt, ok := user["salt"].(string)
	if !ok {
		return false, rsp.Response{Message: "Invalid Credentials.", Code: http.StatusUnauthorized}
	}

	truePass, ok := user["password"].(string)
	if !ok {
		return false, rsp.Response{Message: "Invalid Credentials.", Code: http.StatusUnauthorized}
	}

	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(truePass, shaString) == 0 {
		return true, rsp.Response{Code: http.StatusOK}
	}

	return false, rsp.Response{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
}

func generateSaltPass(password string) (string, string, rsp.Response) {
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

func CheckUser(r *http.Request) (model.Ref, rsp.Response) {
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.ExtractToken(r)

	if err != nil {
		return model.Ref{}, rsp.Response{Message: "Bearer Header not present", Code: http.StatusUnauthorized}
	}

	token := strings.Replace(tokenStrings, "Bearer ", "", 1)

	userRef, resp := VerifyJWT(token)
	if !resp.Ok() {
		log.Println("CheckUser")
		return model.Ref{}, resp
	}

	return userRef, rsp.Response{Code: http.StatusOK}
}

func CreateJWT(u model.Ref) (string, rsp.Response) {

	claims := &jwt.StandardClaims{
		//IssuedAt:  time.Now().Unix(),
		//ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:  "composer",
		Subject: u.Shortcode,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Println(err)
		return "", rsp.Response{Code: http.StatusInternalServerError, Message: "Unable to create token."}
	}

	return ss, rsp.Response{Code: http.StatusOK}
}

func VerifyJWT(tokenString string) (model.Ref, rsp.Response) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		log.Println(err)
		return model.Ref{}, rsp.Response{Message: err.Error(), Code: http.StatusBadRequest}
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			shortcode := claims["sub"].(string)
			id, err := sql.GetUserRef(shortcode)
			if err != nil {
				return model.Ref{}, rsp.Response{
					Code:    http.StatusBadRequest,
					Message: "Token is malformed"}
			}
			return id, rsp.Response{Code: http.StatusOK}
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return model.Ref{}, rsp.Response{Message: "Token is Malformed", Code: http.StatusBadRequest}
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return model.Ref{}, rsp.Response{Message: "Token is inactive", Code: http.StatusBadRequest}
		}
	}

	return model.Ref{}, rsp.Response{Message: "Token is invalid", Code: http.StatusBadRequest}
}
