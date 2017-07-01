package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type CreateUserRequest struct {
	Username string
	Email    string
	Password string
}

func (cf *CreateUserRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
		&cf.Email: binding.Field{
			Form:     "email",
			Required: true,
		},
		&cf.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
	}
}

type PatchUserRequest struct {
	Bio  string
	URL  string
	Name string
}

func (cf *PatchUserRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Bio:  "bio",
		&cf.URL:  "email",
		&cf.Name: "password",
	}
}
