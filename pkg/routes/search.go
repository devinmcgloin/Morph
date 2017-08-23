package routes

import (
	"github.com/fokal/fokal/pkg/cache"
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/search"
	"github.com/fokal/fokal/pkg/security"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterSearchRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	get := api.Methods("GET").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()
	cache := chain.Append(alice.Constructor(handler.Middleware{State: state, M: cache.Handler}.Handler))

	get.Handle("/images/featured",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.FeaturedImageHandler}))
	opts.Handle("/images/featured", chain.Then(handler.Options("GET")))

	get.Handle("/images/recent",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.RecentImageHandler}))
	opts.Handle("/images/recent", chain.Then(handler.Options("GET")))

	get.Handle("/images/trending",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.TrendingImagesHander}))
	opts.Handle("/images/trending", chain.Then(handler.Options("GET")))

	get.Handle("/images/search",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.SearchHandler}))
	opts.Handle("/images/search", chain.Then(handler.Options("GET")))

}
