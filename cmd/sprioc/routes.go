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
	post.HandleFunc("/users", unsecure(handlers.Signup))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", secure(NotImplemented))
	put.HandleFunc("/{username}/favorite", secure(NotImplemented))
	put.HandleFunc("/{username}/follow", secure(NotImplemented))

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", secure(NotImplemented))
	del.HandleFunc("/{username}/favorite", secure(NotImplemented))
	del.HandleFunc("/{username}/follow", secure(NotImplemented))

	patch := usr.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{username}", secure(NotImplemented))
}

func registerCollectionRoutes(api *mux.Router) {
	col := api.PathPrefix("/collections").Subrouter()

	get := col.Methods("GET").Subrouter()
	get.HandleFunc("/{CID}", unsecure(NotImplemented))
	get.HandleFunc("/{CID}/users", unsecure(NotImplemented))
	get.HandleFunc("/{CID}/images", unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/collections", secure(NotImplemented))

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", secure(NotImplemented))
	put.HandleFunc("/{CID}/users", secure(NotImplemented))
	put.HandleFunc("/{CID}/favorite", secure(NotImplemented))
	put.HandleFunc("/{CID}/follow", secure(NotImplemented))

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", secure(NotImplemented))
	del.HandleFunc("/{CID}/images/{IID}", secure(NotImplemented))
	del.HandleFunc("/{CID}/users/{username}", secure(NotImplemented))
	del.HandleFunc("/{CID}/favorite", secure(NotImplemented))
	del.HandleFunc("/{CID}/follow", secure(NotImplemented))

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", secure(NotImplemented))
}

func registerAlbumRoutes(api *mux.Router) {
	alb := api.PathPrefix("/albums").Subrouter()

	get := alb.Methods("GET").Subrouter()
	get.HandleFunc("/{AID}", unsecure(NotImplemented))
	get.HandleFunc("/{AID}/images", unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/albums", secure(NotImplemented))

	put := alb.Methods("PUT").Subrouter()
	put.HandleFunc("/{AID}/images", secure(NotImplemented))
	put.HandleFunc("/{AID}/favorite", secure(NotImplemented))
	put.HandleFunc("/{AID}/follow", secure(NotImplemented))

	del := alb.Methods("DELETE").Subrouter()
	del.HandleFunc("/{AID}", secure(NotImplemented))
	del.HandleFunc("/{AID}/images/{IID}", secure(NotImplemented))
	del.HandleFunc("/{AID}/favorite", secure(NotImplemented))
	del.HandleFunc("/{AID}/follow", secure(NotImplemented))

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
