package main

import (
	"github.com/gorilla/mux"
	"github.com/sprioc/sprioc-core/pkg/handlers"
)

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/images").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/{IID}", unsecure(handlers.GetImage))
	get.HandleFunc("/{IID}/user", unsecure(NotImplemented))
	get.HandleFunc("/{IID}/collections", unsecure(NotImplemented))
	get.HandleFunc("/{IID}/album", unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/images", secure(handlers.UploadImage))

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{IID}/featured", secure(handlers.FeatureImage))
	put.HandleFunc("/{IID}/favorite", secure(handlers.FavoriteImage))

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{IID}", secure(handlers.DeleteImage))
	del.HandleFunc("/{IID}/featured", secure(handlers.UnFeatureImage))
	del.HandleFunc("/{IID}/favorite", secure(handlers.UnFavoriteImage))

	patch := img.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{IID}", secure(handlers.ModifyImage))
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/users").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/{username}", unsecure(handlers.GetUser))
	get.HandleFunc("/{username}/location", unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/users", unsecure(handlers.CreateUser))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", secure(handlers.AvatarUpload))
	put.HandleFunc("/{username}/favorite", secure(handlers.FavoriteUser))
	put.HandleFunc("/{username}/follow", secure(handlers.FollowUser))

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", secure(handlers.DeleteUser))
	del.HandleFunc("/{username}/favorite", secure(handlers.UnFavoriteUser))
	del.HandleFunc("/{username}/follow", secure(handlers.UnFollowUser))

	patch := usr.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{username}", secure(handlers.ModifyUser))
}

func registerCollectionRoutes(api *mux.Router) {
	col := api.PathPrefix("/collections").Subrouter()

	get := col.Methods("GET").Subrouter()
	get.HandleFunc("/{CID}", unsecure(handlers.GetCollection))
	get.HandleFunc("/{CID}/users", unsecure(NotImplemented))
	get.HandleFunc("/{CID}/images", unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/collections", secure(handlers.CreateCollection))

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", secure(handlers.AddImageToCollection))
	put.HandleFunc("/{CID}/users", secure(NotImplemented))
	put.HandleFunc("/{CID}/favorite", secure(handlers.FavoriteCollection))
	put.HandleFunc("/{CID}/follow", secure(handlers.FollowCollection))

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", secure(handlers.DeleteCollection))
	del.HandleFunc("/{CID}/images/{IID}", secure(handlers.DeleteImageFromCollection))
	del.HandleFunc("/{CID}/users/{username}", secure(NotImplemented))
	del.HandleFunc("/{CID}/favorite", secure(handlers.UnFavoriteCollection))
	del.HandleFunc("/{CID}/follow", secure(handlers.UnFollowCollection))

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", secure(NotImplemented))
}

func registerAlbumRoutes(api *mux.Router) {
	alb := api.PathPrefix("/albums").Subrouter()

	get := alb.Methods("GET").Subrouter()
	get.HandleFunc("/{AID}", unsecure(handlers.GetAlbum))
	get.HandleFunc("/{AID}/images", unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/albums", secure(NotImplemented))

	put := alb.Methods("PUT").Subrouter()
	put.HandleFunc("/{AID}/images", secure(handlers.AddImageToAlbum))
	put.HandleFunc("/{AID}/favorite", secure(handlers.FavoriteAlbum))
	put.HandleFunc("/{AID}/follow", secure(handlers.FollowAlbum))

	del := alb.Methods("DELETE").Subrouter()
	del.HandleFunc("/{AID}", secure(handlers.DeleteAlbum))
	del.HandleFunc("/{AID}/images/{IID}", secure(handlers.DeleteImageFromAlbum))
	del.HandleFunc("/{AID}/favorite", secure(handlers.UnFavoriteAlbum))
	del.HandleFunc("/{AID}/follow", secure(handlers.UnFollowAlbum))

	patch := alb.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{AID}", secure(NotImplemented))
}

func registerSearchRoutes(api *mux.Router) {
	get := api.Methods("GET").Subrouter()

	get.HandleFunc("/images", unsecure(NotImplemented))
	get.HandleFunc("/uers", unsecure(NotImplemented))
	get.HandleFunc("/collections", unsecure(NotImplemented))
	get.HandleFunc("/albums", unsecure(NotImplemented))
	get.HandleFunc("/search", unsecure(NotImplemented))

}

// routes that return random results for a given collection.
// TODO redirect to new thing or just return random one like normal.
func registerLuckyRoutes(api *mux.Router) {

}

func registerAuthRoutes(router *mux.Router) {
	get := router.Methods("POST").Subrouter()

	get.HandleFunc("/login", unsecure(handlers.LoginHandler))

}
