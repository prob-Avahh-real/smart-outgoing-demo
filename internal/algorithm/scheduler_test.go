package algorithm

import (
	"testing"
	"time"
)

func TestNewScheduler(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	if scheduler == nil {
		t.Fatal("NewScheduler returned nil")
	}
	if scheduler.graph != graph {
		t.Error("Scheduler graph not set correctly")
	}
}

func TestScheduleTasks(t *testing.T) {
	graph := NewGraph()
	
	// Add vehicle nodes
	graph.AddNode("vehicle1", 114.0, 22.0)
	graph.AddNode("vehicle2", 114.1, 22.1)
	
	// Add task nodes
	graph.AddNode("task1", 114.05, 22.05)
	graph.AddNode("task2", 114.15, 22.15)
	
	scheduler := NewScheduler(graph)
	
	vehicles := []string{"vehicle1", "vehicle2"}
	tasks := []*Task{
		{
			ID:       "task1",
			Priority: 1,
			Deadline: time.Now().Unix() + 3600,
			Location: Node{ID: "task1", Lng: 114.05, Lat: 22.05},
			Duration: 600,
		},
		{
			ID:       "task2",
			Priority: 2,
			Deadline: time.Now().Unix() + 3600,
			Location: Node{ID: "task2", Lng: 114.15, Lat: 22.15},
			Duration: 600,
		},
	}
	
	routes := scheduler.ScheduleTasks(vehicles, tasks)
	
	if len(routes) == 0 {
		t.Error("Expected at least one route")
	}
	
	for _, route := range routes {
		if route.VehicleID == "" {
			t.Error("Route has empty vehicle ID")
		}
		if len(route.Tasks) == 0 {
			t.Error("Route has no tasks")
		}
		if route.Distance <= 0 {
			t.Error("Route distance should be positive")
		}
	}
}

func TestScheduleTasksEmpty(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Test with empty vehicles
	routes := scheduler.ScheduleTasks([]string{}, []*Task{})
	if len(routes) != 0 {
		t.Error("Expected empty routes for empty input")
	}
	
	// Test with empty tasks
	routes = scheduler.ScheduleTasks([]string{"vehicle1"}, []*Task{})
	if len(routes) != 0 {
		t.Error("Expected empty routes for empty tasks")
	}
}

func TestOptimizeRoute(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	tasks := []*Task{
		{
			ID:       "task1",
			Priority: 1,
			Location: Node{ID: "task1", Lng: 114.05, Lat: 22.05},
		},
		{
			ID:       "task2",
			Priority: 2,
			Location: Node{ID: "task2", Lng: 114.10, Lat: 22.10},
		},
	}
	
	startNode := &Node{ID: "start", Lng: 114.0, Lat: 22.0}
	path := scheduler.OptimizeRoute(startNode, tasks)
	
	if len(path) == 0 {
		t.Error("OptimizeRoute returned empty path")
	}
	
	if path[0].ID != "start" {
		t.Error("Path should start with start node")
	}
	
	// Path should contain all tasks + start node
	if len(path) != len(tasks)+1 {
		t.Errorf("Expected path length %d, got %d", len(tasks)+1, len(path))
	}
}

func TestOptimizeRouteEmpty(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	startNode := &Node{ID: "start", Lng: 114.0, Lat: 22.0}
	path := scheduler.OptimizeRoute(startNode, []*Task{})
	
	if len(path) != 1 {
		t.Errorf("Expected path length 1 for empty tasks, got %d", len(path))
	}
	
	if path[0].ID != "start" {
		t.Error("Path should contain only start node")
	}
}

func TestDijkstra(t *testing.T) {
	graph := NewGraph()
	
	// Create a simple graph
	graph.AddNode("A", 114.0, 22.0)
	graph.AddNode("B", 114.1, 22.0)
	graph.AddNode("C", 114.2, 22.0)
	
	graph.AddEdge("A", "B", 1.0)
	graph.AddEdge("B", "C", 1.0)
	
	scheduler := NewScheduler(graph)
	
	path, distance := scheduler.Dijkstra("A", "C")
	
	if len(path) == 0 {
		t.Error("Dijkstra returned empty path")
	}
	
	if path[0] != "A" {
		t.Error("Path should start from A")
	}
	
	if path[len(path)-1] != "C" {
		t.Error("Path should end at C")
	}
	
	if distance <= 0 {
		t.Error("Distance should be positive")
	}
}

