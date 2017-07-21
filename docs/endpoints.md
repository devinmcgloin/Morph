# Endpoints

## Retrieval
| Method | url                 | Semantics |
|--------|---------------------|-----------|
| GET    | `/v0/i/{id}`        |           |
| GET    | `/v0/u/{id}`        |           |
| GET    | `/v0/u/{id}/images` |           |
| GET    | `/v0/u/me`          |           |
| GET    | `/v0/t/{id}`        |           |


## Modification
| Method | url                   | Semantics |
|--------|-----------------------|-----------|
| DEL    | `/v0/i/{id}`          |           |
| PUT    | `/v0/i/{id}/featured` |           |
| DEL    | `/v0/i/{id}/featured` |           |
| PATCH  | `/v0/i/{id}`          |           |
| DEL    | `/v0/u/{id}`          |           |
| PATCH  | `/v0/u/{id}`          |           |

## Social
| Method | url                   | Semantics |
|--------|-----------------------|-----------|
| PUT    | `/v0/i/{id}/favorite` |           |
| DELETE | `/v0/i/{id}/favorite` |           |
| PUT    | `/v0/u/{id}/follow`   |           |
| DELETE | `/v0/u/{id}/follow`   |           |

## Search 
| Method | url              | Semantics |
|--------|------------------|-----------|
| GET    | `/v0/i/featured` |           |
| GET    | `/v0/i/recent`   |           |
| GET    | `/v0/i/color`    |           |
| GET    | `/v0/i/geo`      |           |
| GET    | `/v0/i/hot`      |           |
| GET    | `/v0/i/text`     |           |

## Random
| Method | url            | Semantics |
|--------|----------------|-----------|
| GET    | `/v0/i/random` |           |

## Create
| Method | url                 | Semantics |
|--------|---------------------|-----------|
| POST   | `/v0/u`             |           |
| POST   | `/v0/i`             |           |
| PUT    | `/v0/u/{ID}/avatar` |           |

## Authentication
| Method | url                | Semantics |
|--------|--------------------|-----------|
| POST   | `/v0/auth/token`   |           |
| GET    | `/v0/auth/certs`   |           |
| GET    | `/v0/auth/refresh` |           |
