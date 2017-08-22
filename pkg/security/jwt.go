package security

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"strings"

	"database/sql"

	"github.com/devinmcgloin/fokal/pkg/create"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/dgrijalva/jwt-go"
	jwtreq "github.com/dgrijalva/jwt-go/request"
)

const PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsW3uHvJvqaaMIW8wKP2E
NI3oVRghsNwUV4VN+5UH2oMAEaYaHiUfOvhXXRjPZo3q8f+v3rS4R7gfJXe8efP0
3x87DRB1uJlNNS777xDISnTLzVAOFFkLOTL9bOTJBlb69yCRhHV1NdUIPCGWntWC
WdKZBJ2zHOQUQgPpAn31imsYlvmlrLEoGNqKOPUQjwdtxEqEYpZyN84Hj5/NIhTC
F6rU8FhReQzEL27BHPfbUwTWUApmtfvCtrSc9pVM3MtlsMOf4OfoGg65kF5HJ/S8
tKRtL24z48ya+ntjbwbE3A5pEswm/Vm19wd77qbY5UILLmNf0xMQfwrkT/IcnBoD
pQIDAQAB
-----END PUBLIC KEY-----`

const keyHash = "554b5db484856bfa16e7da70a427dc4d9989678a"

func createJWT(state *handler.State, u model.Ref) (string, error) {

	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(state.SessionLifetime).Unix(),
		Issuer:    "fokal",
		Subject:   u.Shortcode,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = keyHash
	ss, err := token.SignedString(state.PrivateKey)
	if err != nil {
		log.Println(err)
		return "", handler.StatusError{Code: http.StatusInternalServerError, Err: errors.New("Unable to create token.")}
	}

	return ss, nil
}

func verifyJWT(state *handler.State, r *http.Request) (model.Ref, error) {
	var userRef model.Ref
	tokenStrings, err := jwtreq.HeaderExtractor{"Authorization"}.ExtractToken(r)

	log.Println(tokenStrings, r.Header.Get("Authorization"))

	if err != nil {
		return userRef, handler.StatusError{Err: errors.New("Bearer Header not present"), Code: http.StatusUnauthorized}
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
		for k, _ := range state.PublicKeys {
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

	if err != nil {
		return model.Ref{}, handler.StatusError{Err: err, Code: http.StatusBadRequest}
	}

	isGoogle := token.Header["kid"] != keyHash

	if token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if token.Valid && ok {
			if isGoogle {
				email := claims["email"].(string)
				id, err := retrieval.GetUserRefByEmail(state.DB, email)
				if err == sql.ErrNoRows {
					name := claims["name"].(string)
					var username string
					username = strings.Split(email, "@")[0]
					if domain, ok := claims["hd"]; ok {
						username = username + "." + domain.(string)

					}
					log.Printf("Createing new user: {Username: %s, Email: %s, Name: %s}", username, email, name)
					err = create.CommitUser(state.DB, username, email, name)
					if err != nil {
						return model.Ref{}, handler.StatusError{
							Code: http.StatusBadRequest,
							Err:  errors.New("Token is malformed")}
					} else {
						id, err = retrieval.GetUserRefByEmail(state.DB, email)
					}
				} else if err != nil {
					return model.Ref{}, handler.StatusError{
						Code: http.StatusBadRequest,
						Err:  errors.New("Token is malformed")}
				}
				return id, nil
			} else {
				shortcode := claims["sub"].(string)
				id, err := retrieval.GetUserRef(state.DB, shortcode)
				if err != nil {
					return model.Ref{}, handler.StatusError{
						Code: http.StatusBadRequest,
						Err:  errors.New("Token is malformed")}
				}
				return id, nil
			}
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
