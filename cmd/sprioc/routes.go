package main

import (
	"github.com/devinmcgloin/sprioc/src/handlers/authentication"
	"github.com/devinmcgloin/sprioc/src/handlers/images"
	"github.com/devinmcgloin/sprioc/src/handlers/users"
	"github.com/gorilla/mux"
)

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/images").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/{ID}", secureWrappedMiddle(images.GetImage))
	get.HandleFunc("/{ID}/user", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/{ID}/collections", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/{ID}/album", secureWrappedMiddle(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/images", secureWrappedMiddle(images.UploadImage))

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{ID}/featured", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{ID}/favorite", secureWrappedMiddle(NotImplemented))

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{ID}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{ID}/featured", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{ID}/favorite", secureWrappedMiddle(NotImplemented))

	patch := img.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{ID}", secureWrappedMiddle(NotImplemented))
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/users").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/{username}", secureWrappedMiddle(users.GetUserHandler))
	get.HandleFunc("/{username}/location", secureWrappedMiddle(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/users", unsecureWrappedMiddle(users.SignupHandler))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{username}/favorite", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{username}/follow", secureWrappedMiddle(NotImplemented))

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{username}/favorite", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{username}/follow", secureWrappedMiddle(NotImplemented))

	patch := usr.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{username}", secureWrappedMiddle(NotImplemented))
}

func registerCollectionRoutes(api *mux.Router) {
	col := api.PathPrefix("/collections").Subrouter()

	get := col.Methods("GET").Subrouter()
	get.HandleFunc("/{CID}", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/{CID}/users", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/{CID}/images", secureWrappedMiddle(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/collections", secureWrappedMiddle(NotImplemented))

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{CID}/users", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{CID}/favorite", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{CID}/follow", secureWrappedMiddle(NotImplemented))

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/images/{IID}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/users/{username}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/favorite", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/follow", secureWrappedMiddle(NotImplemented))

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", secureWrappedMiddle(NotImplemented))
}

func registerAlbumRoutes(api *mux.Router) {
	alb := api.PathPrefix("/albums").Subrouter()

	get := alb.Methods("GET").Subrouter()
	get.HandleFunc("/{AID}", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/{AID}/images", secureWrappedMiddle(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/albums", secureWrappedMiddle(NotImplemented))

	put := alb.Methods("PUT").Subrouter()
	put.HandleFunc("/{AID}/images", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{AID}/favorite", secureWrappedMiddle(NotImplemented))
	put.HandleFunc("/{AID}/follow", secureWrappedMiddle(NotImplemented))

	del := alb.Methods("DELETE").Subrouter()
	del.HandleFunc("/{AID}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{AID}/images/{IID}", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{AID}/favorite", secureWrappedMiddle(NotImplemented))
	del.HandleFunc("/{AID}/follow", secureWrappedMiddle(NotImplemented))

	patch := alb.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{AID}", secureWrappedMiddle(NotImplemented))
}

func registerSearchRoutes(api *mux.Router) {
	get := api.Methods("GET").Subrouter()

	get.HandleFunc("/images", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/uers", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/collections", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/albums", secureWrappedMiddle(NotImplemented))
	get.HandleFunc("/search", secureWrappedMiddle(NotImplemented))

}

// routes that return random results for a given collection.
// TODO redirect to new thing or just return random one like normal.
func registerLuckyRoutes(api *mux.Router) {

}

func registerAuthRoutes(router *mux.Router) {
	get := router.Methods("POST").Subrouter()

	get.HandleFunc("/login", unsecureWrappedMiddle(authentication.LoginHandler))

}
