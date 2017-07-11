# Search Endpoints

## Default Query Values
If not provided the following url query params receive the following values:

| Param  | Value   |
|--------|---------|
| u      | nil     |
| limit  | 25      |
| offset | 0       |
| radius | 2 miles |

If u is provided for any of these endpoints it will limit the return values to
images associated with that user.

## Explore
Explore allows for viewing images related to a specific landmark, city,
arbitrary point or color.

### Landmarks

| Param  | Required |
|--------|----------|
| limit  | N        |
| offset | N        |

### Cities

| Param  | Required |
|--------|----------|
| limit  | N        |
| offset | N        |

### Geo
The geo endpoint is at `/v0/geo` and takes the following parameters:

| Param  | Required |
|--------|----------|
| lat    | Y        |
| lng    | Y        |
| radius | N        |
| limit  | N        |
| offset | N        |

### Color
| Param         | Required |
|---------------|----------|
| limit         | N        |
| offset        | N        |
| hex           | Y        |
| pixelfraction | N        |

## Hot
| Param  | Required |
|--------|----------|
| limit  | N        |
| offset | N        |


## Recent
| Param  | Required |
|--------|----------|
| limit  | N        |
| offset | N        |


