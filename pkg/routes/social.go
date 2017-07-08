package routes

import (
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/devinmcgloin/fokal/pkg/security/permissions"
	"github.com/devinmcgloin/fokal/pkg/social"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterSocialRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	put := api.Methods("PUT").Subrouter()
	del := api.Methods("DELETE").Subrouter()

	put.Handle("/i/{ID:[a-zA-Z]{12}}/favorite",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: social.FavoriteHandler}))

	del.Handle("/i/{ID:[a-zA-Z]{12}}/favorite",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanView,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: social.UnFavoriteHandler}))

	put.Handle("/u/{ID}/follow", chain.Append(
		handler.Middleware{State: state, M: security.Authenticate}.Handler,
		permissions.Middleware{State: state,
			T:          permissions.CanView,
			TargetType: model.Users,
			M:          permissions.PermissionMiddle}.Handler).
		Then(handler.Handler{
			State: state,
			H:     social.FollowHandler,
		}))
	del.Handle("/u/{ID}/follow", chain.Append(
		handler.Middleware{State: state, M: security.Authenticate}.Handler,
		permissions.Middleware{State: state,
			T:          permissions.CanView,
			TargetType: model.Users,
			M:          permissions.PermissionMiddle}.Handler).
		Then(handler.Handler{
			State: state,
			H:     social.UnFollowHandler,
		}))
}
