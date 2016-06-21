# Morph

A simple photography site written in Go.

# API Overview

## Resources

* Images
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

Default behavior for the standard getters `GET /api/v0/collections/:ID` is to
include all sub fields up to a depth of 1. This means if you ask for a
collection, you'll get the image fields filled out, and not a reference, but you
will get an ID when looking at a images album.

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
| Method | Endpoint                       | Returns |
|:-------|:-------------------------------|:--------|
| GET    | /api/v0/images/:ID             |         |
|        | /api/v0/images/:ID/user        |         |
|        | /api/v0/images/:ID/collections |         |
|        | /api/v0/images/:ID/album       |         |
| POST   | /api/v0/images                 |         |
| PUT    | /api/v0/images/:ID/featured    |         |
| DELTE  | /api/v0/images/:ID             |         |
|        | /api/v0/images/:ID/featured    |         |
| PATCH  | /api/v0/images/:ID             |         |

### Users
| Method | Endpoint                         | Returns |
|:-------|:---------------------------------|:--------|
| GET    | /api/v0/users/:username          |         |
|        | /api/v0/users/:username/location |         |
| POST   | /api/v0/users                    |         |
| PUT    | /api/v0/users/:username/avatar   |         |
| DELETE | /api/v0/users/:username          |         |
| PATCH  | /api/v0/users/:username          |         |

### Collections
| Method | Endpoint                                 | Returns |
|:-------|:-----------------------------------------|:--------|
| GET    | /api/v0/collections/:CID                 |         |
|        | /api/v0/collections/:CID/users           |         |
|        | /api/v0/collections/:CID/images          |         |
| POST   | /api/v0/collections                      |         |
| PUT    | /api/v0/collections/:CID/images          |         |
|        | /api/v0/collections/:CID/users           |         |
| DELETE | /api/v0/collections                      |         |
|        | /api/v0/collections/:CID/images/:IID     |         |
|        | /api/v0/collections/:CID/users/:username |         |
| PATCH  | /api/v0/collections/:CID                 |         |

### User Albums
| Method | Endpoint                        | Returns |
|:-------|:--------------------------------|:--------|
| GET    | /api/v0/albums/:AID             |         |
|        | /api/v0/albums/:AID/images      |         |
| POST   | /api/v0/albums                  |         |
| PUT    | /api/v0/albums/:AID/images      |         |
| DELETE | /api/v0/albums                  |         |
|        | /api/v0/albums/:AID/images/:IID |         |
| PATCH  | /api/v0/albums/:AID             |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint            | Returns |
|:-------|:--------------------|:--------|
| GET    | /api/v0/users       |         |
|        | /api/v0/images      |         |
|        | /api/v0/collections |         |
|        | /api/v0/albums      |         |

# Error Codes
