# Redis schema

## Images

| Key                               | Version | Type      | Semantic Type                                  |
|:----------------------------------|:--------|:----------|:-----------------------------------------------|
| images:{shortcode}                | v0.1    | HSET      | publish time + metadata + publish_time + owner |
| images:{shortcode}:owner          | v0.1    | string    | user:{shortcode}                               |
| images:{shortcode}:favorited_by   |         | SortedSet | users:{shortcode} => timestamp                 |
| images:{shortcode}:collections_in |         | SortedSet | collections:{shortcode} => timestamp           |
| images:{shortcode}:downloads      | v0.1    | int       | num of downloads                               |
| images:{shortcode}:views          | v0.1    | Set       | num of views                                   |
| images:{shortcode}:purchases      |         | Set       | num of purchases                               |
| images:{shortcode}:can_view       | v0.1    | Set       | users:{shortcode}                              |
| images:{shortcode}:can_edit       | v0.1    | Set       | users:{shortcode}                              |
| images:{shortcode}:can_delete     | v0.1    | Set       | users:{shortcode}                              |
| images:{shortcode}:tags           |         | Set       | tag                                            |
| images:{shortcode}:labels         |         | SortedSet | label_id => score                              |
| images:{shortcode}:colors         |         | SortedSet | color_id => score                              |
| images:{shortcode}:landmarks      |         | SortedSet | landmark_id => score                           |
| images:featured                   | v0.1    | Set       | images:{shortcode} => timestamp                |
| images:new                        | v0.1    | List      | images:{shortcode}                             |
| images:location                   | v0.1    | GEO       | Lng, Lat, {shortcode}                          |

## Users

| Key                            | Version | Type      | Semantic Type                                                     |
|:-------------------------------|:--------|:----------|:------------------------------------------------------------------|
| users:{shortcode}              | v0.1    | HSET      | Name, Bio, url, avatar_shortcode, Location, Email, Password, Salt |
| users:{shortcode}:views        | v0.1    | int       |                                                                   |
| users:{shortcode}:images       | v0.1    | SortedSet |                                                                   |
| users:{shortcode}:collections  |         | SortedSet |                                                                   |
| users:{shortcode}:purchased    |         | SortedSet |                                                                   |
| users:{shortcode}:downloaded   | v0.1    | SortedSet |                                                                   |
| users:{shortcode}:seen         | v0.1    | SortedSet |                                                                   |
| users:{shortcode}:followed_by  |         | SortedSet |                                                                   |
| users:{shortcode}:followed     |         | SortedSet |                                                                   |
| users:{shortcode}:favorited    |         | SortedSet |                                                                   |
| users:{shortcode}:favorited_by |         | SortedSet |                                                                   |
| users:{shortcode}:stream       |         | SortedSet |                                                                   |
| users:featured                 | v0.1    | SortedSet | users:{shortcode} => timestamp                                    |
| users:admin                    | v0.1    | Set       | users:{shortcode}                                                 |
| users:location                 | v0.1    | GEO       | Lng, Lat, {shortcode}                                             |
| users:emails                   | v0.1    | Set       | List of all user emails                                           |

## Collections

| Key                                  | Type      | Semantic Type                        |
|:-------------------------------------|:----------|:-------------------------------------|
| collections:{shortcode}              | string    | bson.ObjectId                        |
| collections:{shortcode}:images       | SortedSet |                                      |
| collections:{shortcode}:followed_by  | SortedSet |                                      |
| collections:{shortcode}:favorited_by | SortedSet |                                      |
| collections:{shortcode}:views        | int       |                                      |
| collections:{shortcode}:view_type    | string    |                                      |
| collections:{shortcode}:owner        | string    | users:{shortcode}                    |
| collections:{shortcode}:can_view     | set       | users:{shortcode}                    |
| collections:{shortcode}:can_edit     | set       | users:{shortcode}                    |
| collections:{shortcode}:can_delete   | set       | users:{shortcode}                    |
| collections:featured                 | SortedSet | collections:{shortcode} => timestamp |
| collections:new                      | SortedSet | collections:{shortcode} => timestamp |
