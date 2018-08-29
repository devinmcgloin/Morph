package request

import (
	"net/http"

	"github.com/mholt/binding"
)

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func (cf *CreateUser) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
		&cf.Email: binding.Field{
			Form:     "email",
			Required: true,
		},
		&cf.Token: binding.Field{
			Form:     "token",
			Required: true,
		},
	}
}

type PatchUser struct {
	Username  string `structs:"username,omitempty"`
	Bio       string `structs:"bio,omitempty"`
	URL       string `structs:"url,omitempty"`
	Name      string `structs:"name,omitempty"`
	Location  string `structs:"location,omitempty"`
	Instagram string `structs:"instagram,omitempty"`
	Twitter   string `structs:"twitter,omitempty"`
}

func (cf *PatchUser) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.Bio:       "bio",
		&cf.URL:       "email",
		&cf.Name:      "password",
		&cf.Location:  "location",
		&cf.Twitter:   "twitter",
		&cf.Instagram: "instagram",
		&cf.Username:  "username",
	}
}
