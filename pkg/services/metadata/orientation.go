package metadata

import (
	"image"

	"github.com/disintegration/imaging"
)

func NormalizeOrientatation(image image.Image, orientation uint16) image.Image {
	if orientation == 6 || orientation == 5 {
		image = imaging.Rotate270(image)
	} else if orientation == 8 || orientation == 7 {
		image = imaging.Rotate90(image)
	} else if orientation == 3 || orientation == 4 {
		image = imaging.Rotate180(image)
	}

	if includes(orientation, []uint16{2, 5, 4, 7}) {
		image = imaging.FlipH(image)
	}

	return image
}

func includes(orientation uint16, options []uint16) bool {
	for _, v := range options {
		if v == orientation {
			return true
		}
	}
	return false
}
