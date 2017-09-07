package model

import (
	"fmt"
	"log"
)

type ReferenceType uint32

const (
	Images ReferenceType = iota
	Users
	Tags
	Labels
	Landmarks
	Collections
)

type Ref struct {
	Id         int64
	Collection ReferenceType
	Shortcode  string
}

func (r Ref) ToURL(port int, local bool) string {
	var host string
	if local {
		host = fmt.Sprintf("http://localhost:%d/v0", port)
	} else {
		host = "https://api.fok.al/v0"
	}
	switch r.Collection {
	case Users:
		return fmt.Sprintf("%s/users/%s", host, r.Shortcode)
	case Images:
		return fmt.Sprintf("%s/images/%s", host, r.Shortcode)
	case Tags:
		return fmt.Sprintf("%s/tags/%s", host, r.Shortcode)
	default:
		log.Panic("Invalid Collection Type")
	}
	return ""
}
