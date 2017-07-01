package security

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"encoding/json"

	"github.com/devinmcgloin/fokal/pkg/ferr"
	"github.com/devinmcgloin/fokal/pkg/generator"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/sql"
	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
)

var hmacSecret = []byte(os.Getenv("HMAC_SECRET"))
var dbase = os.Getenv("MONGODB_NAME")

const sessionLifetime = time.Minute * 10
const refreshAt = time.Minute * 1

func GetToken(w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r.Body)

	var creds = make(map[string]string)
	var username, password string
	var ok bool

	err := decoder.Decode(&creds)
	if err != nil {
		return nil, ferr.FError{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if username, ok = creds["username"]; !ok {
		return ferr.FError{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	if password, ok = creds["password"]; !ok {
		return ferr.FError{Message: "Bad Request", Code: http.StatusBadRequest}
	}

	valid, resp := core.ValidateCredentialsByUserName(username, password)
	if !resp.Ok() {
		return ferr.FError{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
	}

	if valid {
		token, resp := core.CreateJWT(model.Ref{Collection: model.Users, Shortcode: username})
		if !resp.Ok() {
			log.Println(resp)
			return resp
		}
		return ferr.FError{Code: 201, Data: tok{Token: token}}
	}

	return ferr.FError{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
}

type tok struct {
	Token string `json:"token"`
}

func ValidateCredentialsByUserName(username string, password string) (bool, error) {
	user, err := sql.GetLogin(username)
	if err != nil {
		return false, ferr.FError{Message: "Invalid Credentials.", Code: http.StatusUnauthorized}
	}
	return validUser(user, password)
}

func validUser(user map[string]interface{}, password string) (bool, error) {
	salt, ok := user["salt"].(string)
	if !ok {
		return false, ferr.FError{Message: "Invalid Credentials.", Code: http.StatusUnauthorized}
	}

	truePass, ok := user["password"].(string)
	if !ok {
		return false, ferr.FError{Message: "Invalid Credentials.", Code: http.StatusUnauthorized}
	}

	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(truePass, shaString) == 0 {
		return true, ferr.FError{Code: http.StatusOK}
	}

	return false, ferr.FError{Message: "Invalid Credentials", Code: http.StatusUnauthorized}
}

func GenerateSaltPass(password string) (string, string, error) {
	salt, err := generator.GenerateSecureString(64)
	if err != nil {
		return "", "", ferr.FError{Message: "Unable to create user", Code: http.StatusInternalServerError}
	}
	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	saltedPass := hex.EncodeToString(sha)
	return saltedPass, salt, ferr.FError{Code: http.StatusOK}
}

func CheckUser(r *http.Request) (model.Ref, error) {
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.ExtractToken(r)

	if err != nil {
		return model.Ref{}, ferr.FError{Message: "Bearer Header not present", Code: http.StatusUnauthorized}
	}

	token := strings.Replace(tokenStrings, "Bearer ", "", 1)

	userRef, resp := VerifyJWT(token)
	if !resp.Ok() {
		log.Println("CheckUser")
		return model.Ref{}, resp
	}

	return userRef, ferr.FError{Code: http.StatusOK}
}

func CreateJWT(u model.Ref) (string, error) {

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
		return "", ferr.FError{Code: http.StatusInternalServerError, Message: "Unable to create token."}
	}

	return ss, ferr.FError{Code: http.StatusOK}
}

func VerifyJWT(tokenString string) (model.Ref, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		log.Println(err)
		return model.Ref{}, ferr.FError{Message: err.Error(), Code: http.StatusBadRequest}
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			shortcode := claims["sub"].(string)
			id, err := sql.GetUserRef(shortcode)
			if err != nil {
				return model.Ref{}, ferr.FError{
					Code:    http.StatusBadRequest,
					Message: "Token is malformed"}
			}
			return id, ferr.FError{Code: http.StatusOK}
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return model.Ref{}, ferr.FError{Message: "Token is Malformed", Code: http.StatusBadRequest}
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return model.Ref{}, ferr.FError{Message: "Token is inactive", Code: http.StatusBadRequest}
		}
	}

	return model.Ref{}, ferr.FError{Message: "Token is invalid", Code: http.StatusBadRequest}
}
