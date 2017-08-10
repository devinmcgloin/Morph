package routes

import (
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/modification"
	"github.com/devinmcgloin/fokal/pkg/security"
	"github.com/devinmcgloin/fokal/pkg/security/permissions"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterModificationRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	put := api.Methods("PUT").Subrouter()
	del := api.Methods("DELETE").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()

	//Image Routes
	put.Handle("/i/{ID:[a-zA-Z]{12}}/featured",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.FeatureImage}))
	del.Handle("/i/{ID:[a-zA-Z]{12}}/featured",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.UnFeatureImage}))

	opts.Handle("/i/{ID:[a-zA-Z]{12}}/featured", chain.Then(handler.Options("DELETE", "PUT")))

	del.Handle("/i/{ID:[a-zA-Z]{12}}",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanDelete,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.DeleteImage}))

	patch := api.Methods("PATCH").Subrouter()
	patch.Handle("/i/{ID:[a-zA-Z]{12}}",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.PatchImage}))
	opts.Handle("/i/{ID:[a-zA-Z]{12}}", chain.Then(handler.Options("PATCH", "DELETE")))

	// User Routes
	del.Handle("/u/{ID}",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanDelete,
				TargetType: model.Users,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.DeleteUser}))
	patch.Handle("/u/{ID}",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Users,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.PatchUser}))
	opts.Handle("/u/{ID}", chain.Then(handler.Options("PATCH", "DELETE")))

}
