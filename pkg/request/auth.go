package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type logrusinRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cf *logrusinRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Username: "username",
		&cf.Password: "password",
	}
}
