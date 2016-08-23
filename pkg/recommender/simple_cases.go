package recommender

import (
	"math"

	"github.com/sprioc/composer/pkg/model"
)

// isFavorited returns a 1 if the image is favorited and a 0 in all other cases.
func isFavorited(image model.Image, user model.User) float64 {
	for _, fav := range user.Favorites {
		if fav.Collection == "images" && fav.Shortcode == image.ShortCode {
			return 1
		}
	}
	return 0
}

func freshness(image model.Image) float64 {
	return 1 * math.Pow(math.E, -.05*float64(image.PublishTime.Hour()))
}

func popularity(image model.Image) float64 {
	return float64(len(image.FavoritedBy) + len(image.Collections) + image.Downloads)
}

func maxPopularity(images []*model.Image) float64 {
	currentMax := 0.0
	for _, image := range images {
		currentPopularity := popularity(*image)
		if currentPopularity > currentMax {
			currentMax = currentPopularity
		}
	}

	return currentMax
}
