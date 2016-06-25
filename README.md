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

Default behavior for the standard getters GET /v0/collections/:ID is to
include all sub fields up to a depth of 1. This means if you ask for a
collection, you'll get the image fields filled out, and not a reference, but you
will get a reference when looking at a images album.

NOTE: You can get the collections and album two different ways. One if you hold
a reference to an image that belongs to an album or collection and the second if
you have direct references to those items themselves.

For GET requests on resources (EG: `GET /v0/images/:ID`) you may specify the
fields you want returned in the JSON body. If not specified you will receive the
default behavior specified above.  

```json
{
  "fields":["tags", "location", "sources"]
}
```

Patch requests have the following format:
```json
{
    "$set": { "metadata.make": "Olympus", "featured": false },
  }
```


### Images
| Method | Endpoint                | Body | Returns                                                            |
|:-------|:------------------------|:-----|:-------------------------------------------------------------------|
| GET    | /v0/images/:ID          |      | Image view that contains a filled out user field.                  |
| POST   | /v0/images              |      | create a new image with the authenticated user                     |
| PUT    | /v0/images/:ID/featured |      | feature this image, only works if the request maker owns the image |
|        | /v0/images/:ID/favorite |      | favorite the image                                                 |
| DELETE | /v0/images/:ID          |      | delete the image, only works if the request maker owns the image   |
|        | /v0/images/:ID/featured |      | un-feature the image.                                              |
|        | /v0/images/:ID/favorite |      | unfavorite the image                                               |
| PATCH  | /v0/images/:ID          |      | make changes to the image (See notes on patch requests)            |

### Users
| Method | Endpoint                      | Body | Returns                                                    |
|:-------|:------------------------------|:-----|:-----------------------------------------------------------|
| GET    | /v0/users/:username           |      | full user view that does not contain filled out sub fields |
| POST   | /v0/users                     |      | create new user                                            |
| PUT    | /v0/users/:username/avatar    |      | update avatar image                                        |
|        | /v0/images/:username/favorite |      | favorite this user                                         |
|        | /v0/images/:username/follow   |      | follow this user                                           |
| DELETE | /v0/users/:username           |      | Delete this user account                                   |
|        | /v0/images/:username/favorite |      | unfavorite this user                                       |
|        | /v0/images/:username/follow   |      | unfollow this user                                         |
| PATCH  | /v0/users/:username           |      | make changes to the user (See notes on patch requests)     |

### Collections
| Method | Endpoint                             | Body | Returns |
|:-------|:-------------------------------------|:-----|:--------|
| GET    | /v0/collections/:CID                 |      |         |
| POST   | /v0/collections                      |      |         |
| PUT    | /v0/collections/:CID/images          |      |         |
|        | /v0/collections/:CID/users           |      |         |
|        | /v0/collections/:CID/favorite        |      |         |
|        | /v0/collections/:CID/follow          |      |         |
| DELETE | /v0/collections                      |      |         |
|        | /v0/collections/:CID/images/:IID     |      |         |
|        | /v0/collections/:CID/users/:username |      |         |
|        | /v0/collections/:CID/favorite        |      |         |
|        | /v0/collections/:CID/follow          |      |         |
| PATCH  | /v0/collections/:CID                 |      |         |

### Albums
| Method | Endpoint                      | Body | Returns |
|:-------|:------------------------------|:-----|:--------|
| GET    | /v0/albums/:AID               |      |         |
| POST   | /v0/albums                    |      |         |
| PUT    | /v0/albums/:AID/images        |      |         |
|        | /v0/collections/:AID/favorite |      |         |
|        | /v0/collections/:AID/follow   |      |         |
| DELETE | /v0/albums                    |      |         |
|        | /v0/albums/:AID/images/:IID   |      |         |
|        | /v0/collections/:AID/favorite |      |         |
|        | /v0/collections/:AID/follow   |      |         |
| PATCH  | /v0/albums/:AID               |      |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint        | Body | Returns |
|:-------|:----------------|:-----|:--------|
| POST   | /v0/users       |      |         |
|        | /v0/images      |      |         |
|        | /v0/collections |      |         |
|        | /v0/albums      |      |         |
|        | /v0/search      |      |         |
|        | /v0/stream      |      |         |
