package main

import (
	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/handlers"
	"github.com/sprioc/sprioc-core/pkg/middleware"
)

// TODO lock these routes down to alphabetical only with regex.
// TODO add names to linked routes.

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/images").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/{IID}", middleware.Unsecure(handlers.GetImage)).Name("image")
	get.HandleFunc("/{IID}/user", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/{IID}/collections", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/{IID}/album", middleware.Unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/images", middleware.Secure(handlers.UploadImage))

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{IID}/featured", middleware.Secure(handlers.FeatureImage))
	put.HandleFunc("/{IID}/favorite", middleware.Secure(handlers.FavoriteImage))

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{IID}", middleware.Secure(handlers.DeleteImage))
	del.HandleFunc("/{IID}/featured", middleware.Secure(handlers.UnFeatureImage))
	del.HandleFunc("/{IID}/favorite", middleware.Secure(handlers.UnFavoriteImage))

	patch := img.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{IID}", middleware.Secure(handlers.ModifyImage))
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/users").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/{username}", middleware.Unsecure(handlers.GetUser))
	get.HandleFunc("/{username}/location", middleware.Unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/users", middleware.Unsecure(handlers.CreateUser))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", middleware.Secure(handlers.AvatarUpload))
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
	get.HandleFunc("/{CID}/users", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/{CID}/images", middleware.Unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/collections", middleware.Secure(handlers.CreateCollection))

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", middleware.Secure(handlers.AddImageToCollection))
	put.HandleFunc("/{CID}/users", middleware.Secure(NotImplemented))
	put.HandleFunc("/{CID}/favorite", middleware.Secure(handlers.FavoriteCollection))
	put.HandleFunc("/{CID}/follow", middleware.Secure(handlers.FollowCollection))

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", middleware.Secure(handlers.DeleteCollection))
	del.HandleFunc("/{CID}/images/{IID}", middleware.Secure(handlers.DeleteImageFromCollection))
	del.HandleFunc("/{CID}/users/{username}", middleware.Secure(NotImplemented))
	del.HandleFunc("/{CID}/favorite", middleware.Secure(handlers.UnFavoriteCollection))
	del.HandleFunc("/{CID}/follow", middleware.Secure(handlers.UnFollowCollection))

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", middleware.Secure(NotImplemented))
}

func registerAlbumRoutes(api *mux.Router) {
	alb := api.PathPrefix("/albums").Subrouter()

	get := alb.Methods("GET").Subrouter()
	get.HandleFunc("/{AID}", middleware.Unsecure(handlers.GetAlbum))
	get.HandleFunc("/{AID}/images", middleware.Unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/albums", middleware.Secure(NotImplemented))

	put := alb.Methods("PUT").Subrouter()
	put.HandleFunc("/{AID}/images", middleware.Secure(handlers.AddImageToAlbum))
	put.HandleFunc("/{AID}/favorite", middleware.Secure(handlers.FavoriteAlbum))
	put.HandleFunc("/{AID}/follow", middleware.Secure(handlers.FollowAlbum))

	del := alb.Methods("DELETE").Subrouter()
	del.HandleFunc("/{AID}", middleware.Secure(handlers.DeleteAlbum))
	del.HandleFunc("/{AID}/images/{IID}", middleware.Secure(handlers.DeleteImageFromAlbum))
	del.HandleFunc("/{AID}/favorite", middleware.Secure(handlers.UnFavoriteAlbum))
	del.HandleFunc("/{AID}/follow", middleware.Secure(handlers.UnFollowAlbum))

	patch := alb.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{AID}", middleware.Secure(NotImplemented))
}

func registerSearchRoutes(api *mux.Router) {
	get := api.Methods("GET").Subrouter()

	get.HandleFunc("/images", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/uers", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/collections", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/albums", middleware.Unsecure(NotImplemented))
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
