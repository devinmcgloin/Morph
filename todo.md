| Text                                                                         | Type       | Path                                                         |
|:-----------------------------------------------------------------------------|:-----------|:-------------------------------------------------------------|
| add names to linked routes.                                                  | __TODO__   | [cmd/sprioc/routes.go](cmd/sprioc/routes.go)                 |
| lock these routes down to alphabetical only with regex.                      | __TODO__   | [cmd/sprioc/routes.go](cmd/sprioc/routes.go)                 |
| redirect to new thing or just return random one like normal.                 | __TODO__   | [cmd/sprioc/routes.go](cmd/sprioc/routes.go)                 |
| In the future it should also spin off worker threads to                      | __TODO__   | [pkg/contentStorage/upload.go](pkg/contentStorage/upload.go) |
| this does not match properly to the mediaTypeOptions                         | __TODO__   | [pkg/contentStorage/upload.go](pkg/contentStorage/upload.go) |
| need to implement refreshtoken                                               | __TODO__   | [pkg/core/authentication.go](pkg/core/authentication.go)     |
| need to think about JWT refresh                                              | __TODO__   | [pkg/core/authentication.go](pkg/core/authentication.go)     |
| this code is narly and slow, not sure how well the db will maintain          | __REVIEW__ | [pkg/core/delete.go](pkg/core/delete.go)                     |
| would like to tell users the request failed due to the target not existsing. | __TODO__   | [pkg/core/favorite.go](pkg/core/favorite.go)                 |
| would like to tell users the request failed due to the target not existsing. | __TODO__   | [pkg/core/follow.go](pkg/core/follow.go)                     |
| would really like to lock this down more and do more content validation.     | __TODO__   | [pkg/core/modify.go](pkg/core/modify.go)                     |
| NEED TO ABSTRACT THIS FURTHER                                                | __TODO__   | [pkg/core/upload.go](pkg/core/upload.go)                     |
| need to send more of this functionality to core                              | __TODO__   | [pkg/handlers/create.go](pkg/handlers/create.go)             |
| Could add altitude data here.                                                | __TODO__   | [pkg/metadata/exif.go](pkg/metadata/exif.go)                 |
| check if point matches anything in the db.                                   | __TODO__   | [pkg/metadata/location.go](pkg/metadata/location.go)         |
| also needs to escape html tags.                                              | __TODO__   | [pkg/middleware/decoding.go](pkg/middleware/decoding.go)     |
| should validate all user inputs here.                                        | __TODO__   | [pkg/middleware/decoding.go](pkg/middleware/decoding.go)     |
| verify structure of post requests and adding things to collections.          | __TODO__   | [pkg/middleware/decoding.go](pkg/middleware/decoding.go)     |
| fill out owner information.                                                  | __TODO__   | [pkg/middleware/encoding.go](pkg/middleware/encoding.go)     |
| need to change metadata return types to strings not ratios.                  | __TODO__   | [pkg/middleware/encoding.go](pkg/middleware/encoding.go)     |
| need to take dbrefs and form them as links.                                  | __TODO__   | [pkg/middleware/encoding.go](pkg/middleware/encoding.go)     |
| this writes null if the resp.Data is null.                                   | __TODO__   | [pkg/middleware/secure.go](pkg/middleware/secure.go)         |
| this writes null if the resp.Data is null.                                   | __TODO__   | [pkg/middleware/secure.go](pkg/middleware/secure.go)         |
| it would be good to have both public and private collections / images.       | __TODO__   | [pkg/model/schema.go](pkg/model/schema.go)                   |
| need to define external representations of these types and functions         | __TODO__   | [pkg/model/schema.go](pkg/model/schema.go)                   |
| need to check if modification already exists and that types are correct.     | __TODO__   | [pkg/store/common.go](pkg/store/common.go)                   |
| should say something if the operation does not do anything.                  | __TODO__   | [pkg/store/common.go](pkg/store/common.go)                   |
