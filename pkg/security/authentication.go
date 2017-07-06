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

	"errors"

	"github.com/devinmcgloin/fokal/pkg/generator"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/jmoiron/sqlx"
)

var hmacSecret = []byte(os.Getenv("HMAC_SECRET"))
var dbase = os.Getenv("MONGODB_NAME")

const sessionLifetime = time.Minute * 10
const refreshAt = time.Minute * 1

func GetToken(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (map[string]string, error) {
	decoder := json.NewDecoder(r.Body)

	var creds = make(map[string]string)
	var username, password string
	var ok bool

	err := decoder.Decode(&creds)
	if err != nil {
		return nil, handler.StatusError{Err: errors.New("Bad Request"), Code: http.StatusBadRequest}
	}

	if username, ok = creds["username"]; !ok {
		return nil, handler.StatusError{Err: errors.New("Bad Request"), Code: http.StatusBadRequest}
	}

	if password, ok = creds["password"]; !ok {
		return nil, handler.StatusError{Err: errors.New("Bad Request"), Code: http.StatusBadRequest}
	}

	valid, err := ValidateCredentialsByUserName(db, username, password)
	if err != nil {
		return nil, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
	}

	if valid {
		token, err := CreateJWT(model.Ref{Collection: model.Users, Shortcode: username})
		if err != nil {
			return nil, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
		}
		return map[string]string{"token": token}, nil
	}

	return nil, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
}

type tok struct {
	Token string `json:"token"`
}

func ValidateCredentialsByUserName(db *sqlx.DB, username string, password string) (bool, error) {
	user, err := GetLogin(db, username)
	if err != nil {
		return false, handler.StatusError{Err: errors.New("Invalid Credentials."), Code: http.StatusUnauthorized}
	}
	return validUser(user, password)
}

func validUser(user map[string]interface{}, password string) (bool, error) {
	salt, ok := user["salt"].(string)
	if !ok {
		return false, handler.StatusError{Err: errors.New("Invalid Credentials."), Code: http.StatusUnauthorized}
	}

	truePass, ok := user["password"].(string)
	if !ok {
		return false, handler.StatusError{Err: errors.New("Invalid Credentials."), Code: http.StatusUnauthorized}
	}

	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(truePass, shaString) == 0 {
		return true, nil
	}

	return false, handler.StatusError{Err: errors.New("Invalid Credentials"), Code: http.StatusUnauthorized}
}

func GenerateSaltPass(password string) (string, string, error) {
	salt, err := generator.GenerateSecureString(64)
	if err != nil {
		return "", "", handler.StatusError{Err: errors.New("Unable to create user"), Code: http.StatusInternalServerError}
	}
	hasher := sha512.New()

	passwordSalt := append([]byte(password), []byte(salt)...)

	sha := hasher.Sum(passwordSalt)

	saltedPass := hex.EncodeToString(sha)
	return saltedPass, salt, nil
}

func CheckUser(db *sqlx.DB, r *http.Request) (model.Ref, error) {
	var userRef model.Ref
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.ExtractToken(r)

	if err != nil {
		return userRef, handler.StatusError{Err: errors.New("Bearer Header not present"), Code: http.StatusUnauthorized}
	}

	token := strings.Replace(tokenStrings, "Bearer ", "", 1)

	userRef, err = VerifyJWT(db, token)
	if err != nil {
		log.Println("CheckUser")
		return model.Ref{}, err
	}

	return userRef, nil
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
		return "", handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("Unable to create token.")}
	}

	return ss, nil
}

func VerifyJWT(db *sqlx.DB, tokenString string) (model.Ref, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return model.Ref{}, handler.StatusError{Err: err, Code: http.StatusBadRequest}
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			shortcode := claims["sub"].(string)
			id, err := retrieval.GetUserRef(db, shortcode)
			if err != nil {
				return model.Ref{}, handler.StatusError{
					Code: http.StatusBadRequest,
					Err:  errors.New("Token is malformed")}
			}
			return id, nil
		}
	} else if err, ok := err.(*jwt.ValidationError); ok {
		if err.Errors&jwt.ValidationErrorMalformed != 0 {
			// Token is malformed
			return model.Ref{}, handler.StatusError{Err: errors.New("Token is Malformed"), Code: http.StatusBadRequest}
		} else if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return model.Ref{}, handler.StatusError{Err: errors.New("Token is inactive"), Code: http.StatusBadRequest}
		}
	}

	return model.Ref{}, handler.StatusError{Err: errors.New("Token is invalid"), Code: http.StatusBadRequest}
}

// GetLogin returns the salt, password, email and username for a given user.
func GetLogin(db *sqlx.DB, ref string) (map[string]interface{}, error) {
	userInfo := make(map[string]interface{})
	rows, err := db.Query("SELECT id, username, salt, password, email FROM content.users WHERE username = $1 LIMIT 1;", ref)
	if err != nil {
		log.Print(err)
		return userInfo, err
	}
	defer rows.Close()
	var id int64
	var username string
	var salt string
	var password string
	var email string
	for rows.Next() {
		if err := rows.Scan(&id, &username, &salt, &password, &email); err != nil {
			log.Print(err)
			return userInfo, err
		}
	}
	if err := rows.Err(); err != nil {
		log.Print(err)
		return userInfo, err
	}
	userInfo["id"] = id
	userInfo["username"] = username
	userInfo["salt"] = salt
	userInfo["password"] = password
	userInfo["email"] = email

	return userInfo, nil
}
