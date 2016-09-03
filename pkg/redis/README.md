# Redis schema

## Images

| Key                               | Type      | Semantic Meaning                     |
|:----------------------------------|:-------------------------------------------------|
| images:{shortcode}                | string    | bson.ObjectId                        |
| images:{shortcode}:favorited_by   | SortedSet | users:{shortcode} => timestamp       |
| images:{shortcode}:collections_in | SortedSet | collections:{shortcode} => timestamp |
| images:{shortcode}:downloads      | int       | num of downloads                     |
| images:{shortcode}:views          | int       | num of views                         |
| images:{shortcode}:purchases      | int       | num of purchases                     |
| images:{shortcode}:owner          | string    | users:{shortcode}                    |
| images:{shortcode}:can_view       | set       | users:{shortcode}                    |
| images:{shortcode}:can_edit       | set       | users:{shortcode}                    |
| images:{shortcode}:can_delete     | set       | users:{shortcode}                    |

## Users

| Key                            | Type      | Semantic Meaning              |
|:-------------------------------|:------------------------------------------|
| users:{shortcode}              | HSET      | bson.ObjectId, Password, Salt |
| users:{shortcode}:admin        | bool      |                               |
| users:{shortcode}:views        | int       |                               |
| users:{shortcode}:images       | SortedSet |                               |
| users:{shortcode}:collections  | SortedSet |                               |
| users:{shortcode}:purchased    | SortedSet |                               |
| users:{shortcode}:downloaded   | SortedSet |                               |
| users:{shortcode}:seen         | SortedSet |                               |
| users:{shortcode}:followed_by  | SortedSet |                               |
| users:{shortcode}:followed     | SortedSet |                               |
| users:{shortcode}:favorited    | SortedSet |                               |
| users:{shortcode}:favorited_by | SortedSet |                               |
| users:{shortcode}:stream       | SortedSet |                               |

## Collections

| Key                                  | Type      | Semantic Meaning  |
|:-------------------------------------|:------------------------------|
| collections:{shortcode}              | string    | bson.ObjectId     |
| collections:{shortcode}:images       | SortedSet |                   |
| collections:{shortcode}:followed_by  | SortedSet |                   |
| collections:{shortcode}:favorited_by | SortedSet |                   |
| collections:{shortcode}:views        | int       |                   |
| collections:{shortcode}:view_type    | string    |                   |
| collections:{shortcode}:owner        | string    |                   |
| collections:{shortcode}:can_view     | set       | users:{shortcode} |
| collections:{shortcode}:can_edit     | set       | users:{shortcode} |
| collections:{shortcode}:can_delete   | set       | users:{shortcode} |
