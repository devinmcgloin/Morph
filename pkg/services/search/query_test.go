package search

import "testing"

func TestQuery(t *testing.T) {
	tables := []struct {
		RequiredTerms   []string
		OptionalTerms   []string
		ExcludedTerms   []string
		CorrectTSVector string
	}{
		{
			RequiredTerms:   []string{"mountains"},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{"city", "nyc"},
			CorrectTSVector: "(mountains) & !city & !nyc",
		},
		{
			RequiredTerms:   []string{},
			OptionalTerms:   []string{"lake"},
			ExcludedTerms:   []string{"city", "nyc"},
			CorrectTSVector: "(lake) & !city & !nyc",
		},
		{
			RequiredTerms:   []string{},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{"city", "nyc"},
			CorrectTSVector: "!city & !nyc",
		},
		{
			RequiredTerms:   []string{},
			OptionalTerms:   []string{"lake", "mountain"},
			ExcludedTerms:   []string{"city", "nyc"},
			CorrectTSVector: "(lake | mountain) & !city & !nyc",
		},
		{
			RequiredTerms:   []string{"lake", "mountain"},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{"city", "nyc"},
			CorrectTSVector: "(lake & mountain) & !city & !nyc",
		},
		{
			RequiredTerms:   []string{"lake", "mountain", "st,"},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{"city", "nyc"},
			CorrectTSVector: "(lake & mountain & st) & !city & !nyc",
		}, {
			RequiredTerms:   []string{},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{},
			CorrectTSVector: "",
		},
		{
			RequiredTerms:   []string{""},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{},
			CorrectTSVector: "",
		},
		{
			RequiredTerms:   []string{"Clay", "St", "&", "Front", "St,", "San", "Francisco,", "CA", "94111,", "USA"},
			OptionalTerms:   []string{},
			ExcludedTerms:   []string{},
			CorrectTSVector: "(Clay & St & Front & St & San & Francisco & CA & 94111 & USA)",
		},
	}

	for _, test := range tables {
		tsVector := formatQueryString(test.RequiredTerms, test.OptionalTerms, test.ExcludedTerms)
		if tsVector != test.CorrectTSVector {
			t.Errorf(`
Expected: %s
Got     : %s
Required Terms : %+v
Optional Terms : %+v
Excluded Terms : %+v`, test.CorrectTSVector, tsVector, test.RequiredTerms, test.OptionalTerms, test.ExcludedTerms)
		}
	}
}
