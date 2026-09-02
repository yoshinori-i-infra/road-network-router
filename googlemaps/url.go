package googlemaps

import (
	"fmt"
	"net/url"
	"strings"

	"road-network-router/types"
)

func GenerateURL(path []int64, coords map[int64]types.Coordinate) string {
	if len(path) < 2 {
		return ""
	}

	startNode := path[0]
	goalNode := path[len(path)-1]

	startCoord := coords[startNode]
	goalCoord := coords[goalNode]

	origin := fmt.Sprintf("%f,%f", startCoord.Lat, startCoord.Lon)
	destination := fmt.Sprintf("%f,%f", goalCoord.Lat, goalCoord.Lon)

	maxWaypoints := 8
	waypoints := []string{}
	intermediateNodes := path[1 : len(path)-1]

	if len(intermediateNodes) > 0 {
		step := len(intermediateNodes) / (maxWaypoints + 1)
		if step == 0 {
			step = 1
		}

		for i := step; i < len(intermediateNodes) && len(waypoints) < maxWaypoints; i += step {
			wpNode := intermediateNodes[i]
			wpCoord := coords[wpNode]
			waypoints = append(waypoints, fmt.Sprintf("%f,%f", wpCoord.Lat, wpCoord.Lon))
		}
	}

	baseURL := "https://www.google.com/maps/dir/?api=1&travelmode=driving"
	baseURL += "&origin=" + origin
	baseURL += "&destination=" + destination

	if len(waypoints) > 0 {
		joinedWaypoints := strings.Join(waypoints, "|")
		baseURL += "&waypoints=" + url.QueryEscape(joinedWaypoints)
	}

	return baseURL
}
