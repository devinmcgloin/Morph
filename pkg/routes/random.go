package routes

import (
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/random"
	"github.com/fokal/fokal-core/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterRandomRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	get := api.Methods("GET").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	get.Handle("/images/random",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).
			Then(handler.Handler{State: state, H: random.ImageHandler}))
	opts.Handle("/images/random", chain.Then(handler.Options("GET")))
}
