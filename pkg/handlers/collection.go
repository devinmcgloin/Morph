package handlers

import "net/http"

func GetCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func CreateCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

// TODO these need to pull targets from request body
func AddImageToCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func AddUserToCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func FavoriteCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func FollowCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func DeleteCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

// TODO these need to pull targets from request body
func DeleteImageFromCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func DeleteUserFromCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func UnFavoriteCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}

func UnFollowCollection(w http.ResponseWriter, r *http.Request) Response {
	return Response{Code: http.StatusNotImplemented, Message: "This should be implemented soon!"}
}
