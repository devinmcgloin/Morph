# API Overview

### Images
| Method | Endpoint                | Body           | Semantics                                               |
|:-------|:------------------------|:---------------|:--------------------------------------------------------|
| GET    | /v0/images/:ID          |                | Image view that contains a filled out user field.       |
| POST   | /v0/images              | raw image data | create a new image with the authenticated user          |
| PUT    | /v0/images/:ID/favorite |                | favorite the image                                      |
|        | /v0/images/:ID/tags     |                | Add a tag to this image                                 |
| DELETE | /v0/images/:ID          |                | delete the image                                        |
|        | /v0/images/:ID/favorite |                | unfavorite the image                                    |
|        | /v0/images/:ID/tags     |                | remove a tag to this image                              |
| PATCH  | /v0/images/:ID          |                | make changes to the image (See notes on patch requests) |

### Users
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

### Collections

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


