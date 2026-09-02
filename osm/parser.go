package osm

import (
	"context"
	"os"

	"road-network-router/geo"
	"road-network-router/types"

	"github.com/paulmach/osm"
	"github.com/paulmach/osm/osmpbf"
)

var defaultDrivableRoads = map[string]bool{
	"motorway":       true,
	"motorway_link":  true,
	"trunk":          true,
	"trunk_link":     true,
	"primary":        true,
	"primary_link":   true,
	"secondary":      true,
	"secondary_link": true,
	"tertiary":       true,
	"tertiary_link":  true,
	"unclassified":   true,
	"residential":    true,
	"living_street":  true,
}

func isOneway(tags osm.Tags, highwayType string) (bool, bool) {
	onewayTag := ""
	isRoundabout := false

	for _, tag := range tags {
		if tag.Key == "oneway" {
			onewayTag = tag.Value
		}
		if tag.Key == "junction" && tag.Value == "roundabout" {
			isRoundabout = true
		}
	}

	switch onewayTag {
	case "yes", "1", "true":
		return true, false
	case "-1", "reverse":
		return false, true
	case "no", "0", "false":
		return true, true
	}

	if highwayType == "motorway" || highwayType == "motorway_link" || isRoundabout {
		return true, false
	}

	return true, true
}

func LoadGraph(filepath string) (types.Graph, map[int64]types.Coordinate, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	scanner := osmpbf.New(context.Background(), f, 3)
	defer scanner.Close()

	coords := make(map[int64]types.Coordinate)
	graph := make(types.Graph)

	for scanner.Scan() {
		switch o := scanner.Object().(type) {
		case *osm.Node:
			coords[int64(o.ID)] = types.Coordinate{
				Lat: o.Lat,
				Lon: o.Lon,
			}
		case *osm.Way:
			highwayType := ""
			for _, tag := range o.Tags {
				if tag.Key == "highway" && defaultDrivableRoads[tag.Value] {
					highwayType = tag.Value
					break
				}
			}

			if highwayType != "" && len(o.Nodes) > 1 {
				forward, backward := isOneway(o.Tags, highwayType)

				for i := 0; i < len(o.Nodes)-1; i++ {
					n1ID := int64(o.Nodes[i].ID)
					n2ID := int64(o.Nodes[i+1].ID)

					c1, ok1 := coords[n1ID]
					c2, ok2 := coords[n2ID]

					if ok1 && ok2 {
						dist := geo.Distance(c1, c2)
						if forward {
							graph[n1ID] = append(graph[n1ID], types.Edge{To: n2ID, Weight: dist})
						}
						if backward {
							graph[n2ID] = append(graph[n2ID], types.Edge{To: n1ID, Weight: dist})
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return graph, coords, nil
}
