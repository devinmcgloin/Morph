package qmgo

type Filter interface {
	Valid() bool
}

// Items for describing relationships

// Ord discribes filters that can be used for ordered datatypes.
type Ord string

func (ord Ord) Valid() bool {
	for _, op := range []string{"$eq", "$gt", "$gte", "$lt", "$lte", "$ne"} {
		if string(ord) == op {
			return true
		}
	}
	return false
}

// Geo describes filters that are used for geojson points or features.
type Geo string

func (g Geo) Valid() bool {
	for _, op := range []string{"$near"} {
		if string(g) == op {
			return true
		}
	}
	return false
}

type Eq string

func (eq Eq) Valid() bool {
	for _, op := range []string{"$eq", "$ne"} {
		if string(eq) == op {
			return true
		}
	}
	return false
}
