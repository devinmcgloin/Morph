package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cf *LoginRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Username: "username",
		&cf.Password: "password",
	}
}
