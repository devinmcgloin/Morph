package views

import "github.com/devinmcgloin/morph/src/content"

type Page struct {
	Img content.Img
	Src content.ImgSource
}

type Collection struct {
	Title   string
	NumImg  int
	Images  []content.Img
	Sources []content.ImgSource
}
