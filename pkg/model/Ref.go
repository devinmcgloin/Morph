package model

import (
	"fmt"
	"log"
)

const (
	Images uint32 = iota
	Users
	Tags
	Labels
	Landmarks
)

type Ref struct {
	Id         int64
	Collection uint32
	Shortcode  string
}

func (r Ref) ToURL(local bool) string {
	var host string
	if local {
		host = "http://localhost:8080"
	} else {
		host = "https://api.sprioc.xyz"
	}
	switch r.Collection {
	case Users:
		return fmt.Sprintf("%s/u/%s", host, r.Shortcode)
	case Images:
		return fmt.Sprintf("%s/i/%s", host, r.Shortcode)
	default:
		log.Panic("Invalid Collection Type")
	}
	return ""
}
