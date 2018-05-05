package routes

import (
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/model"
	"github.com/fokal/fokal-core/pkg/modification"
	"github.com/fokal/fokal-core/pkg/security"
	"github.com/fokal/fokal-core/pkg/security/permissions"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterModificationRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {
	put := api.Methods("PUT").Subrouter()
	del := api.Methods("DELETE").Subrouter()
	opts := api.Methods("OPTIONS").Subrouter()
	patch := api.Methods("PATCH").Subrouter()

	//Image Routes
	put.Handle("/images/{ID:[a-zA-Z]{12}}/featured",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.FeatureImage}))
	del.Handle("/images/{ID:[a-zA-Z]{12}}/featured",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.UnFeatureImage}))

	opts.Handle("/images/{ID:[a-zA-Z]{12}}/featured", chain.Then(handler.Options("DELETE", "PUT")))

	del.Handle("/images/{ID:[a-zA-Z]{12}}",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanDelete,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.DeleteImage}))

	patch.Handle("/images/{ID:[a-zA-Z]{12}}",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler,
			permissions.Middleware{State: state,
				T:          permissions.CanEdit,
				TargetType: model.Images,
				M:          permissions.PermissionMiddle}.Handler).
			Then(handler.Handler{State: state, H: modification.PatchImage}))

	put.Handle("/images/{ID:[a-zA-Z]{12}}/download",
		chain.
			Then(handler.Handler{State: state, H: modification.DownloadHandler}))

	opts.Handle("/images/{ID:[a-zA-Z]{12}}", chain.Then(handler.Options("PATCH", "DELETE")))
	opts.Handle("/images/{ID:[a-zA-Z]{12}}/download", chain.Then(handler.Options("PUT")))

	// User Routes
	del.Handle("/users/me",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler).
			Then(handler.Handler{State: state, H: modification.DeleteUser}))
	patch.Handle("/users/me",
		chain.Append(
			handler.Middleware{State: state, M: security.Authenticate}.Handler).
			Then(handler.Handler{State: state, H: modification.PatchUser}))
	opts.Handle("/users/me", chain.Then(handler.Options("PATCH", "DELETE")))

}
