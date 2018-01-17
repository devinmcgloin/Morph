package search

type Rank struct {
	ID        int64   `json:"id"`
	Rank      float64 `json:"rank"`
	Type      string  `json:"type"`
	ColorDist float64 `json:"-" db:"color_dist"`
}

type ByRankColor []Rank

func (a ByRankColor) Len() int {
	return len(a)
}
func (a ByRankColor) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByRankColor) Less(i, j int) bool {
	return -a[i].Rank+a[i].ColorDist/50 < -a[j].Rank+a[j].ColorDist/50
}
