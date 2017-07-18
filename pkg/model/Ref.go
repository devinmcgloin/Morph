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

func (r Ref) ToURL(port int, local bool) string {
	var host string
	if local {
		host = fmt.Sprintf("http://localhost:%d/v0", port)
	} else {
		host = "https://api.sprioc.xyz/v0"
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
