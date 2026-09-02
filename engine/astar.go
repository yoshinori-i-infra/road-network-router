package engine

import (
	"container/heap"
	"errors"

	"road-network-router/geo"
	"road-network-router/queue"
	"road-network-router/types"
)

var ErrPathNotFound = errors.New("path not found")

func AStar(
	graph types.Graph,
	coords map[int64]types.Coordinate,
	startID, goalID int64,
) (*types.RouteResult, error) {
	goalCoord, ok := coords[goalID]
	if !ok {
		return nil, ErrPathNotFound
	}
	startCoord, ok := coords[startID]
	if !ok {
		return nil, ErrPathNotFound
	}

	gScore := make(map[int64]float64)
	gScore[startID] = 0

	cameFrom := make(map[int64]int64)
	pq := make(queue.PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &queue.Item{
		NodeID:   startID,
		Priority: geo.Distance(startCoord, goalCoord),
	})

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*queue.Item).NodeID

		if current == goalID {
			path := []int64{current}
			totalDist := gScore[current]
			for {
				prev, exists := cameFrom[current]
				if !exists {
					break
				}
				path = append([]int64{prev}, path...)
				current = prev
			}
			return &types.RouteResult{
				Path:           path,
				DistanceMeters: totalDist,
			}, nil
		}

		for _, edge := range graph[current] {
			tentativeGScore := gScore[current] + edge.Weight

			if currentG, exists := gScore[edge.To]; !exists || tentativeGScore < currentG {
				toCoord, ok := coords[edge.To]
				if !ok {
					continue
				}

				cameFrom[edge.To] = current
				gScore[edge.To] = tentativeGScore

				hScore := geo.Distance(toCoord, goalCoord)
				heap.Push(&pq, &queue.Item{
					NodeID:   edge.To,
					Priority: tentativeGScore + hScore,
				})
			}
		}
	}

	return nil, ErrPathNotFound
}
