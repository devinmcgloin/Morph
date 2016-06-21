package main

import "github.com/gorilla/mux"

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/images").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/", NotImplemented)
	get.HandleFunc("/{ID}", NotImplemented)
	get.HandleFunc("/{ID}/user", NotImplemented)
	get.HandleFunc("/{ID}/collections", NotImplemented)
	get.HandleFunc("/{ID}/album", NotImplemented)

	post := img.Methods("POST").Subrouter()
	post.HandleFunc("/{ID}", NotImplemented)

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{ID}/featured", NotImplemented)

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{ID}", NotImplemented)
	del.HandleFunc("/{ID}/featured", NotImplemented)

	patch := img.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{ID}", NotImplemented)
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/users").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/", NotImplemented)
	get.HandleFunc("/{username}", NotImplemented)
	get.HandleFunc("/{username}/location", NotImplemented)

	post := usr.Methods("POST").Subrouter()
	post.HandleFunc("/", NotImplemented)

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", NotImplemented)

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", NotImplemented)

	patch := usr.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{username}", NotImplemented)
}

func registerCollectionRoutes(api *mux.Router) {
	col := api.PathPrefix("/collections").Subrouter()

	get := col.Methods("GET").Subrouter()
	get.HandleFunc("/", NotImplemented)
	get.HandleFunc("/{CID}", NotImplemented)
	get.HandleFunc("/{CID}/users", NotImplemented)
	get.HandleFunc("/{CID}/images", NotImplemented)

	post := col.Methods("POST").Subrouter()
	post.HandleFunc("/", NotImplemented)

	put := col.Methods("PUT").Subrouter()
	put.HandleFunc("/{CID}/images", NotImplemented)
	put.HandleFunc("/{CID}/users", NotImplemented)

	del := col.Methods("DELETE").Subrouter()
	del.HandleFunc("/{CID}", NotImplemented)
	del.HandleFunc("/{CID}/images/{IID}", NotImplemented)
	del.HandleFunc("/{CID}/users/{username}", NotImplemented)

	patch := col.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{CID}", NotImplemented)
}

func registerAlbumRoutes(api *mux.Router) {
	alb := api.PathPrefix("/collections").Subrouter()

	get := alb.Methods("GET").Subrouter()
	get.HandleFunc("/", NotImplemented)
	get.HandleFunc("/{AID}", NotImplemented)
	get.HandleFunc("/{AID}/images", NotImplemented)

	post := alb.Methods("POST").Subrouter()
	post.HandleFunc("/", NotImplemented)

	put := alb.Methods("PUT").Subrouter()
	put.HandleFunc("/{AID}/images", NotImplemented)

	del := alb.Methods("DELETE").Subrouter()
	del.HandleFunc("/{AID}", NotImplemented)
	del.HandleFunc("/{AID}/images/{IID}", NotImplemented)

	patch := alb.Methods("PATCH").Subrouter()
	patch.HandleFunc("/{AID}", NotImplemented)
}
