package metadata

//func SetLocation(point *gj.Point) {

//if point == nil {
//return
//} else if mongo.Exists("locations", bson.M{"bounds": bson.M{"$geoIntersects": bson.M{"$geometry": point}}}) {
//log.Println("Location already found")
//return
//} else if mongo.Exists("locations", bson.M{"bounds": bson.M{"$near": bson.M{"$geometry": point, "$maxDistance": 1000}}}) {
//log.Println("Location already found")
//return
//}

//log.Println("Finding new location")

//r := &maps.GeocodingRequest{
//LatLng: &maps.LatLng{
//Lat: float64(point.Coordinates[1]),
//Lng: float64(point.Coordinates[0]),
//},
//LocationType: []maps.GeocodeAccuracy{maps.GeocodeAccuracyApproximate, maps.GeocodeAccuracyGeometricCenter},
//ResultType: []string{"point_of_interest", "airport", "natural_feature", "route",
//"neighborhood", "political"},
//}
//result, err := mapsClient.Geocode(context.Background(), r)
//if err != nil {
//log.Println(err)
//return
//}

//bounds := result[0].Geometry.Viewport

//location := model.Location{
//GoogleLoc: result[0],
//Bounds:    getBounds(bounds),
//}

//err = mongo.Create("locations", location)
//if err != nil {
//log.Println(err)
//}
//}

//func getBounds(bounds maps.LatLngBounds) gj.Polygon {
//poly := &gj.Polygon{
//Type: "Polygon",
//}

//box := []gj.Coordinate{
//gj.Coordinate{gj.CoordType(bounds.NorthEast.Lng), gj.CoordType(bounds.NorthEast.Lat)},
//gj.Coordinate{gj.CoordType(bounds.SouthWest.Lng), gj.CoordType(bounds.NorthEast.Lat)},
//gj.Coordinate{gj.CoordType(bounds.SouthWest.Lng), gj.CoordType(bounds.SouthWest.Lat)},
//gj.Coordinate{gj.CoordType(bounds.NorthEast.Lng), gj.CoordType(bounds.SouthWest.Lat)},
//gj.Coordinate{gj.CoordType(bounds.NorthEast.Lng), gj.CoordType(bounds.NorthEast.Lat)},
//}

//poly.AddCoordinates(gj.Coordinates(box))
//return *poly
//}
