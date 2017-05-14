package daemon

import (
	"github.com/gorilla/mux"
	"github.com/sprioc/composer/pkg/handlers"
	"github.com/sprioc/composer/pkg/middleware"
)

func registerImageRoutes(api *mux.Router) {
	img := api.PathPrefix("/i").Subrouter()

	get := img.Methods("GET").Subrouter()
	get.HandleFunc("/featured/{limit:[0-9]+}", middleware.Unsecure(handlers.GetFeaturedImages))
	get.HandleFunc("/recent/{limit:[0-9]+}", middleware.Unsecure(handlers.GetRecentImages))
	get.HandleFunc("/{IID:[a-zA-Z]{12}}", middleware.Unsecure(handlers.GetImage))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/i", middleware.Secure(handlers.UploadImage))

	put := img.Methods("PUT").Subrouter()
	put.HandleFunc("/{IID:[a-zA-Z]{12}}/tags/{tag}", middleware.Secure(NotImplemented))
	put.HandleFunc("/{IID:[a-zA-Z]{12}}/featured", middleware.Secure(NotImplemented))
	put.HandleFunc("/{IID:[a-zA-Z]{12}}/favorite", middleware.Secure(NotImplemented))

	del := img.Methods("DELETE").Subrouter()
	del.HandleFunc("/{IID:[a-zA-Z]{12}}", middleware.Secure(handlers.DeleteImage))
	del.HandleFunc("/{IID:[a-zA-Z]{12}}/tags/{tag}", middleware.Secure(NotImplemented))
	del.HandleFunc("/{IID:[a-zA-Z]{12}}/featured", middleware.Secure(NotImplemented))
	del.HandleFunc("/{IID:[a-zA-Z]{12}}/favorite", middleware.Secure(NotImplemented))

	//patch := img.Methods("PATCH").Subrouter()
	//patch.HandleFunc("/{IID:[a-zA-Z]{12}}", middleware.Secure(handlers.ModifyImage))
}

func registerUserRoutes(api *mux.Router) {
	usr := api.PathPrefix("/u").Subrouter()

	get := usr.Methods("GET").Subrouter()
	get.HandleFunc("/me", middleware.Secure(handlers.GetLoggedInUser))
	get.HandleFunc("/{username}", middleware.Unsecure(handlers.GetUser))
	get.HandleFunc("/{username}/images", middleware.Unsecure(handlers.GetUserImages))
	get.HandleFunc("/{username}/favorite", middleware.Unsecure(NotImplemented))
	get.HandleFunc("/{username}/follow", middleware.Unsecure(NotImplemented))

	post := api.Methods("POST").Subrouter()
	post.HandleFunc("/u", middleware.Unsecure(handlers.CreateUser))

	put := usr.Methods("PUT").Subrouter()
	put.HandleFunc("/{username}/avatar", middleware.Secure(handlers.UploadAvatar))
	put.HandleFunc("/{username}/follow", middleware.Secure(NotImplemented))

	del := usr.Methods("DELETE").Subrouter()
	del.HandleFunc("/{username}", middleware.Secure(handlers.DeleteUser))
	del.HandleFunc("/{username}/follow", middleware.Secure(NotImplemented))

	// patch := usr.Methods("PATCH").Subrouter()
	// patch.HandleFunc("/{username}", middleware.Secure(handlers.ModifyUser))
}

//
// func registerCollectionRoutes(api *mux.Router) {
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
// func registerSearchRoutes(api *mux.Router) {
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
func registerLuckyRoutes(api *mux.Router) {

}

func registerAuthRoutes(api *mux.Router) {
	post := api.Methods("POST").Subrouter()

	post.HandleFunc("/get_token", middleware.Unsecure(handlers.GetToken))

}
