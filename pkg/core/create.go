package core

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"gopkg.in/mgo.v2/bson"

	"github.com/sprioc/sprioc-core/pkg/model"
	"github.com/sprioc/sprioc-core/pkg/refs"
	"github.com/sprioc/sprioc-core/pkg/rsp"
	"github.com/sprioc/sprioc-core/pkg/store"
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

	if !(validPassword(password) || validPassPhrase(password)) {
		return rsp.Response{Message: "Invalid password", Code: http.StatusBadRequest}
	}

	if store.ExistsUserID(username) || store.ExistsEmail(email) {
		return rsp.Response{Message: "Username or Email already exist", Code: http.StatusConflict}
	}

	securePassword, salt, resp := GetSaltPass(password)
	if !resp.Ok() {
		log.Println(resp)
		return rsp.Response{Message: "Error adding user", Code: http.StatusConflict}
	}

	usr := model.User{
		ID:        bson.NewObjectId(),
		Email:     email,
		Pass:      securePassword,
		Salt:      salt,
		ShortCode: username,
		AvatarURL: formatSources("default", "avatars"),
	}

	log.Printf("%+v", usr)

	err := store.Create("users", usr)
	if err != nil {
		return rsp.Response{Message: "Error adding user", Code: http.StatusConflict}
	}

	var response = make(map[string]string)
	response["link"] = refs.GetURL(refs.GetUserRef(username))

	return rsp.Response{Code: http.StatusAccepted, Data: response}
}

func CreateCollection(requestuser model.User, colData map[string]string) rsp.Response {
	var title, desc string
	var ok bool

	if title, ok = colData["title"]; !ok {
		return rsp.Response{Message: "Title not present", Code: http.StatusBadRequest}
	}

	if desc, ok = colData["desc"]; !ok {
		desc = ""
	}

	userRef := refs.GetUserRef(requestuser.ShortCode)
	colRef := refs.GetCollectionRef(store.GetNewCollectionShortCode())

	col := model.Collection{
		ID:        bson.NewObjectId(),
		Title:     title,
		Desc:      desc,
		Owner:     refs.GetUserRef(requestuser.ShortCode),
		ShortCode: colRef.Shortcode,
	}

	err := store.Create("collections", col)
	if err != nil {
		return rsp.Response{Code: http.StatusInternalServerError}
	}

	resp := Modify(userRef,
		bson.M{"$addToSet": bson.M{"collections": colRef}})
	if !resp.Ok() {
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
	letters := 0

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false
		}
	}

	return hasLower && hasUpper && hasNumber && hasSpecial && letters > 8
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
