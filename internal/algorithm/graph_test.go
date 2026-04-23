package algorithm

import (
	"testing"
)

func TestNewGraph(t *testing.T) {
	graph := NewGraph()
	if graph == nil {
		t.Fatal("NewGraph returned nil")
	}
	if graph.Nodes == nil {
		t.Error("Nodes map is nil")
	}
	if graph.Edges == nil {
		t.Error("Edges map is nil")
	}
}

func TestAddNode(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("node1", 114.0, 22.0)
	
	if !graph.HasNode("node1") {
		t.Error("Node not added to graph")
	}
	
	node := graph.Nodes["node1"]
	if node.ID != "node1" {
		t.Errorf("Expected node ID 'node1', got '%s'", node.ID)
	}
	if node.Lng != 114.0 {
		t.Errorf("Expected lng 114.0, got %f", node.Lng)
	}
	if node.Lat != 22.0 {
		t.Errorf("Expected lat 22.0, got %f", node.Lat)
	}
}

func TestAddEdge(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("node1", 114.0, 22.0)
	graph.AddNode("node2", 114.1, 22.1)
	graph.AddEdge("node1", "node2", 1.0)
	
	neighbors := graph.GetNeighbors("node1")
	if len(neighbors) != 1 {
		t.Errorf("Expected 1 neighbor, got %d", len(neighbors))
	}
	
	if neighbors[0].To != "node2" {
		t.Errorf("Expected edge to 'node2', got '%s'", neighbors[0].To)
	}
}

func TestCalculateDistance(t *testing.T) {
	// Test distance between two known points
	// Shenzhen to Guangzhou approximate distance
	distance := CalculateDistance(114.0579, 22.5431, 113.2644, 23.1291)
	
	if distance <= 0 {
		t.Error("Distance should be positive")
	}
	
	// Distance should be approximately 100km (100,000 meters)
	if distance < 80000 || distance > 150000 {
		t.Errorf("Distance seems incorrect: %f meters", distance)
	}
	
	// Test distance to same point
	sameDistance := CalculateDistance(114.0, 22.0, 114.0, 22.0)
	if sameDistance != 0 {
		t.Errorf("Distance to same point should be 0, got %f", sameDistance)
	}
}

func TestCalculateDistance3D(t *testing.T) {
	// Test 3D distance with altitude
	distance := CalculateDistance3D(114.0, 22.0, 0, 114.1, 22.0, 100)
	
	if distance <= 0 {
		t.Error("3D distance should be positive")
	}
	
	// 3D distance should be greater than 2D distance
	distance2D := CalculateDistance(114.0, 22.0, 114.1, 22.0)
	if distance <= distance2D {
		t.Errorf("3D distance (%f) should be greater than 2D distance (%f)", distance, distance2D)
	}
	
	// Test same point with altitude
	samePointDistance := CalculateDistance3D(114.0, 22.0, 0, 114.0, 22.0, 0)
	if samePointDistance != 0 {
		t.Errorf("Distance to same point should be 0, got %f", samePointDistance)
	}
	
	// Test vertical only distance
	verticalDistance := CalculateDistance3D(114.0, 22.0, 0, 114.0, 22.0, 100)
	if verticalDistance != 100 {
		t.Errorf("Vertical distance should be 100, got %f", verticalDistance)
	}
}

func TestGetNeighbors(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("node1", 114.0, 22.0)
	graph.AddNode("node2", 114.1, 22.1)
	graph.AddNode("node3", 114.2, 22.2)
	
	graph.AddEdge("node1", "node2", 1.0)
	graph.AddEdge("node1", "node3", 2.0)
	
	neighbors := graph.GetNeighbors("node1")
	if len(neighbors) != 2 {
		t.Errorf("Expected 2 neighbors, got %d", len(neighbors))
	}
	
	neighbors2 := graph.GetNeighbors("node2")
	if len(neighbors2) != 0 {
		t.Errorf("Expected 0 neighbors for node2, got %d", len(neighbors2))
	}
}

