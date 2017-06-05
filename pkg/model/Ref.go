package model

import (
	"fmt"
	"log"
)

const (
	Images uint32 = iota
	Users
	Tags
)

type Ref struct {
	Id         int64
	Collection uint32
	Shortcode  string
}

func (r Ref) ToURL() string {
	switch r.Collection {
	case Users:
		return fmt.Sprintf("http://localhost:8080/u/%s", r.Shortcode)
	case Images:
		return fmt.Sprintf("http://localhost:8080/i/%s", r.Shortcode)
	default:
		log.Panic("Invalid Collection Type")
	}
	return ""
}
