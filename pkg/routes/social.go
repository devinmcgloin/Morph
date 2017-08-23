package routes

import (
	"github.com/fokal/fokal/pkg/handler"
	"github.com/fokal/fokal/pkg/model"
	"github.com/fokal/fokal/pkg/security"
	"github.com/fokal/fokal/pkg/security/permissions"
	"github.com/fokal/fokal/pkg/social"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterSocialRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	put := api.Methods("PUT").Subrouter()
	del := api.Methods("DELETE").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	put.Handle("/images/{ID:[a-zA-Z]{12}}/favorite",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: social.FavoriteHandler}))
	del.Handle("/images/{ID:[a-zA-Z]{12}}/favorite",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: social.UnFavoriteHandler}))
	opts.Handle("/images/{ID:[a-zA-Z]{12}}/favorite", chain.Then(handler.Options("PUT", "DELETE")))

	put.Handle("/users/{ID}/follow", chain.Append(
		handler.Middleware{State: state, M: security.Authenticate}.Handler,
		permissions.Middleware{State: state,
			T:          permissions.CanView,
			TargetType: model.Users,
			M:          permissions.PermissionMiddle}.Handler).
		Then(handler.Handler{
			State: state,
			H:     social.FollowHandler,
		}))
	del.Handle("/users/{ID}/follow", chain.Append(
		handler.Middleware{State: state, M: security.Authenticate}.Handler,
		permissions.Middleware{State: state,
			T:          permissions.CanView,
			TargetType: model.Users,
			M:          permissions.PermissionMiddle}.Handler).
		Then(handler.Handler{
			State: state,
			H:     social.UnFollowHandler,
		}))
	opts.Handle("/users/{ID}/follow", chain.Then(handler.Options("PUT", "DELETE")))

}
