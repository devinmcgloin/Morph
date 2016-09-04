package handlers

//
// import (
// 	"net/http"
//
// 	"github.com/gorilla/context"
// 	"github.com/sprioc/composer/pkg/core"
// 	"github.com/sprioc/composer/pkg/model"
// 	"github.com/sprioc/composer/pkg/refs"
// 	"github.com/sprioc/composer/pkg/rsp"
// )
//
// func GetStream(w http.ResponseWriter, r *http.Request) rsp.Response {
// 	var user model.User
// 	val, ok := context.GetOk(r, "auth")
// 	if !ok {
// 		return rsp.Response{Code: http.StatusUnauthorized, Message: "Must be logged in to view a user stream"}
// 	}
//
// 	user = val.(model.User)
//
// 	stream, resp := core.GetStream(user)
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	for _, img := range stream {
// 		refs.FillExternalImage(img)
// 	}
//
// 	return rsp.Response{Code: http.StatusOK, Data: stream}
// }
