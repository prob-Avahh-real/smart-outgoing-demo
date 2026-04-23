
package algorithm

import (
	"math"
)

// Node represents a graph node (location)
type Node struct {
	ID    string
	Lng   float64
	Lat   float64
	Alt   float64 // Altitude in meters (0 for 2D)
}

// Edge represents a graph edge with weight
type Edge struct {
	From     string
	To       string
	Weight   float64
	Distance float64
}

// Graph represents a graph data structure
type Graph struct {
	Nodes map[string]*Node
	Edges map[string][]*Edge
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]*Edge),
	}
}

// AddNode adds a node to the graph
func (g *Graph) AddNode(id string, lng, lat float64) {
	g.AddNode3D(id, lng, lat, 0) // Default altitude 0 for 2D
}

// AddNode3D adds a node to the graph with altitude (3D support)
func (g *Graph) AddNode3D(id string, lng, lat, alt float64) {
	g.Nodes[id] = &Node{
		ID:  id,
		Lng: lng,
		Lat: lat,
		Alt: alt,
	}
}

// AddEdge adds an edge to the graph
func (g *Graph) AddEdge(from, to string, weight float64) {
	distance := CalculateDistance3D(
		g.Nodes[from].Lng, g.Nodes[from].Lat, g.Nodes[from].Alt,
		g.Nodes[to].Lng, g.Nodes[to].Lat, g.Nodes[to].Alt,
	)
	
	edge := &Edge{
		From:     from,
		To:       to,
		Weight:   weight,
		Distance: distance,
	}
	
	g.Edges[from] = append(g.Edges[from], edge)
}

// CalculateDistance calculates distance between two points using Haversine formula (2D)
func CalculateDistance(lng1, lat1, lng2, lat2 float64) float64 {
	return CalculateDistance3D(lng1, lat1, 0, lng2, lat2, 0)
}

// CalculateDistance3D calculates 3D distance including altitude
func CalculateDistance3D(lng1, lat1, alt1, lng2, lat2, alt2 float64) float64 {
	const earthRadius = 6371000.0 // meters
	
	// Calculate horizontal distance using Haversine formula
	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0
	deltaLat := (lat2 - lat1) * math.Pi / 180.0
	deltaLng := (lng2 - lng1) * math.Pi / 180.0
	
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	horizontalDistance := earthRadius * c
	
	// Calculate vertical distance
	verticalDistance := alt2 - alt1
	
	// Calculate 3D distance using Pythagorean theorem
	threeDDistance := math.Sqrt(horizontalDistance*horizontalDistance + verticalDistance*verticalDistance)
	
	return threeDDistance
}

// GetNeighbors returns neighbors of a node
func (g *Graph) GetNeighbors(nodeID string) []*Edge {
	return g.Edges[nodeID]
}

// HasNode checks if a node exists
func (g *Graph) HasNode(id string) bool {
	_, exists := g.Nodes[id]
	return exists
}
