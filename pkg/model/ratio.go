package model

import "fmt"

type Ratio struct {
	Num int64  `bson:"num" json:"num"`
	Den int64  `bson:"den" json:"den"`
	Rep string `bson:"rep" json:"rep"`
}

func (r Ratio) Format() string {
	return fmt.Sprintf("%d", (r.Num / r.Den))
}
