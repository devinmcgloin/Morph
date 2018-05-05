package security

import (
	"net/http"

	"crypto/x509"

	"encoding/pem"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/tokens"
)

func PublicKeyHandler(state *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	key := make(map[string]string)
	keyBytes, err := x509.MarshalPKIXPublicKey(state.PublicKeys[state.KeyHash])
	if err != nil {
		return handler.Response{}, handler.StatusError{Code: http.StatusInternalServerError}
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: keyBytes,
	})

	key[state.KeyHash] = string(pemBytes)
	return handler.Response{Code: http.StatusOK, Data: key}, nil
}

func RefreshHandler(state *handler.State, w http.ResponseWriter, r *http.Request) (handler.Response, error) {
	user, err := tokens.Verify(state, r)
	if err != nil {
		return handler.Response{}, err

	}
	oldJWT, _ := tokens.Parse(state, r)
	claims := oldJWT.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	jwt, err := tokens.Create(state, user, email)
	if err != nil {
		return handler.Response{}, err
	}
	return handler.Response{Code: http.StatusOK, Data: map[string]string{"token": jwt}}, nil
}
