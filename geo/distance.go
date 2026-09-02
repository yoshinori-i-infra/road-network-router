package geo

import (
	"math"

	"road-network-router/types"
)

const EarthRadiusMeters = 6371000.0

func Distance(c1, c2 types.Coordinate) float64 {
	phi1 := c1.Lat * math.Pi / 180
	phi2 := c2.Lat * math.Pi / 180
	deltaPhi := (c2.Lat - c1.Lat) * math.Pi / 180
	deltaLambda := (c2.Lon - c1.Lon) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusMeters * c
}
