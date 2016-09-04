# Sprioc

A simple photography site with a focus on geolocation, metadata and searching.
Use the API with your own front end or use ours!

# API Overview

## Resources

Sprioc is built around images, users and collections. Users can favorite
anything else, and follow other users and collections.

# Routes

## CRUD

Default behavior for the standard getters `GET /v0/collections/:ID` is to
include the owner field by default. This means if you ask for a collection,
you'll get links to the images in the collection, and information about the
curator, and of course the default information about the collection.

Patch requests have the following format:

```json
{
    "$set": { "metadata.make": "Olympus", "featured": false },
}
```

Currently valid commands are `$set` and `$unset`. I don't plan on adding any
others.

Ref bodies have the following format:

```json
{
  "images": [ "https://sprioc.xyz/v0/images/nKMSewUkOXBY", "https://sprioc.xyz/v0/images/lPtjPPUFVVUR" ]
}
```


### Images
| Method | Endpoint                | Body           | Returns                                                          |
|:-------|:------------------------|:---------------|:-----------------------------------------------------------------|
| GET    | /v0/images/:ID          |                | Image view that contains a filled out user field.                |
| POST   | /v0/images              | raw image data | create a new image with the authenticated user                   |
| PUT    | /v0/images/:ID/favorite |                | favorite the image                                               |
|        | /v0/images/:ID/tags     |                | Add a tag to this image                                          |
| DELETE | /v0/images/:ID          |                | delete the image, only works if the request maker owns the image |
|        | /v0/images/:ID/favorite |                | unfavorite the image                                             |
|        | /v0/images/:ID/tags     |                | remove a tag to this image                                       |
| PATCH  | /v0/images/:ID          |                | make changes to the image (See notes on patch requests)          |

### Users
| Method | Endpoint                     | Body           | Returns                                                |
|:-------|:-----------------------------|:---------------|:-------------------------------------------------------|
| GET    | /v0/users/:username          |                | full user view                                         |
|        | /v0/users/me                 |                | full user view of the logged in user                   |
| POST   | /v0/users                    |                | create new user                                        |
| PUT    | /v0/users/:username/avatar   | raw image data | update avatar image                                    |
|        | /v0/users/:username/favorite |                | favorite this user                                     |
|        | /v0/users/:username/follow   |                | follow this user                                       |
| DELETE | /v0/users/:username          |                | Delete this user account                               |
|        | /v0/users/:username/favorite |                | unfavorite this user                                   |
|        | /v0/users/:username/follow   |                | unfollow this user                                     |
| PATCH  | /v0/users/:username          |                | make changes to the user (See notes on patch requests) |

### Collections
| Method | Endpoint                      | Body                      | Returns |
|:-------|:------------------------------|:--------------------------|:--------|
| GET    | /v0/collections/:CID          |                           |         |
| POST   | /v0/collections               |                           |         |
|        | /v0/collections/:CID/images   | array of image shortcodes |         |
| PUT    | /v0/collections/:CID/favorite |                           |         |
|        | /v0/collections/:CID/follow   |                           |         |
| DELETE | /v0/collections               |                           |         |
|        | /v0/collections/:CID/images   | array of image shortcodes |         |
|        | /v0/collections/:CID/favorite |                           |         |
|        | /v0/collections/:CID/follow   |                           |         |
| PATCH  | /v0/collections/:CID          |                           |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint   | Body | Returns                                 |
|:-------|:-----------|:-----|:----------------------------------------|
| POST   | /v0/search |      |                                         |
|        | /v0/stream |      | stream of images from things you follow |

### Query Syntax
