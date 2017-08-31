package routes

import (
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterAuthRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	opts := api.Methods("OPTIONS").Subrouter()

	get := api.Methods("GET").Subrouter()
	get.Handle("/auth/certs", chain.Then(handler.Handler{State: state, H: security.PublicKeyHandler}))
	opts.Handle("/auth/certs", chain.Then(handler.Options("GET")))

	get.Handle("/auth/refresh", chain.Then(handler.Handler{State: state, H: security.RefreshHandler}))
	opts.Handle("/auth/refresh", chain.Then(handler.Options("GET")))

}
