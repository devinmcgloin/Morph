package routes

import (
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterAuthRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	post := api.Methods("POST").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	post.Handle("/auth/token", chain.Then(handler.Handler{State: state, H: security.LoginHandler}))
	opts.Handle("/auth/token", chain.Then(handler.Options("POST")))

	get := api.Methods("GET").Subrouter()
	get.Handle("/auth/certs", chain.Then(handler.Handler{State: state, H: security.PublicKeyHandler}))
	opts.Handle("/auth/certs", chain.Then(handler.Options("GET")))

	get.Handle("/auth/refresh", chain.Then(handler.Handler{State: state, H: security.RefreshHandler}))
	opts.Handle("/auth/refresh", chain.Then(handler.Options("GET")))

}
