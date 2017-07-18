package routes

import (
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterAuthRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	post := api.Methods("POST").Subrouter()
	post.Handle("/auth/token", chain.Then(handler.Handler{State: state, H: security.LoginHandler}))

	get := api.Methods("GET").Subrouter()
	get.Handle("/auth/certs", chain.Then(handler.Handler{State: state, H: security.PublicKeyHandler}))
}
