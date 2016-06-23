# Sprioc

A simple photography site written in Go.

# API Overview

## Resources

* Images1
  * Location
  * Tags
  * MachineTags
  * Featured
  * MetaData
  * Collections
  * Album
* Users
  * Images
  * Location
* Albums
  * Images
  * Users
* Collections
  * Images
  * Users

# Routes

## CRUD

Default behavior for the standard getters GET /api/v0/collections/:ID is to
include all sub fields up to a depth of 1. This means if you ask for a
collection, you'll get the image fields filled out, and not a reference, but you
will get a reference when looking at a images album.

NOTE: You can get the collections and album two different ways. One if you hold
a reference to an image that belongs to an album or collection and the second if
you have direct references to those items themselves.

For GET requests on resources (EG: `GET /api/v0/images/:ID`) you may specify the
fields you want returned in the JSON body. If not specified you will receive the
default behavior specified above.  

```json
{
  "fields":["tags", "location", "sources"]
}
```

### Images
| Method | Endpoint                    | Body | Returns                                                            |
|:-------|:----------------------------|:-----|:-------------------------------------------------------------------|
| GET    | /api/v0/images/:ID          |      | Image view that contains a filled out user field.                  |
| POST   | /api/v0/images              |      | create a new image with the authenticated user                     |
| PUT    | /api/v0/images/:ID/featured |      | feature this image, only works if the request maker owns the image |
|        | /api/v0/images/:ID/favorite |      | favorite the image                                                 |
| DELETE | /api/v0/images/:ID          |      | delete the image, only works if the request maker owns the image   |
|        | /api/v0/images/:ID/featured |      | un-feature the image.                                              |
|        | /api/v0/images/:ID/favorite |      | unfavorite the image                                               |
| PATCH  | /api/v0/images/:ID          |      | make changes to the image (See notes on patch requests)            |

### Users
| Method | Endpoint                          | Body | Returns                                                    |
|:-------|:----------------------------------|:-----|:-----------------------------------------------------------|
| GET    | /api/v0/users/:username           |      | full user view that does not contain filled out sub fields |
| POST   | /api/v0/users                     |      | create new user                                            |
| PUT    | /api/v0/users/:username/avatar    |      | update avatar image                                        |
|        | /api/v0/images/:username/favorite |      | favorite this user                                         |
|        | /api/v0/images/:username/follow   |      | follow this user                                           |
| DELETE | /api/v0/users/:username           |      | Delete this user account                                   |
|        | /api/v0/images/:username/favorite |      | unfavorite this user                                       |
|        | /api/v0/images/:username/follow   |      | unfollow this user                                         |
| PATCH  | /api/v0/users/:username           |      | make changes to the user (See notes on patch requests)     |

### Collections
| Method | Endpoint                                 | Body | Returns |
|:-------|:-----------------------------------------|:-----|:--------|
| GET    | /api/v0/collections/:CID                 |      |         |
| POST   | /api/v0/collections                      |      |         |
| PUT    | /api/v0/collections/:CID/images          |      |         |
|        | /api/v0/collections/:CID/users           |      |         |
|        | /api/v0/collections/:CID/favorite        |      |         |
|        | /api/v0/collections/:CID/follow          |      |         |
| DELETE | /api/v0/collections                      |      |         |
|        | /api/v0/collections/:CID/images/:IID     |      |         |
|        | /api/v0/collections/:CID/users/:username |      |         |
|        | /api/v0/collections/:CID/favorite        |      |         |
|        | /api/v0/collections/:CID/follow          |      |         |
| PATCH  | /api/v0/collections/:CID                 |      |         |

### Albums
| Method | Endpoint                          | Body | Returns |
|:-------|:----------------------------------|:-----|:--------|
| GET    | /api/v0/albums/:AID               |      |         |
| POST   | /api/v0/albums                    |      |         |
| PUT    | /api/v0/albums/:AID/images        |      |         |
|        | /api/v0/collections/:AID/favorite |      |         |
|        | /api/v0/collections/:AID/follow   |      |         |
| DELETE | /api/v0/albums                    |      |         |
|        | /api/v0/albums/:AID/images/:IID   |      |         |
|        | /api/v0/collections/:AID/favorite |      |         |
|        | /api/v0/collections/:AID/follow   |      |         |
| PATCH  | /api/v0/albums/:AID               |      |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint            | Body | Returns |
|:-------|:--------------------|:-----|:--------|
| GET    | /api/v0/users       |      |         |
|        | /api/v0/images      |      |         |
|        | /api/v0/collections |      |         |
|        | /api/v0/albums      |      |         |
|        | /api/v0/search      |      |         |
|        | /api/v0/stream      |      |         |

# Error Codes
