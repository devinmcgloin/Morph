| Text                                                                      | Type     | Path                                                         |
|:--------------------------------------------------------------------------|:---------|:-------------------------------------------------------------|
| add names to linked routes.                                               | __TODO__ | [cmd/sprioc/routes.go](cmd/sprioc/routes.go)                 |
| lock these routes down to alphabetical only with regex.                   | __TODO__ | [cmd/sprioc/routes.go](cmd/sprioc/routes.go)                 |
| redirect to new thing or just return random one like normal.              | __TODO__ | [cmd/sprioc/routes.go](cmd/sprioc/routes.go)                 |
| this should be formatted in json                                          | __TODO__ | [cmd/sprioc/sprioc.go](cmd/sprioc/sprioc.go)                 |
| need to think about JWT refresh                                           | __TODO__ | [pkg/authentication/auth.go](pkg/authentication/auth.go)     |
| In the future it should also spin off worker threads to                   | __TODO__ | [pkg/contentStorage/upload.go](pkg/contentStorage/upload.go) |
| this does not match properly to the mediaTypeOptions                      | __TODO__ | [pkg/contentStorage/upload.go](pkg/contentStorage/upload.go) |
| these need to pull targets from request body                              | __TODO__ | [pkg/handlers/album.go](pkg/handlers/album.go)               |
| these need to pull targets from request body                              | __TODO__ | [pkg/handlers/album.go](pkg/handlers/album.go)               |
| these need to pull targets from request body                              | __TODO__ | [pkg/handlers/collection.go](pkg/handlers/collection.go)     |
| these need to pull targets from request body                              | __TODO__ | [pkg/handlers/collection.go](pkg/handlers/collection.go)     |
| also needs to escape html tags.                                           | __TODO__ | [pkg/handlers/decoding.go](pkg/handlers/decoding.go)         |
| should validate all user inputs here.                                     | __TODO__ | [pkg/handlers/decoding.go](pkg/handlers/decoding.go)         |
| verify structure of post requests and adding things to collections.       | __TODO__ | [pkg/handlers/decoding.go](pkg/handlers/decoding.go)         |
| fill out owner information.                                               | __TODO__ | [pkg/handlers/encoding.go](pkg/handlers/encoding.go)         |
| need to change metadata return types to strings not ratios.               | __TODO__ | [pkg/handlers/encoding.go](pkg/handlers/encoding.go)         |
| need to take dbrefs and form them as links.                               | __TODO__ | [pkg/handlers/encoding.go](pkg/handlers/encoding.go)         |
| need to trim and verify usernames, passwords and emails.                  | __TODO__ | [pkg/handlers/user.go](pkg/handlers/user.go)                 |
| Could add altitude data here.                                             | __TODO__ | [pkg/metadata/exif.go](pkg/metadata/exif.go)                 |
| it would be good to have both public and private collections / images.    | __TODO__ | [pkg/model/schema.go](pkg/model/schema.go)                   |
| need to define external representations of these types and functions      | __TODO__ | [pkg/model/schema.go](pkg/model/schema.go)                   |
| need to check if modification already exists and that types are correct.  | __TODO__ | [pkg/store/common.go](pkg/store/common.go)                   |
| should say something if the operation does not do anything.               | __TODO__ | [pkg/store/common.go](pkg/store/common.go)                   |
| You can get the collections and album two different ways. One if you hold | __NOTE__ | [README.md](README.md)                                       |
