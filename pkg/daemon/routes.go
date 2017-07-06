package daemon

import (
	"net/http"

	"github.com/devinmcgloin/fokal/pkg/create"
	"github.com/devinmcgloin/fokal/pkg/handler"
	"github.com/devinmcgloin/fokal/pkg/model"
	"github.com/devinmcgloin/fokal/pkg/modification"
	"github.com/devinmcgloin/fokal/pkg/retrieval"
	"github.com/devinmcgloin/fokal/pkg/security/permissions"
	"github.com/devinmcgloin/fokal/pkg/social"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func registerImageRoutes(api *mux.Router, chain alice.Chain) {
	img := api.PathPrefix("/i").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.Handle("/{ID:[a-zA-Z]{12}}", permission(chain, permissions.CanView, model.Images).Then(handler.Handler{State: &AppState, H: retrieval.ImageHandler}))

	post := api.Methods("POST").Subrouter()
	post.Handle("/i", auth(chain).Then(handler.Handler{State: &AppState, H: create.ImageHandler}))

	put := img.Methods("PUT").Subrouter()
	put.Handle("/{ID:[a-zA-Z]{12}}/featured",
		permission(chain, permissions.CanEdit, model.Images).Then(handler.Handler{State: &AppState, H: modification.FeatureImage}))
	put.Handle("/{ID:[a-zA-Z]{12}}/favorite",
		permission(chain, permissions.CanView, model.Images).Then(handler.Handler{State: &AppState, H: social.FavoriteHandler}))

	del := img.Methods("DELETE").Subrouter()
	del.Handle("/{ID:[a-zA-Z]{12}}",
		permission(chain, permissions.CanDelete, model.Images).Then(handler.Handler{State: &AppState, H: modification.DeleteImage}))

	del.Handle("/{ID:[a-zA-Z]{12}}/featured",
		permission(chain, permissions.CanEdit, model.Images).Then(handler.Handler{State: &AppState, H: modification.UnFeatureImage}))

	del.Handle("/{ID:[a-zA-Z]{12}}/favorite",
		permission(chain, permissions.CanView, model.Images).Then(handler.Handler{State: &AppState, H: social.UnFavoriteHandler}))

	patch := img.Methods("PATCH").Subrouter()
	patch.Handle("/{ID:[a-zA-Z]{12}}",
		permission(chain, permissions.CanEdit, model.Images).Then(handler.Handler{State: &AppState, H: modification.PatchImage}))

}

func registerUserRoutes(api *mux.Router, chain alice.Chain) {
	usr := api.PathPrefix("/u").Subrouter()

	get := usr.Methods("GET").Subrouter()
	//get.Handle("/me", chain.Then(handler.Handler{State: &AppState, H: retrieval.UserHandler}))
	get.Handle("/{ID}", chain.Then(handler.Handler{State: &AppState, H: retrieval.UserHandler}))

	post := api.Methods("POST").Subrouter()
	post.Handle("/u", auth(chain).Then(handler.Handler{State: &AppState, H: create.UserHandler}))

	put := usr.Methods("PUT").Subrouter()
	put.Handle("/{ID}/avatar", permission(chain, permissions.CanEdit, model.Users).Then(handler.Handler{
		State: &AppState,
		H:     create.AvatarHandler,
	}))
	put.Handle("/{ID}/follow", permission(chain, permissions.CanView, model.Users).Then(handler.Handler{
		State: &AppState,
		H:     social.FollowHandler,
	}))
	del := usr.Methods("DELETE").Subrouter()
	del.Handle("/{ID}", permission(chain, permissions.CanDelete, model.Users).Then(handler.Handler{
		State: &AppState,
		H:     modification.DeleteUser,
	}))
	del.Handle("/{ID}/follow",
		permission(chain, permissions.CanView, model.Users).Then(handler.Handler{State: &AppState, H: social.UnFollowHandler}))

	patch := usr.Methods("PATCH").Subrouter()
	patch.Handle("/{ID}",
		permission(chain, permissions.CanEdit, model.Users).Then(handler.Handler{State: &AppState, H: modification.PatchUser}))

}

//
// func registerCollectionRoutes(api *mux.Router, chain http.Handler) {
// 	col := api.PathPrefix("/collections").Subrouter()
//
// 	get := col.Methods("GET").Subrouter()
// 	get.HandleFunc("/{CID:[a-zA-Z]{12}}", middleware.Unsecure(handlers.GetCollection))
// 	get.HandleFunc("/{CID:[a-zA-Z]{12}}/images", middleware.Unsecure(handlers.GetCollectionImages))
//
// 	post := api.Methods("POST").Subrouter()
// 	post.HandleFunc("/collections", middleware.Secure(handlers.CreateCollection))
//
// 	put := col.Methods("PUT").Subrouter()
// 	put.HandleFunc("/{CID:[a-zA-Z]{12}}/images", middleware.Secure(handlers.AddImageToCollection))
// 	put.HandleFunc("/{CID:[a-zA-Z]{12}}/favorite", middleware.Secure(handlers.FavoriteCollection))
// 	put.HandleFunc("/{CID:[a-zA-Z]{12}}/follow", middleware.Secure(handlers.FollowCollection))
//
// 	del := col.Methods("DELETE").Subrouter()
// 	del.HandleFunc("/{CID:[a-zA-Z]{12}}", middleware.Secure(handlers.DeleteCollection))
// 	del.HandleFunc("/{CID:[a-zA-Z]{12}}/images", middleware.Secure(handlers.DeleteImageFromCollection))
// 	del.HandleFunc("/{CID:[a-zA-Z]{12}}/favorite", middleware.Secure(handlers.UnFavoriteCollection))
// 	del.HandleFunc("/{CID:[a-zA-Z]{12}}/follow", middleware.Secure(handlers.UnFollowCollection))
//
// 	patch := col.Methods("PATCH").Subrouter()
// 	patch.HandleFunc("/{CID:[a-zA-Z]{12}}", middleware.Secure(handlers.ModifyCollection))
// }
//
// func registerSearchRoutes(api *mux.Router, chain http.Handler) {
// 	get := api.Methods("GET").Subrouter()
//
// 	get.HandleFunc("/stream", middleware.Secure(handlers.GetStream))
//
// 	post := api.Methods("POST").Subrouter()
// 	post.HandleFunc("/search", middleware.Unsecure(handlers.Search))
//
// }

// routes that return random results for a given collection.
// TODO redirect to new thing or just return random one like normal.
func registerLuckyRoutes(api *mux.Router, chain http.Handler) {

}

//
//func registerAuthRoutes(api *mux.Router, chain http.Handler) {
//	post := api.Methods("POST").Subrouter()
//
//	post.HandleFunc("/get_token", middleware.Unsecure(handlers.GetToken))
//
//}
