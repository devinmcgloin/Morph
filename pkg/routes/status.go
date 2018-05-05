package routes

import (
	"github.com/fokal/fokal-core/pkg/handler"
	"github.com/fokal/fokal-core/pkg/status"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func RegisterStatusRoutes(state *handler.State, api *mux.Router, chain alice.Chain) {

	head := api.Methods("HEAD").Subrouter()
	head.Handle("/status",
		chain.Then(handler.Handler{State: state, H: status.StatusHandler}))

}
