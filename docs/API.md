# API Overview

## Patch Requests
Fokal works with patch requests as defined in
[RFC #7396][https://tools.ietf.org/html/rfc7396]

Consider a user resource as follows:

```json
{
  "permalink": "wRbjZBxEiZGO",
  "tags": ["lake", "forrest", "road"],
  "publish_time": "2017-05-06T06:49:30.864002-04:00",
  "last_modified": "2017-05-06T06:49:30.864002-04:00",
  "owner_link": "devinmcgloin",
  "featured": false,
  "downloads": 0,
  "views": 0,
  "aperture": "2",
  "exposure_time": "1/60",
  "focal_length": "4",
  "iso": 32,
  "make": "Apple",
  "model": "iPhone 6s",
  "lens_make": "Apple",
  "lens_model": "iPhone 6s back camera 4.15mm f/2.2",
  "pixel_xd": 4032,
  "pixel_yd": 4032,
  "capture_time": "2017-03-27T19:08:21Z"
  }
```

In order to change the tags, iso and lens_model you would send the following
json body to `/v0/i/wRbjZBxEiZGO` as a PATCH request:
```json
{
    "tags": ["boat", "lake", "forrest"],
    "iso": 100,
    "lens_model": "iPhone 6s 4.15mm f/2.2"
}
```

This would create the following document:

```json
{
  "permalink": "wRbjZBxEiZGO",
  "tags": ["lake", "forrest", "boat"],
  "publish_time": "2017-05-06T06:49:30.864002-04:00",
  "last_modified": "2017-05-06T06:49:30.864002-04:00",
  "owner_link": "devinmcgloin",
  "featured": false,
  "downloads": 0,
  "views": 0,
  "aperture": "2",
  "exposure_time": "1/60",
  "focal_length": "4",
  "iso": 100,
  "make": "Apple",
  "model": "iPhone 6s",
  "lens_make": "Apple",
  "lens_model": "iPhone 6s 4.15mm f/2.2",
  "pixel_xd": 4032,
  "pixel_yd": 4032,
  "capture_time": "2017-03-27T19:08:21Z"
  }
```

## Images
| Method | Endpoint                | Body           | Semantics                                               |
|:-------|:------------------------|:---------------|:--------------------------------------------------------|
| GET    | /v0/images/:ID          |                | Image view that contains a filled out user field.       |
| POST   | /v0/images              | raw image data | create a new image with the authenticated user          |
| PUT    | /v0/images/:ID/favorite |                | favorite the image                                      |
| DELETE | /v0/images/:ID          |                | delete the image                                        |
|        | /v0/images/:ID/favorite |                | unfavorite the image                                    |
| PATCH  | /v0/images/:ID          |                | make changes to the image (See notes on patch requests) |

For images valid patch targets are as follows: tags, aperature, exposure_time,
focal_length, iso, make, model, lens_make, lens_model, capture_time.

## Users
| Method | Endpoint                     | Body           | Semantics |
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

For users valid patch request fields are bio, url, and name. All other fields
should be modified through their specific endpoints.

## Collections

Note: Collections are not implemented yet and none of these endpoints are
active.

| Method | Endpoint                      | Body | Semantics |
|:-------|:------------------------------|:-----|:--------|
| GET    | /v0/collections/:CID          |      |         |
| POST   | /v0/collections               |      |         |
|        | /v0/collections/:CID/images   |      |         |
| PUT    | /v0/collections/:CID/favorite |      |         |
|        | /v0/collections/:CID/follow   |      |         |
| DELETE | /v0/collections               |      |         |
|        | /v0/collections/:CID/images   |      |         |
|        | /v0/collections/:CID/favorite |      |         |
|        | /v0/collections/:CID/follow   |      |         |
| PATCH  | /v0/collections/:CID          |      |         |


## Queries

All queries are called on the root path for that resource with a JSON body
indicating the query.

I'm planning on allowing sorting, filtering and searching, but the format for
these queries is TBD.

| Method | Endpoint   | Body | Returns                                 |
|:-------|:-----------|:-----|:----------------------------------------|
| POST   | /v0/search |      |                                         |
|        | /v0/stream |      | stream of images from things you follow |


