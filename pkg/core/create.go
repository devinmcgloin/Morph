package core

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/composer/pkg/model"
	"github.com/sprioc/composer/pkg/mongo"
	"github.com/sprioc/composer/pkg/redis"
	"github.com/sprioc/composer/pkg/refs"
	"github.com/sprioc/composer/pkg/rsp"
)

var validEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Anything but special characters and spaces.
var validUsername = regexp.MustCompile(`^[^\<\>\!\{\}\[\]\!\@\#\$\%\^\&\*\(\)\.\ ]{3,16}$`)

var letters = regexp.MustCompile("^[a-zA-Z]+$")

func CreateUser(userData map[string]string) rsp.Response {
	var username, email, password string
	var ok bool

	if username, ok = userData["username"]; !ok {
		return rsp.Response{Message: "Username not present", Code: http.StatusBadRequest}
	}
	username = strings.ToLower(username)

	if !validUsername.MatchString(username) {
		return rsp.Response{Message: "Invalid username", Code: http.StatusBadRequest}
	}

	if email, ok = userData["email"]; !ok {
		return rsp.Response{Message: "Email not present", Code: http.StatusBadRequest}
	}
	email = strings.ToLower(email)

	if !validEmail.MatchString(email) {
		return rsp.Response{Message: "Invalid email", Code: http.StatusBadRequest}
	}

	if password, ok = userData["password"]; !ok {
		return rsp.Response{Message: "Password not present", Code: http.StatusBadRequest}
	}

	password = strings.TrimSpace(password)

	if !(validPassword(password) || validPassPhrase(password)) {
		return rsp.Response{Message: "Invalid password", Code: http.StatusBadRequest}
	}

	userRef := model.Ref{Collection: model.Users, ShortCode: username}

	exists, err := redis.Exists(userRef)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if exists {
		return rsp.Response{Message: "Username already exist", Code: http.StatusConflict}
	}

	exists, err = redis.ExistsEmail(email)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	if exists {
		return rsp.Response{Message: "Email already exist", Code: http.StatusConflict}
	}

	securePassword, salt, resp := generateSaltPass(password)
	if !resp.Ok() {
		log.Println(resp)
		return rsp.Response{Message: "Error adding user", Code: http.StatusConflict}
	}

	usr := model.User{
		ID:        bson.NewObjectId(),
		Email:     email,
		AvatarURL: formatSources("default", "avatars"),
	}

	err = redis.CreateUser(userRef, usr.ID, email, securePassword, salt)
	if err != nil {
		return rsp.Response{Message: "Error adding user", Code: http.StatusInternalServerError}
	}

	err = mongo.Create("users", usr)
	if err != nil {
		return rsp.Response{Message: "Error adding user", Code: http.StatusConflict}
	}

	var response = make(map[string]string)
	response["link"] = refs.GetURL(refs.GetUserRef(username))

	return rsp.Response{Code: http.StatusAccepted, Data: response}
}

func CreateCollection(requestuser model.Ref, colData map[string]string) rsp.Response {
	var title, desc string
	var ok bool

	if title, ok = colData["title"]; !ok {
		return rsp.Response{Message: "Title not present", Code: http.StatusBadRequest}
	}

	if desc, ok = colData["desc"]; !ok {
		desc = ""
	}

	colRef, err := redis.GenerateShortCode(model.Collections)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}
	col := model.Collection{
		ID:    bson.NewObjectId(),
		Title: title,
		Desc:  desc,
	}

	err = redis.CreateCollection(requestuser, colRef, col.ID)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	err = mongo.Create("collections", col)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	var response = make(map[string]string)
	response["link"] = refs.GetURL(colRef)

	return rsp.Response{Code: http.StatusAccepted, Data: response}
}

func validPassword(password string) bool {
	var hasUpper bool
	var hasLower bool
	var hasSpecial bool
	var hasNumber bool

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		default:
			return false
		}
	}

	log.Println(hasLower, hasUpper, hasNumber, hasSpecial, len(password))
	return hasLower && hasUpper && hasNumber && hasSpecial && len(password) > 8
}

func validPassPhrase(passphrase string) bool {
	sections := strings.Split(passphrase, "-")

	for _, sect := range sections {
		if !letters.MatchString(sect) {
			return false
		} else if len(sect) < 5 {
			return false
		}
	}
	return len(sections) >= 3
}