func TestHasNode(t *testing.T) {
	graph := NewGraph()
	
	if graph.HasNode("nonexistent") {
		t.Error("Should not have nonexistent node")
	}
	
	graph.AddNode("node1", 114.0, 22.0)
	if !graph.HasNode("node1") {
		t.Error("Should have node1")
	}
}

func TestAddNodeEdgeCases(t *testing.T) {
	graph := NewGraph()
	
	// Test adding node with zero coordinates
	graph.AddNode("zero", 0, 0)
	if !graph.HasNode("zero") {
		t.Error("Should accept zero coordinates")
	}
	
	// Test adding node with negative coordinates
	graph.AddNode("negative", -180, -90)
	if !graph.HasNode("negative") {
		t.Error("Should accept negative coordinates")
	}
	
	// Test adding node with extreme coordinates
	graph.AddNode("extreme", 180, 90)
	if !graph.HasNode("extreme") {
		t.Error("Should accept extreme coordinates")
	}
}

func TestAddEdgeEdgeCases(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("node1", 114.0, 22.0)
	graph.AddNode("node2", 114.1, 22.1)
	
	// Test edge with zero weight
	graph.AddEdge("node1", "node2", 0)
	neighbors := graph.GetNeighbors("node1")
	if len(neighbors) != 1 || neighbors[0].Weight != 0 {
		t.Error("Should accept zero weight")
	}
	
	// Test edge with negative weight
	graph.AddEdge("node1", "node2", -1)
	neighbors = graph.GetNeighbors("node1")
	if len(neighbors) != 2 {
		t.Error("Should accept negative weight")
	}
}

func TestDistanceEdgeCases(t *testing.T) {
	// Test distance to same point with altitude
	samePoint := CalculateDistance3D(114.0, 22.0, 100, 114.0, 22.0, 100)
	if samePoint != 0 {
		t.Errorf("Distance to same point should be 0, got %f", samePoint)
	}
	
	// Test distance with negative altitude
	negativeAlt := CalculateDistance3D(114.0, 22.0, 100, 114.0, 22.0, -50)
	if negativeAlt != 150 {
		t.Errorf("Distance with negative altitude should be 150, got %f", negativeAlt)
	}
	
	// Test distance with very large altitude
	largeAlt := CalculateDistance3D(114.0, 22.0, 0, 114.0, 22.0, 10000)
	if largeAlt != 10000 {
		t.Errorf("Distance with large altitude should be 10000, got %f", largeAlt)
	}
}

func TestGraphRobustness(t *testing.T) {
	graph := NewGraph()
	
	// Test adding duplicate nodes (should not crash)
	graph.AddNode("duplicate", 114.0, 22.0)
	graph.AddNode("duplicate", 114.1, 22.1)
	
	// Second add should overwrite first
	node := graph.Nodes["duplicate"]
	if node.Lng != 114.1 || node.Lat != 22.1 {
		t.Error("Duplicate node should be overwritten")
	}
	
	// Test adding multiple edges between same nodes
	graph.AddNode("A", 114.0, 22.0)
	graph.AddNode("B", 114.1, 22.1)
	graph.AddEdge("A", "B", 1.0)
	graph.AddEdge("A", "B", 2.0)
	
	neighbors := graph.GetNeighbors("A")
	if len(neighbors) != 2 {
		t.Error("Should allow multiple edges between same nodes")
	}
	
	// Test getting neighbors of non-existent node (should not crash)
	neighbors = graph.GetNeighbors("nonexistent")
	if neighbors != nil {
		t.Error("Neighbors of non-existent node should be nil")
	}
}

func TestDistanceRobustness(t *testing.T) {
	// Test with invalid coordinates (should not crash)
	distance := CalculateDistance(999, 999, -999, -999)
	if distance < 0 {
		t.Error("Distance calculation should not return negative for invalid coords")
	}
	
	// Test with very large coordinates
	distance = CalculateDistance(1000, 1000, -1000, -1000)
	if distance < 0 {
		t.Error("Distance calculation should handle large coordinates")
	}
	
	// Test 3D distance with extreme altitudes
	distance = CalculateDistance3D(0, 0, -100000, 0, 0, 100000)
	if distance != 200000 {
		t.Errorf("3D distance with extreme altitudes should be 200000, got %f", distance)
	}
}
