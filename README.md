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

Default behavior for the standard getters `GET /collections/:ID` is to include
all sub fields up to a depth of 1. This means if you ask for a collection,
you'll get the image fields filled out, and not a reference, but you will get an
ID when looking at a images album.

NOTE: You can get the collections and album two different ways. One if you hold
a reference to an image that belongs to an album or collection and the second if
you have direct references to those items themselves.

For GET requests on resources (EG: `GET /images/:ID`) you may specify the fields
you want returned, and you'll only get the data you need.

```json
{
  "fields":["tags", "location", "sources"]
}
```

### Images
| Method | Endpoint                | Returns |
|:-------|:------------------------|:--------|
| GET    | /images/:ID             |         |
|        | /images/:ID             |         |
|        | /images/:ID/user        |         |
|        | /images/:ID/collections |         |
|        | /images/:ID/album       |         |
| POST   | /images                 |         |
| PUT    | /images/:ID/featured    |         |
| DELTE  | /images/:ID             |         |
|        | /images/:ID/featured    |         |
|        | /images/:ID?params      |         |
| PATCH  | /images/:ID             |         |
|        | /images/:ID/metadata    |         |
|        | /images/:ID/tags        |         |

### Users
| Method | Endpoint            | Returns |
|:-------|:--------------------|:--------|
| GET    | /users/:ID          |         |
|        | /users/:ID?fields   |         |
|        | /users/:ID/location |         |
| POST   | /users              |         |
| PUT    | /users/:ID/avatar   |         |
| DELETE | /users/:ID          |         |
|        | /users/:ID          |         |
| PATCH  | /users/:ID/bio      |         |
|        | /users/:ID/name     |         |

### Collections
| Method | Endpoint                     | Returns |
|:-------|:-----------------------------|:--------|
| GET    | /collections/:ID             |         |
|        | /collections/:ID?fields      |         |
|        | /collections/:ID/users       |         |
|        | /collections/:ID/images      |         |
| POST   | /collections                 |         |
| PUT    | /collections/:ID/images      |         |
| DELETE | /collections                 |         |
|        | /collections/:ID/images/:ID  |         |
| PATCH  | /collections/:ID/title       |         |
|        | /collections/:ID/description |         |
|        | /collections/:ID/tags        |         |

### User Albums
| Method | Endpoint                | Returns |
|:-------|:------------------------|:--------|
| GET    | /albums/:ID             |         |
|        | /albums/:ID?fields      |         |
|        | /albums/:ID/images      |         |
| POST   | /albums                 |         |
| PUT    | /albums/:ID/images      |         |
| DELETE | /albums                 |         |
|        | /albums/:ID/images/:ID  |         |
| PATCH  | /albums/:ID/title       |         |
|        | /albums/:ID/description |         |
|        | /albums/:ID/tags        |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint     | Returns |
|:-------|:-------------|:--------|
| GET    | /users       |         |
|        | /images      |         |
|        | /collections |         |
|        | /albums      |         |

# Error Codes
