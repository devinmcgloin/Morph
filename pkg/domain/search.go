package domain

type SearchService interface {
	FullSearch(req Request) Response
	GeoSearch(geo GeoParams) Response
	ColorSearch(color ColorParams) Response

	SimilarImages(imageID uint64) Response
}

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
	NE Point `json:"ne"`
	SW Point `json:"sw"`
}

type ColorParams struct {
	HexCode       string  `json:"hex"`
	PixelFraction float64 `json:"pixel_fraction"`
}

type Rank struct {
	ID   int64   `json:"id"`
	Rank float64 `json:"rank"`
}

type Response struct {
	Images  []Rank
	Streams []Rank
	Tags    []Rank
	Users   []Rank
}
