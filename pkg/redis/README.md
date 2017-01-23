<!-- # Redis schema -->

## Images

| Key                               | Type      | Semantic Type                                  |
|:----------------------------------|:----------|:-----------------------------------------------|
| images:{shortcode}                | HSET      | publish time + metadata + publish_time + owner |
| images:{shortcode}:owner          | string    | user:{shortcode}                               |
| images:{shortcode}:favorited_by   | SortedSet | users:{shortcode} => timestamp                 |
| images:{shortcode}:collections_in | SortedSet | collections:{shortcode} => timestamp           |
| images:{shortcode}:downloads      | int       | num of downloads                               |
| images:{shortcode}:views          | int       | num of views                                   |
| images:{shortcode}:purchases      | int       | num of purchases                               |
| images:{shortcode}:can_view       | Set       | users:{shortcode}                              |
| images:{shortcode}:can_edit       | Set       | users:{shortcode}                              |
| images:{shortcode}:can_delete     | Set       | users:{shortcode}                              |
| images:{shortcode}:tags           | Set       | tag                                            |
| images:{shortcode}:labels         | SortedSet | label_id => score                              |
| images:{shortcode}:colors         | SortedSet | color_id => score                              |
| images:{shortcode}:landmarks      | SortedSet | landmark_id => score                           |
| images:featured                   | Set       | images:{shortcode} => timestamp                |
| images:new                        | List      | images:{shortcode}                             |
| images:location                   | GEO       | Lng, Lat, {shortcode}                          |

## Users

| Key                            | Type      | Semantic Type                                                     |
|:-------------------------------|:----------|:------------------------------------------------------------------|
| users:{shortcode}              | HSET      | Name, Bio, url, avatar_shortcode, Location, Email, Password, Salt |
| users:{shortcode}:views        | int       |                                                                   |
| users:{shortcode}:images       | SortedSet |                                                                   |
| users:{shortcode}:collections  | SortedSet |                                                                   |
| users:{shortcode}:purchased    | SortedSet |                                                                   |
| users:{shortcode}:downloaded   | SortedSet |                                                                   |
| users:{shortcode}:seen         | SortedSet |                                                                   |
| users:{shortcode}:followed_by  | SortedSet |                                                                   |
| users:{shortcode}:followed     | SortedSet |                                                                   |
| users:{shortcode}:favorited    | SortedSet |                                                                   |
| users:{shortcode}:favorited_by | SortedSet |                                                                   |
| users:{shortcode}:stream       | SortedSet |                                                                   |
| users:featured                 | SortedSet | users:{shortcode} => timestamp                                    |
| users:admin                    | SortedSet | users:{shortcode} => timestamp                                    |
| users:location                 | GEO       | Lng, Lat, {shortcode}                                             |
| users:emails                   | Set       | List of all user emails                                           |

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
