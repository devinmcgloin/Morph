package search

import "github.com/fokal/fokal/pkg/model"

type Score struct {
	ID    int64       `json:"-"`
	Score float64     `json:"score"`
	Image model.Image `json:"image"`
}

type ByScores []Score

func (s ByScores) Len() int {
	return len(s)
}

func (s ByScores) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByScores) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}
