package search

import (
	"github.com/fokal/fokal/pkg/model"
)

const (
	User       = "user"
	Image      = "image"
	Tag        = "tag"
	Collection = "collections"
)

type Request struct {
	RequiredTerms []string `json:"required_terms"`
	OptionalTerms []string `json:"optional_terms"`
	ExcludedTerms []string `json:"excluded_terms"`

	Color *ColorParams `json:"color"`
	Geo   *GeoParams   `json:"geo"`

	Limit *int     `json:"limit"`
	Types []string `json:"document_types"`
	User  *string  `json:"user"`
}

type GeoParams struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Radius    float64 `json:"radius"`
}

type ColorParams struct {
	HexCode       string  `json:"hex"`
	PixelFraction float64 `json:"pixel_fraction"`
}

type TagResponse struct {
	Id         string      `json:"id"`
	Permalink  string      `json:"permalink"`
	TitleImage model.Image `json:"image"`
}

type Response struct {
	Images []model.Image `json:"images"`
	Users  []model.User  `json:"users"`
	Tags   []TagResponse `json:"tags"`
}
