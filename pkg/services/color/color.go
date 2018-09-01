package color

import "github.com/devinmcgloin/clr/clr"

const (
	Shade uint8 = iota
	SpecificColor
)

type Service interface {
	Iterate() []clr.Color
	Lookup(hex string) clr.ColorSpace
	AddColor(name, hex string) error
	AddColors(colors map[string]string) error
	Colors() (map[string]string, error)
}
