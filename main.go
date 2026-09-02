package main

import (
	"fmt"
	"os"
	"time"

	"road-network-router/engine"
	"road-network-router/googlemaps"
	"road-network-router/osm"
)

func main() {
	filepath := "shiga.osm.pbf"
	fmt.Printf("Loading graph from '%s'...\n", filepath)

	start := time.Now()
	graph, coords, err := osm.LoadGraph(filepath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading graph: %v\n", err)
		os.Exit(1)
	}

	totalEdges := 0
	for _, edges := range graph {
		totalEdges += len(edges)
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("Graph loaded in %v\n", time.Since(start))
	fmt.Printf("Nodes: %d\n", len(graph))
	fmt.Printf("Edges: %d\n", totalEdges)
	fmt.Println("----------------------------------------")

	var startNode, goalNode int64
	count := 0
	for id := range graph {
		if count == 0 {
			startNode = id
		} else if count == 1000 {
			goalNode = id
			break
		}
		count++
	}

	fmt.Printf("Start Node ID: %d\n", startNode)
	fmt.Printf("Goal Node ID:  %d\n", goalNode)

	searchStart := time.Now()
	result, err := engine.AStar(graph, coords, startNode, goalNode)
	searchElapsed := time.Since(searchStart)

	if err != nil {
		fmt.Printf("Route not found: %v (elapsed: %v)\n", err, searchElapsed)
		fmt.Println("----------------------------------------")
		return
	}

	fmt.Printf("Route found! (elapsed: %v)\n", searchElapsed)
	fmt.Printf("Total Distance: %.2f m\n", result.DistanceMeters)
	fmt.Printf("Path Node Count: %d\n", len(result.Path))

	fmt.Println("\n--- Google Maps URL ---")
	navURL := googlemaps.GenerateURL(result.Path, coords)
	fmt.Println(navURL)
	fmt.Println("----------------------------------------")
}
