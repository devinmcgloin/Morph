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


## User Similarity


## Preference Search


## Freshness

Freshness is a simple decay function based upon how recently the image was
posted.

## Popularity


## Favorited Already

If an image is already favorited then it gets a 1, else a 0.
