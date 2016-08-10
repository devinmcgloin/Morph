# Recommender System

The recommender systems outputs a score for a specific image and user that
represents the likely hood of the user liking that image.

The system does this in memory and takes an array of images and a user to
recommend to.

## Weights

Each category is given a decimal score from 0-1 and the final score is
extrapolated from those values by multiplying by weights for each category. The
total result should not be greater than 100%.

No images that receive a 0% score should be included.

## Text Search

Text search is essentially fuzzy search over the image tags, content and
location description.

## Collaborative Filtering

Still evaluating possible solutions.

## Freshness

Freshness is a simple decay function based upon how recently the image was
posted.

## Popularity

Popularity is the max popularity of the set over the popularity of the specific
image being ranked.

## Favorited Already

If an image is already favorited then it gets a 1, else a 0.
