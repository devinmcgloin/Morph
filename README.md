# Morph

A simple photography site written in Go.

## API Reference

### Resources

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

### Routes

#### CRUD

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

##### Images
GET /images/:ID
GET /images/:ID?fields
GET /images/:ID/user
GET /images/:ID/collections
GET /images/:ID/album

POST /images

PUT /images/:ID/featured

DELETE /images/:ID
DELETE /images/:ID/featured
DELETE /images/:ID?params

PATCH /images/:ID
PATCH /images/:ID/metadata
PATCH /images/:ID/tags

##### Users
GET /users/:ID
GET /users/:ID?fields
GET /users/:ID/location

POST /users

PUT /users/:ID/avatar

DELETE /users/:ID
DELETE /users/:ID?fields

PATCH /users/:ID/bio
PATCH /users/:ID/name

##### Collections
GET /collections/:ID
GET /collections/:ID?fields
GET /collections/:ID/users
GET /collections/:ID/images

POST /collections

PUT /collections/:ID/images

DELETE /collections
DELETE /collections/:ID/images/:ID

PATCH /collections/:ID/title
PATCH /collections/:ID/description
PATCH /collections/:ID/tags

##### User Albums
GET /albums/:ID
GET /albums/:ID?fields
GET /albums/:ID/images

POST /albums

PUT /albums/:ID/images

DELETE /albums
DELETE /albums/:ID/images/:ID

PATCH /albums/:ID/title
PATCH /albums/:ID/description
PATCH /albums/:ID/tags


#### Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

GET /users
GET /images
GET /collections
GET /albums
