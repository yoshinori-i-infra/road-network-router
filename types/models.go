package types

type Coordinate struct {
	Lat float64
	Lon float64
}

type Edge struct {
	To     int64
	Weight float64
}

type Graph map[int64][]Edge

type RouteResult struct {
	Path           []int64
	DistanceMeters float64
}
