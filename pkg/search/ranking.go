package search

type Rank struct {
	ID   int64   `json:"id"`
	Rank float64 `json:"rank"`
	Type string  `json:"type"`
}
