package routes

import (
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

	get.Handle("/i/{ID:[a-zA-Z]{12}}",
		chain.Append(
			handler.Middleware{
				State: state,
				M:     security.SetAuthenticatedUser,
			}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle,
			}.Handler).Then(handler.Handler{State: state, H: retrieval.ImageHandler}))

	get.Handle("/u/me", chain.Append(
		handler.Middleware{
			State: state,
			M:     security.Authenticate,
		}.Handler).Then(handler.Handler{State: state, H: retrieval.LoggedInUserHandler}))

	get.Handle("/u/{ID}", chain.Then(handler.Handler{State: state, H: retrieval.UserHandler}))

	get.Handle("/t/{ID}", chain.Then(handler.Handler{State: state, H: retrieval.TagHandler}))

	//TODO user images leaks hidden images.
	get.Handle("/u/{ID}/images", chain.Then(handler.Handler{State: state, H: retrieval.UserImagesHandler}))

}
