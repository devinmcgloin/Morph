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
	chain = chain.Append(alice.Constructor(handler.Middleware{State: state, M: cache.Handler}.Handler))

	get.Handle("/i/featured",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.FeaturedImageHandler}))

	get.Handle("/i/recent",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.RecentImageHandler}))

	get.Handle("/i/color",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.ColorHandler}))

	get.Handle("/i/geo",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.GeoDistanceHandler}))

	get.Handle("/i/hot",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.HotImagesHander}))

	get.Handle("/i/text",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler).Then(handler.Handler{State: state, H: search.TextHandler}))

}
