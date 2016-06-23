package main

import (
	"github.com/devinmcgloin/sprioc/src/handlers/images"
	"github.com/gorilla/mux"
)

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/images").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/{ID}", wrappedMiddle(images.GetImage))
	get.HandleFunc("/{ID}/user", wrappedMiddle(NotImplemented))
	get.HandleFunc("/{ID}/collections", wrappedMiddle(NotImplemented))
	get.HandleFunc("/{ID}/album", wrappedMiddle(NotImplemented))

	post := img.Methods("POST").Subrouter()
	post.HandleFunc("/", wrappedMiddle(images.UploadImage))

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{ID}/featured", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{ID}/favorite", wrappedMiddle(NotImplemented))

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{ID}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{ID}/featured", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{ID}/favorite", wrappedMiddle(NotImplemented))

	patch := img.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{ID}", wrappedMiddle(NotImplemented))
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/users").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/{username}", wrappedMiddle(NotImplemented))
	get.HandleFunc("/{username}/location", wrappedMiddle(NotImplemented))

	post := usr.Methods("POST").Subrouter()
	post.HandleFunc("/", wrappedMiddle(NotImplemented))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{username}/favorite", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{username}/follow", wrappedMiddle(NotImplemented))

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{username}/favorite", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{username}/follow", wrappedMiddle(NotImplemented))

	patch := usr.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{username}", wrappedMiddle(NotImplemented))
}

func registerCollectionRoutes(api *mux.Router) {
	col := api.PathPrefix("/collections").Subrouter()

	get := col.Methods("GET").Subrouter()
	get.HandleFunc("/{CID}", wrappedMiddle(NotImplemented))
	get.HandleFunc("/{CID}/users", wrappedMiddle(NotImplemented))
	get.HandleFunc("/{CID}/images", wrappedMiddle(NotImplemented))

	post := col.Methods("POST").Subrouter()
	post.HandleFunc("/", wrappedMiddle(NotImplemented))

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{CID}/users", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{CID}/favorite", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{CID}/follow", wrappedMiddle(NotImplemented))

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/images/{IID}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/users/{username}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/favorite", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{CID}/follow", wrappedMiddle(NotImplemented))

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", wrappedMiddle(NotImplemented))
}

func registerAlbumRoutes(api *mux.Router) {
	alb := api.PathPrefix("/collections").Subrouter()

	get := alb.Methods("GET").Subrouter()
	get.HandleFunc("/{AID}", wrappedMiddle(NotImplemented))
	get.HandleFunc("/{AID}/images", wrappedMiddle(NotImplemented))

	post := alb.Methods("POST").Subrouter()
	post.HandleFunc("/", wrappedMiddle(NotImplemented))

	put := alb.Methods("PUT").Subrouter()
	put.HandleFunc("/{AID}/images", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{AID}/favorite", wrappedMiddle(NotImplemented))
	put.HandleFunc("/{AID}/follow", wrappedMiddle(NotImplemented))

	del := alb.Methods("DELETE").Subrouter()
	del.HandleFunc("/{AID}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{AID}/images/{IID}", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{AID}/favorite", wrappedMiddle(NotImplemented))
	del.HandleFunc("/{AID}/follow", wrappedMiddle(NotImplemented))

	patch := alb.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{AID}", wrappedMiddle(NotImplemented))
}

func registerSearchRoutes(api *mux.Router) {
	get := api.Methods("GET").Subrouter()

	get.HandleFunc("/images", wrappedMiddle(NotImplemented))
	get.HandleFunc("/uers", wrappedMiddle(NotImplemented))
	get.HandleFunc("/collections", wrappedMiddle(NotImplemented))
	get.HandleFunc("/albums", wrappedMiddle(NotImplemented))
	get.HandleFunc("/search", wrappedMiddle(NotImplemented))

}

// routes that return random results for a given collection.
// TODO redirect to new thing or just return random one like normal.
func registerLuckyRoutes(api *mux.Router) {

}
