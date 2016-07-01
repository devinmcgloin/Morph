package main

import (
	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/handlers"
	"github.com/sprioc/sprioc-core/pkg/middleware"
)

// TODO lock these routes down to alphabetical only with regex.

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/images").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/{IID}", middleware.Unsecure(handlers.GetImage))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/images", middleware.Secure(handlers.UploadImage))

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{IID}/tags", middleware.Secure(NotImplemented))
	put.HandleFunc("/{IID}/featured", middleware.Secure(handlers.FeatureImage))
	put.HandleFunc("/{IID}/favorite", middleware.Secure(handlers.FavoriteImage))

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{IID}", middleware.Secure(handlers.DeleteImage))
	del.HandleFunc("/{IID}/tags", middleware.Secure(NotImplemented))
	del.HandleFunc("/{IID}/featured", middleware.Secure(handlers.UnFeatureImage))
	del.HandleFunc("/{IID}/favorite", middleware.Secure(handlers.UnFavoriteImage))

	patch := img.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{IID}", middleware.Secure(handlers.ModifyImage))
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/users").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/{username}", middleware.Unsecure(handlers.GetUser))
	get.HandleFunc("/{username}/images", middleware.Unsecure(handlers.GetUserImages))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/users", middleware.Unsecure(handlers.CreateUser))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", middleware.Secure(handlers.UploadAvatar))
	put.HandleFunc("/{username}/favorite", middleware.Secure(handlers.FavoriteUser))
	put.HandleFunc("/{username}/follow", middleware.Secure(handlers.FollowUser))

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", middleware.Secure(handlers.DeleteUser))
	del.HandleFunc("/{username}/favorite", middleware.Secure(handlers.UnFavoriteUser))
	del.HandleFunc("/{username}/follow", middleware.Secure(handlers.UnFollowUser))

	patch := usr.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{username}", middleware.Secure(handlers.ModifyUser))
}

func registerCollectionRoutes(api *mux.Router) {
	col := api.PathPrefix("/collections").Subrouter()

	get := col.Methods("GET").Subrouter()
	get.HandleFunc("/{CID}", middleware.Unsecure(handlers.GetCollection))
	get.HandleFunc("/{CID}/images", middleware.Unsecure(handlers.GetCollectionImages))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/collections", middleware.Secure(handlers.CreateCollection))

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", middleware.Secure(handlers.AddImageToCollection))
	put.HandleFunc("/{CID}/favorite", middleware.Secure(handlers.FavoriteCollection))
	put.HandleFunc("/{CID}/follow", middleware.Secure(handlers.FollowCollection))

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", middleware.Secure(handlers.DeleteCollection))
	del.HandleFunc("/{CID}/images/{IID}", middleware.Secure(handlers.DeleteImageFromCollection))
	del.HandleFunc("/{CID}/favorite", middleware.Secure(handlers.UnFavoriteCollection))
	del.HandleFunc("/{CID}/follow", middleware.Secure(handlers.UnFollowCollection))

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", middleware.Secure(handlers.ModifyCollection))
}

func registerSearchRoutes(api *mux.Router) {
	get := api.Methods("GET").Subrouter()

	get.HandleFunc("/stream", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/search", middleware.Unsecure(NotImplemented))

}

// routes that return random results for a given collection.
// TODO redirect to new thing or just return random one like normal.
func registerLuckyRoutes(api *mux.Router) {

}

func registerAuthRoutes(api *mux.Router) {
	post := api.Methods("POST").Subrouter()

	post.HandleFunc("/get_token", middleware.Unsecure(handlers.GetToken))

}
