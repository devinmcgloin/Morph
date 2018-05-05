package security

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"
	"strings"

	"errors"

	"github.com/fokal/fokal-core/pkg/generator"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/request"
	"github.com/jmoiron/sqlx"
)

type Credentials struct {
	Id       int64  `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Salt     string `db:"salt"`
	Email    string `db:"email"`
}

func ValidateCredentials(db *sqlx.DB, req request.LoginRequest) (bool, error) {
	creds, err := GetLogin(db, req.Username)
	if err != nil {
		return false, handler.StatusError{Err: errors.New("Invalid Credentials."), Code: http.StatusUnauthorized}
	}

	hasher := sha512.New()

	passwordSalt := append([]byte(req.Password), []byte(creds.Salt)...)

	sha := hasher.Sum(passwordSalt)

	shaString := hex.EncodeToString(sha)

	if strings.Compare(creds.Password, shaString) == 0 {
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

// GetLogin returns the salt, password, email and username for a given user.
func GetLogin(db *sqlx.DB, username string) (*Credentials, error) {
	userInfo := new(Credentials)
	err := db.Get(userInfo, "SELECT id, username, salt, password, email FROM content.users WHERE username = $1 LIMIT 1;", username)
	if err != nil {
		log.Print(err)
		return userInfo, err
	}

	return userInfo, nil
}