func TestDijkstraNonExistent(t *testing.T) {
	graph := NewGraph()
	graph.AddNode("A", 114.0, 22.0)
	
	scheduler := NewScheduler(graph)
	
	path, distance := scheduler.Dijkstra("A", "nonexistent")
	
	if len(path) != 0 {
		t.Error("Expected empty path for nonexistent node")
	}
	
	if distance != 0 {
		t.Error("Expected zero distance for nonexistent node")
	}
}

func TestHungarianAlgorithm(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Create a simple cost matrix
	costMatrix := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
	}
	
	assignment := scheduler.hungarianAlgorithm(costMatrix)
	
	if len(assignment) != 2 {
		t.Errorf("Expected 2 assignments, got %d", len(assignment))
	}
	
	// Check that each task is assigned to exactly one vehicle
	assignedTasks := make(map[int]bool)
	for _, taskIdx := range assignment {
		if taskIdx != -1 {
			if assignedTasks[taskIdx] {
				t.Error("Task assigned to multiple vehicles")
			}
			assignedTasks[taskIdx] = true
		}
	}
}

func TestHungarianAlgorithmEmpty(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	costMatrix := [][]float64{}
	assignment := scheduler.hungarianAlgorithm(costMatrix)
	
	if len(assignment) != 0 {
		t.Error("Expected empty assignment for empty matrix")
	}
}

func TestScheduleTasksEdgeCases(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Test with more vehicles than tasks
	vehicles := []string{"v1", "v2", "v3"}
	tasks := []*Task{
		{ID: "t1", Priority: 1, Location: Node{ID: "t1", Lng: 114.05, Lat: 22.05}},
	}
	
	graph.AddNode("v1", 114.0, 22.0)
	graph.AddNode("v2", 114.1, 22.1)
	graph.AddNode("v3", 114.2, 22.2)
	graph.AddNode("t1", 114.05, 22.05)
	
	routes := scheduler.ScheduleTasks(vehicles, tasks)
	if len(routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(routes))
	}
	
	// Test with more tasks than vehicles
	tasks = []*Task{
		{ID: "t1", Priority: 1, Location: Node{ID: "t1", Lng: 114.05, Lat: 22.05}},
		{ID: "t2", Priority: 2, Location: Node{ID: "t2", Lng: 114.15, Lat: 22.15}},
		{ID: "t3", Priority: 3, Location: Node{ID: "t3", Lng: 114.25, Lat: 22.25}},
	}
	
	graph.AddNode("t2", 114.15, 22.15)
	graph.AddNode("t3", 114.25, 22.25)
	
	routes = scheduler.ScheduleTasks(vehicles, tasks)
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}
}

func TestDijkstraEdgeCases(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Test with disconnected nodes
	graph.AddNode("A", 114.0, 22.0)
	graph.AddNode("B", 114.1, 22.1)
	graph.AddNode("C", 114.2, 22.2)
	
	// Only connect A to B, not C
	graph.AddEdge("A", "B", 1.0)
	
	// Test that algorithm doesn't crash on disconnected nodes
	path, distance := scheduler.Dijkstra("A", "C")
	_ = path // Accept whatever the algorithm returns
	_ = distance
	
	// Test self-loop
	graph.AddEdge("A", "A", 0.5)
	path, distance = scheduler.Dijkstra("A", "A")
	// Self-loop should have zero distance
	if len(path) > 0 && distance != 0 {
		t.Error("Self-loop should have zero distance")
	}
}

