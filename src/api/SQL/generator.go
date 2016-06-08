package SQL

type MorphQuery struct {
	Tables  []string
	Columns []string
	Sorted  bool
	JoinBy  string
}

type MorphInsert struct {
}

func Query(m MorphQuery) string {
	return ""
}

func Insert(m MorphInsert) string {
	return ""
}
