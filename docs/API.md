# API Overview

## Patch Requests
Fokal works with patch requests as defined in
[RFC #7396][https://tools.ietf.org/html/rfc7396]

Consider a user resource as follows:

```json
{
    "id": "devin",
    "permalink": "https://api.fok.al/v0/users/devin",
    "name": "Devin McGloin",
    "bio": "Working on Fokal",
    "url": "https://devinmcgloin.com",
    "instagram": "devinmcgloin",
    "twitter": "devinmcgloin",
    "location": "New York, NY",
    "avatar_links": {
        "thumb": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
        "small": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
        "medium": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
        "large": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
        "raw": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360"
    },
    "images_links": [
        "https://api.fok.al/v0/images/UmCoFGzEVUIc",
        "https://api.fok.al/v0/images/GOeGOHZrhpQl",
        "https://api.fok.al/v0/images/svrftapwsijq",
        "https://api.fok.al/v0/images/ZerXpxcUpOzP",
        "https://api.fok.al/v0/images/zTWfoZXJwlOT",
        "https://api.fok.al/v0/images/LVXlAOlmkTKR",
        "https://api.fok.al/v0/images/EtUzqRqSkZGB",
        "https://api.fok.al/v0/images/DjXiFqiVgFmX",
        "https://api.fok.al/v0/images/bVQOctiXBTPw",
        "https://api.fok.al/v0/images/tcNgvMvItujL",
        "https://api.fok.al/v0/images/kwvHKBzTXkFe",
        "https://api.fok.al/v0/images/pAZSFKAMqgMO",
        "https://api.fok.al/v0/images/aGTfljGAqmxR",
        "https://api.fok.al/v0/images/dPdrDMJZWugx",
        "https://api.fok.al/v0/images/zMANRzvPTkDF",
        "https://api.fok.al/v0/images/lzvSBEcEVBCr",
        "https://api.fok.al/v0/images/vkNdnysMjfYC",
        "https://api.fok.al/v0/images/oobIKfadRAAw",
        "https://api.fok.al/v0/images/zrfPuLsYKFPx",
        "https://api.fok.al/v0/images/LvkvNHHjSFYq",
        "https://api.fok.al/v0/images/QGLhwQZkaBfc",
        "https://api.fok.al/v0/images/XtMuhzUoietc",
        "https://api.fok.al/v0/images/pjQYthDqvzyv",
        "https://api.fok.al/v0/images/SadHrDYYRMCj",
        "https://api.fok.al/v0/images/DWzfwTCaFIXq",
        "https://api.fok.al/v0/images/edziXgZtrysX",
        "https://api.fok.al/v0/images/iFMzeQPUSLXa",
        "https://api.fok.al/v0/images/JRWZqSsZTRYg",
        "https://api.fok.al/v0/images/BuUiWPoJWVUe",
        "https://api.fok.al/v0/images/VmCJUtDkhxZK",
        "https://api.fok.al/v0/images/uotNwdQJfEqC",
        "https://api.fok.al/v0/images/SHqZVnQXJObk",
        "https://api.fok.al/v0/images/prYbsyvzkOIH",
        "https://api.fok.al/v0/images/GNxMVGjvaKwI",
        "https://api.fok.al/v0/images/XjhLQvQSdCXS",
        "https://api.fok.al/v0/images/jjXiNDWEXhdi"
    ],
    "favorite_links": [
        "https://api.fok.al/v0/images/EtUzqRqSkZGB",
        "https://api.fok.al/v0/images/XoQJBpyPkGDE",
        "https://api.fok.al/v0/images/XtMuhzUoietc",
        "https://api.fok.al/v0/images/INYKRowhMTpA",
        "https://api.fok.al/v0/images/DWzfwTCaFIXq",
        "https://api.fok.al/v0/images/VmCJUtDkhxZK",
        "https://api.fok.al/v0/images/prYbsyvzkOIH",
        "https://api.fok.al/v0/images/WKCpnyANfLLd",
        "https://api.fok.al/v0/images/jjXiNDWEXhdi"
    ],
    "featured": false,
    "admin": true,
    "created_at": "2017-06-29T05:33:39.926126Z",
    "last_modified": "2017-06-29T05:33:39.926126Z"
    }
```

And an image as: 
```json
{
    "id": "jjXiNDWEXhdi",
    "permalink": "https://api.fok.al/v0/images/jjXiNDWEXhdi",
    "publish_time": "2017-08-17T08:26:24.70716Z",
    "last_modified": "2017-08-17T08:26:24.70716Z",
    "landmarks": [],
    "colors": [
        {
            "sRGB": {
                "r": 231,
                "g": 240,
                "b": 250
            },
            "hex": "#E7F0FA",
            "hsv": {
                "h": 211,
                "s": 7,
                "v": 98
            },
            "shade": "White",
            "color_name": "Solitude",
            "pixel_fraction": 0.2048157,
            "score": 0.4184974
        },
        {
            "sRGB": {
                "r": 95,
                "g": 78,
                "b": 91
            },
            "hex": "#5F4E5B",
            "hsv": {
                "h": 314,
                "s": 17,
                "v": 37
            },
            "shade": "Grey",
            "color_name": "Don Juan",
            "pixel_fraction": 0.0039635357,
            "score": 0.1278347
        },
        {
            "sRGB": {
                "r": 124,
                "g": 111,
                "b": 122
            },
            "hex": "#7C6F7A",
            "hsv": {
                "h": 309,
                "s": 10,
                "v": 48
            },
            "shade": "Grey",
            "color_name": "Fedora",
            "pixel_fraction": 0.0024772098,
            "score": 0.10874551
        },
        {
            "sRGB": {
                "r": 63,
                "g": 46,
                "b": 57
            },
            "hex": "#3F2E39",
            "hsv": {
                "h": 321,
                "s": 26,
                "v": 24
            },
            "shade": "Black",
            "color_name": "Thunder",
            "pixel_fraction": 0.080955215,
            "score": 0.11354378
        },
        {
            "sRGB": {
                "r": 100,
                "g": 86,
                "b": 96
            },
            "hex": "#645660",
            "hsv": {
                "h": 317,
                "s": 13,
                "v": 39
            },
            "shade": "Grey",
            "color_name": "Scorpion",
            "pixel_fraction": 0.0018826793,
            "score": 0.06885082
        },
        {
            "sRGB": {
                "r": 65,
                "g": 53,
                "b": 61
            },
            "hex": "#41353D",
            "hsv": {
                "h": 320,
                "s": 18,
                "v": 25
            },
            "shade": "Black",
            "color_name": "Ship Gray",
            "pixel_fraction": 0.0049544196,
            "score": 0.06626194
        },
        {
            "sRGB": {
                "r": 120,
                "g": 104,
                "b": 118
            },
            "hex": "#786876",
            "hsv": {
                "h": 307,
                "s": 13,
                "v": 47
            },
            "shade": "Grey",
            "color_name": "Fedora",
            "pixel_fraction": 0.00019817677,
            "score": 0.04631115
        },
        {
            "sRGB": {
                "r": 200,
                "g": 195,
                "b": 205
            },
            "hex": "#C8C3CD",
            "hsv": {
                "h": 270,
                "s": 4,
                "v": 80
            },
            "shade": "White",
            "color_name": "Ghost",
            "pixel_fraction": 0.00307174,
            "score": 0.018846901
        },
        {
            "sRGB": {
                "r": 220,
                "g": 234,
                "b": 250
            },
            "hex": "#DCEAFA",
            "hsv": {
                "h": 212,
                "s": 11,
                "v": 98
            },
            "shade": "White",
            "color_name": "Link Water",
            "pixel_fraction": 0.03438367,
            "score": 0.017131416
        },
        {
            "sRGB": {
                "r": 157,
                "g": 150,
                "b": 159
            },
            "hex": "#9D969F",
            "hsv": {
                "h": 286,
                "s": 5,
                "v": 62
            },
            "shade": "Grey",
            "color_name": "Mountain Mist",
            "pixel_fraction": 0.001783591,
            "score": 0.008555592
        }
    ],
    "tags": [
        "architecture",
        "barclays center",
        "nyc",
        "brooklyn",
        "new york city",
        "new york"
    ],
    "labels": [
        {
            "description": "sky",
            "score": 0.8568742
        },
        {
            "description": "line",
            "score": 0.6771448
        },
        {
            "description": "daylighting",
            "score": 0.71713966
        },
        {
            "description": "architecture",
            "score": 0.80497634
        },
        {
            "description": "roof",
            "score": 0.69042677
        }
    ],
    "user": {
        "id": "devin",
        "permalink": "https://api.fok.al/v0/users/devin",
        "name": "Devin McGloin",
        "bio": "Working on Fokal",
        "url": "https://devinmcgloin.com",
        "instagram": "devinmcgloin",
        "twitter": "devinmcgloin",
        "location": "New York, NY",
        "avatar_links": {
            "thumb": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
            "small": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
            "medium": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
            "large": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
            "raw": "https://images.fok.al/avatar/82bc39fa-5d5d-4ea1-98c9-76cce02e0360"
        },
        "images_links": [
            "https://api.fok.al/v0/images/UmCoFGzEVUIc",
            "https://api.fok.al/v0/images/GOeGOHZrhpQl",
            "https://api.fok.al/v0/images/svrftapwsijq",
            "https://api.fok.al/v0/images/ZerXpxcUpOzP",
            "https://api.fok.al/v0/images/zTWfoZXJwlOT",
            "https://api.fok.al/v0/images/LVXlAOlmkTKR",
            "https://api.fok.al/v0/images/EtUzqRqSkZGB",
            "https://api.fok.al/v0/images/DjXiFqiVgFmX",
            "https://api.fok.al/v0/images/bVQOctiXBTPw",
            "https://api.fok.al/v0/images/tcNgvMvItujL",
            "https://api.fok.al/v0/images/kwvHKBzTXkFe",
            "https://api.fok.al/v0/images/pAZSFKAMqgMO",
            "https://api.fok.al/v0/images/aGTfljGAqmxR",
            "https://api.fok.al/v0/images/dPdrDMJZWugx",
            "https://api.fok.al/v0/images/zMANRzvPTkDF",
            "https://api.fok.al/v0/images/lzvSBEcEVBCr",
            "https://api.fok.al/v0/images/vkNdnysMjfYC",
            "https://api.fok.al/v0/images/oobIKfadRAAw",
            "https://api.fok.al/v0/images/zrfPuLsYKFPx",
            "https://api.fok.al/v0/images/LvkvNHHjSFYq",
            "https://api.fok.al/v0/images/QGLhwQZkaBfc",
            "https://api.fok.al/v0/images/XtMuhzUoietc",
            "https://api.fok.al/v0/images/pjQYthDqvzyv",
            "https://api.fok.al/v0/images/SadHrDYYRMCj",
            "https://api.fok.al/v0/images/DWzfwTCaFIXq",
            "https://api.fok.al/v0/images/edziXgZtrysX",
            "https://api.fok.al/v0/images/iFMzeQPUSLXa",
            "https://api.fok.al/v0/images/JRWZqSsZTRYg",
            "https://api.fok.al/v0/images/BuUiWPoJWVUe",
            "https://api.fok.al/v0/images/VmCJUtDkhxZK",
            "https://api.fok.al/v0/images/uotNwdQJfEqC",
            "https://api.fok.al/v0/images/SHqZVnQXJObk",
            "https://api.fok.al/v0/images/prYbsyvzkOIH",
            "https://api.fok.al/v0/images/GNxMVGjvaKwI",
            "https://api.fok.al/v0/images/XjhLQvQSdCXS",
            "https://api.fok.al/v0/images/jjXiNDWEXhdi"
        ],
        "favorite_links": [
            "https://api.fok.al/v0/images/EtUzqRqSkZGB",
            "https://api.fok.al/v0/images/XoQJBpyPkGDE",
            "https://api.fok.al/v0/images/XtMuhzUoietc",
            "https://api.fok.al/v0/images/INYKRowhMTpA",
            "https://api.fok.al/v0/images/DWzfwTCaFIXq",
            "https://api.fok.al/v0/images/VmCJUtDkhxZK",
            "https://api.fok.al/v0/images/prYbsyvzkOIH",
            "https://api.fok.al/v0/images/WKCpnyANfLLd",
            "https://api.fok.al/v0/images/jjXiNDWEXhdi"
        ],
        "featured": false,
        "admin": true,
        "created_at": "2017-06-29T05:33:39.926126Z",
        "last_modified": "2017-06-29T05:33:39.926126Z"
    },
    "featured": true,
    "favorited_by": [
        "https://api.fok.al/v0/users/devin"
    ],
    "stats": {
        "downloads": 0,
        "views": 454,
        "favorites": 1
    },
    "src_links": {
        "thumb": "https://images.fok.al/content/jjXiNDWEXhdi?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=200&fit=max",
        "small": "https://images.fok.al/content/jjXiNDWEXhdi?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=400&fit=max",
        "medium": "https://images.fok.al/content/jjXiNDWEXhdi?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy&w=1080&fit=max",
        "large": "https://images.fok.al/content/jjXiNDWEXhdi?ixlib=rb-0.3.5&q=80&fm=jpg&crop=entropy",
        "raw": "https://images.fok.al/content/jjXiNDWEXhdi"
    },
    "metadata": {
        "aperture": 4.5,
        "exposure_time": "1/125",
        "focal_length": 12,
        "iso": 200,
        "make": "OLYMPUS IMAGING CORP.",
        "model": "E-M5",
        "lens_model": "OLYMPUS M.12-50mm F3.5-6.3",
        "pixel_xd": 2048,
        "pixel_yd": 1536,
        "capture_time": "2017-05-21T15:47:24Z",
        "location": {
            "point": {
                "lat": 40.6839371,
                "lng": -73.9788628
            },
            "description": "Atlantic Av - Barclays Ctr, Brooklyn, NY, United States"
        }
    }
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


