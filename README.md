# Sprioc

A simple photography site with a focus on geolocation, metadata and searching.
Use the API with your own front end or use ours!

# API Overview

## Resources

Sprioc is built around images, users and collections. Users can favorite
anything else, and follow other users and collections.

# Routes

## CRUD

Default behavior for the standard getters `GET /api/v0/collections/:ID` is to
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

DBRef bodies have the following format:

```json
{
  "images": [ "https://sprioc.xyz/api/v0/images/nKMSewUkOXBY", "https://sprioc.xyz/api/v0/images/lPtjPPUFVVUR" ]
}
```


### Images
| Method | Endpoint                    | Body           | Returns                                                          |
|:-------|:----------------------------|:---------------|:-----------------------------------------------------------------|
| GET    | /api/v0/images/:ID          |                | Image view that contains a filled out user field.                |
| POST   | /api/v0/images              | raw image data | create a new image with the authenticated user                   |
| PUT    | /api/v0/images/:ID/favorite |                | favorite the image                                               |
|        | /api/v0/images/:ID/tags     |                | Add a tag to this image                                          |
| DELETE | /api/v0/images/:ID          |                | delete the image, only works if the request maker owns the image |
|        | /api/v0/images/:ID/favorite |                | unfavorite the image                                             |
|        | /api/v0/images/:ID/tags     |                | remove a tag to this image                                       |
| PATCH  | /api/v0/images/:ID          |                | make changes to the image (See notes on patch requests)          |

### Users
| Method | Endpoint                          | Body           | Returns                                                |
|:-------|:----------------------------------|:---------------|:-------------------------------------------------------|
| GET    | /api/v0/users/:username           |                | full user view                                         |
|        | /api/v0/users/me                  |                | full user view of the logged in user                   |
| POST   | /api/v0/users                     |                | create new user                                        |
| PUT    | /api/v0/users/:username/avatar    | raw image data | update avatar image                                    |
|        | /api/v0/images/:username/favorite |                | favorite this user                                     |
|        | /api/v0/images/:username/follow   |                | follow this user                                       |
| DELETE | /api/v0/users/:username           |                | Delete this user account                               |
|        | /api/v0/images/:username/favorite |                | unfavorite this user                                   |
|        | /api/v0/images/:username/follow   |                | unfollow this user                                     |
| PATCH  | /api/v0/users/:username           |                | make changes to the user (See notes on patch requests) |

### Collections
| Method | Endpoint                          | Body                      | Returns |
|:-------|:----------------------------------|:--------------------------|:--------|
| GET    | /api/v0/collections/:CID          |                           |         |
| POST   | /api/v0/collections               |                           |         |
|        | /api/v0/collections/:CID/images   | array of image shortcodes |         |
| PUT    | /api/v0/collections/:CID/favorite |                           |         |
|        | /api/v0/collections/:CID/follow   |                           |         |
| DELETE | /api/v0/collections               |                           |         |
|        | /api/v0/collections/:CID/images   | array of image shortcodes |         |
|        | /api/v0/collections/:CID/favorite |                           |         |
|        | /api/v0/collections/:CID/follow   |                           |         |
| PATCH  | /api/v0/collections/:CID          |                           |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint       | Body | Returns                                 |
|:-------|:---------------|:-----|:----------------------------------------|
| POST   | /api/v0/search |      |                                         |
|        | /api/v0/stream |      | stream of images from things you follow |

### Query Syntax
