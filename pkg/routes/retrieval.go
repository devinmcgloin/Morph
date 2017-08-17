package routes

import (
	"github.com/devinmcgloin/fokal/pkg/cache"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/devinmcgloin/fokal/pkg/security/permissions"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterRetrievalRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	get := api.Methods("GET").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	cache := chain.Append(alice.Constructor(handler.Middleware{State: state, M: cache.Handler}.Handler))
	get.Handle("/i/{ID:[a-zA-Z]{12}}",
		cache.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle,
			}.Handler).Then(handler.Handler{State: state, H: retrieval.ImageHandler}))
	opts.Handle("/i/{ID:[a-zA-Z]{12}}", chain.Then(handler.Options("GET")))

	get.Handle("/u/me", chain.Append(
		handler.Middleware{
			State: state,
			M:     security.Authenticate,
		}.Handler).Then(handler.Handler{State: state, H: retrieval.LoggedInUserHandler}))
	opts.Handle("/u/me", chain.Then(handler.Options("GET")))

	get.Handle("/u/me/images", chain.Append(
		handler.Middleware{
			State: state,
			M:     security.Authenticate,
		}.Handler).Then(handler.Handler{State: state, H: retrieval.LoggedInUserImagesHandler}))
	opts.Handle("/u/me/images", chain.Then(handler.Options("GET")))

	get.Handle("/u/{ID}", cache.Then(handler.Handler{State: state, H: retrieval.UserHandler}))
	opts.Handle("/u/{ID}", chain.Then(handler.Options("GET")))

	get.Handle("/t/{ID}", cache.Then(handler.Handler{State: state, H: retrieval.TagHandler}))
	opts.Handle("/t/{ID}", chain.Then(handler.Options("GET")))

	get.Handle("/u/{ID}/images", cache.Then(handler.Handler{State: state, H: retrieval.UserImagesHandler}))
	opts.Handle("/u/{ID}/images", chain.Then(handler.Options("GET")))
}
