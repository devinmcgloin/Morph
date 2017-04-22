package model

const (
	Images uint32 = iota
	Users
)

type Ref struct {
	Id         uint32
	Collection uint32
}
