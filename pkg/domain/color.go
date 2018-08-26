package domain

import "github.com/devinmcgloin/clr/clr"

type ColorCategory uint8

const (
	Shade uint = iota
	SpecificColor
)

type ColorService interface {
	Iterate() []clr.Color
	Lookup(hex string) clr.ColorSpace
	AddColor(name, hex string) error
	AddColors(colors map[string]string) error
	Colors() (map[string]string, error)
}
