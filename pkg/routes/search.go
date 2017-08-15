package routes

import (
	"github.com/devinmcgloin/fokal/pkg/cache"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/search"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterSearchRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	get := api.Methods("GET").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()
	cache := chain.Append(alice.Constructor(handler.Middleware{State: state, M: cache.Handler}.Handler))

	get.Handle("/i/featured",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.FeaturedImageHandler}))
	opts.Handle("/i/featured", chain.Then(handler.Options("GET")))

	get.Handle("/i/recent",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.RecentImageHandler}))
	opts.Handle("/i/recent", chain.Then(handler.Options("GET")))

	get.Handle("/i/hot",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.HotImagesHander}))
	opts.Handle("/i/hot", chain.Then(handler.Options("GET")))

	get.Handle("/i/search",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.SearchHandler}))
	opts.Handle("/i/search", chain.Then(handler.Options("GET")))

}
