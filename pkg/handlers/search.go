package handlers

//
// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
//
// 	"github.com/gorilla/context"
// 	"github.com/sprioc/composer/pkg/core"
// 	"github.com/sprioc/composer/pkg/model"
// 	"github.com/sprioc/composer/pkg/qmgo"
// 	"github.com/sprioc/composer/pkg/refs"
// 	"github.com/sprioc/composer/pkg/rsp"
// )
//
// func Search(w http.ResponseWriter, r *http.Request) rsp.Response {
//
// 	var user model.User
// 	val, ok := context.GetOk(r, "auth")
// 	if !ok {
// 		user = model.User{}
// 	} else {
// 		user = val.(model.User)
//
// 	}
//
// 	var query qmgo.ImageSearch
// 	decoder := json.NewDecoder(r.Body)
// 	err := decoder.Decode(&query)
// 	if err != nil {
// 		log.Println(err)
// 		return rsp.Response{Code: http.StatusBadRequest, Message: "Invalid search request body"}
// 	}
//
// 	images, resp := core.Search(user, query)
// 	if !resp.Ok() {
// 		return resp
// 	}
//
// 	for _, img := range images {
// 		refs.FillExternalImage(img)
// 	}
//
// 	return rsp.Response{Code: http.StatusOK, Data: images}
// }
