package tokens

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/retrieval"
)

func Create(state *handler.State, u model.Ref, email string) (string, error) {

	claims := &jwt.MapClaims{
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(state.SessionLifetime).Unix(),
		"iss":   "fokal",
		"sub":   u.Shortcode,
		"email": email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = state.KeyHash
	ss, err := token.SignedString(state.PrivateKey)
	if err != nil {
		log.Println(err)
		return "", handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("Unable to create token.")}
	}

	return ss, nil
}

func Parse(state *handler.State, r *http.Request) (*jwt.Token, error) {
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.ExtractToken(r)

	if err != nil {
		return nil, handler.StatusError{Err: errors.New("Bearer Header not present"), Code: http.StatusUnauthorized}
	}

	tokenStr := strings.Replace(tokenStrings, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		interkid, ok := token.Header["kid"]
		if !ok {
			return nil, fmt.Errorf("Missing kid header in token.\n")
		}

		kid, ok := interkid.(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("Invalid kid type.\n")
		}

		valid := false
		for k := range state.PublicKeys {
			if k == kid {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("Invalid kid type.\n")

		}

		return state.PublicKeys[kid], nil
	})
	return token, err
}

func Verify(state *handler.State, r *http.Request) (model.Ref, error) {

	token, err := Parse(state, r)

	if err != nil {
		return model.Ref{}, handler.StatusError{Err: err, Code: http.StatusBadRequest}
	}

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			email := claims["email"].(string)
			id, err := retrieval.GetUserRefByEmail(state.DB, email)
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
