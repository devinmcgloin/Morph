package routes

import (
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/search"
	"github.com/fokal/fokal/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterSearchRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	post := api.Methods("POST").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	post.Handle("/search",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.SearchHandler}))
	opts.Handle("/search", chain.Then(handler.Options("POST")))

}