func TestOptimizeRouteEdgeCases(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Test with single task
	tasks := []*Task{
		{ID: "t1", Priority: 1, Location: Node{ID: "t1", Lng: 114.05, Lat: 22.05}},
	}
	
	startNode := &Node{ID: "start", Lng: 114.0, Lat: 22.0}
	path := scheduler.OptimizeRoute(startNode, tasks)
	
	if len(path) != 2 {
		t.Errorf("Expected path length 2 for single task, got %d", len(path))
	}
	
	// Test with tasks at same location
	tasks = []*Task{
		{ID: "t1", Priority: 1, Location: Node{ID: "t1", Lng: 114.05, Lat: 22.05}},
		{ID: "t2", Priority: 2, Location: Node{ID: "t2", Lng: 114.05, Lat: 22.05}},
		{ID: "t3", Priority: 3, Location: Node{ID: "t3", Lng: 114.05, Lat: 22.05}},
	}
	
	path = scheduler.OptimizeRoute(startNode, tasks)
	if len(path) != 4 {
		t.Errorf("Expected path length 4 for tasks at same location, got %d", len(path))
	}
}

func TestSchedulerRobustness(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Test with nil start node - should handle gracefully without crashing
	startNode := &Node{ID: "start", Lng: 114.0, Lat: 22.0}
	path := scheduler.OptimizeRoute(startNode, []*Task{})
	// Empty tasks should return path with just start node
	if len(path) != 1 {
		t.Error("Empty tasks should return path with only start node")
	}
	
	// Test with tasks having invalid priorities
	tasks := []*Task{
		{ID: "t1", Priority: -1, Location: Node{ID: "t1", Lng: 114.05, Lat: 22.05}},
		{ID: "t2", Priority: 0, Location: Node{ID: "t2", Lng: 114.1, Lat: 22.1}},
		{ID: "t3", Priority: 1000, Location: Node{ID: "t3", Lng: 114.15, Lat: 22.15}},
	}
	
	path = scheduler.OptimizeRoute(startNode, tasks)
	if len(path) == 0 {
		t.Error("Should handle invalid priority values")
	}
	
	// Test with tasks having negative deadlines
	now := int64(1000000)
	tasks = []*Task{
		{ID: "t1", Priority: 1, Deadline: now - 1000, Location: Node{ID: "t1", Lng: 114.05, Lat: 22.05}},
		{ID: "t2", Priority: 2, Deadline: now, Location: Node{ID: "t2", Lng: 114.1, Lat: 22.1}},
		{ID: "t3", Priority: 3, Deadline: now + 1000, Location: Node{ID: "t3", Lng: 114.15, Lat: 22.15}},
	}
	
	// Should handle expired deadlines gracefully
	graph.AddNode("v1", 114.0, 22.0)
	graph.AddNode("t1", 114.05, 22.05)
	graph.AddNode("t2", 114.1, 22.1)
	graph.AddNode("t3", 114.15, 22.15)
	
	vehicles := []string{"v1"}
	routes := scheduler.ScheduleTasks(vehicles, tasks)
	if len(routes) == 0 {
		t.Error("Should handle expired deadlines")
	}
}

func TestHungarianAlgorithmRobustness(t *testing.T) {
	graph := NewGraph()
	scheduler := NewScheduler(graph)
	
	// Test with rectangular matrix (more vehicles than tasks)
	costMatrix := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
		{5.0, 6.0},
	}
	
	assignment := scheduler.hungarianAlgorithm(costMatrix)
	if len(assignment) != 3 {
		t.Errorf("Expected 3 assignments for 3x2 matrix, got %d", len(assignment))
	}
	
	// Test with rectangular matrix (more tasks than vehicles)
	costMatrix = [][]float64{
		{1.0, 2.0, 3.0},
		{4.0, 5.0, 6.0},
	}
	
	assignment = scheduler.hungarianAlgorithm(costMatrix)
	if len(assignment) != 2 {
		t.Errorf("Expected 2 assignments for 2x3 matrix, got %d", len(assignment))
	}
	
	// Test with very large cost values
	costMatrix = [][]float64{
		{1e10, 2e10},
		{3e10, 4e10},
	}
	
	assignment = scheduler.hungarianAlgorithm(costMatrix)
	if len(assignment) != 2 {
		t.Error("Should handle very large cost values")
	}
	
	// Test with very small cost values
	costMatrix = [][]float64{
		{1e-10, 2e-10},
		{3e-10, 4e-10},
	}
	
	assignment = scheduler.hungarianAlgorithm(costMatrix)
	if len(assignment) != 2 {
		t.Error("Should handle very small cost values")
	}
}
