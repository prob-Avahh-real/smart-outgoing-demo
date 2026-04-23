package algorithm

import (
	"container/heap"
	"math"
	"sort"
)

// Task represents a scheduling task
type Task struct {
	ID        string
	Priority  int
	Deadline  int64
	Location  Node
	Duration  int64
}

// Route represents a vehicle route
type Route struct {
	VehicleID  string
	Tasks      []*Task
	Path       []*Node
	Distance   float64
	Duration   int64
}

// Scheduler implements vehicle scheduling algorithms
type Scheduler struct {
	graph *Graph
}

// NewScheduler creates a new scheduler
func NewScheduler(graph *Graph) *Scheduler {
	return &Scheduler{
		graph: graph,
	}
}

// ScheduleTasks assigns tasks to vehicles using Hungarian algorithm
func (s *Scheduler) ScheduleTasks(vehicles []string, tasks []*Task) []*Route {
	if len(vehicles) == 0 || len(tasks) == 0 {
		return []*Route{}
	}

	// Create cost matrix
	costMatrix := make([][]float64, len(vehicles))
	for i := range costMatrix {
		costMatrix[i] = make([]float64, len(tasks))
		for j, task := range tasks {
			vehicleNode := s.graph.Nodes[vehicles[i]]
			costMatrix[i][j] = CalculateDistance(
				vehicleNode.Lng, vehicleNode.Lat,
				task.Location.Lng, task.Location.Lat,
			)
		}
	}

	// Solve assignment problem
	assignment := s.hungarianAlgorithm(costMatrix)

	// Create routes
	routes := make([]*Route, 0)
	for vehicleIdx, taskIdx := range assignment {
		if taskIdx == -1 {
			continue
		}
		
		vehicleID := vehicles[vehicleIdx]
		task := tasks[taskIdx]
		
		route := &Route{
			VehicleID: vehicleID,
			Tasks:     []*Task{task},
			Path:      []*Node{s.graph.Nodes[vehicleID], &task.Location},
			Distance: costMatrix[vehicleIdx][taskIdx],
			Duration: task.Duration,
		}
		routes = append(routes, route)
	}

	return routes
}

// hungarianAlgorithm solves the assignment problem
func (s *Scheduler) hungarianAlgorithm(costMatrix [][]float64) []int {
	n := len(costMatrix)
	if n == 0 {
		return []int{}
	}

	m := len(costMatrix[0])
	if m == 0 {
		return []int{}
	}

	// Ensure square matrix by padding
	maxN := max(n, m)
	matrix := make([][]float64, maxN)
	for i := range matrix {
		matrix[i] = make([]float64, maxN)
		for j := range matrix[i] {
			if i < n && j < m {
				matrix[i][j] = costMatrix[i][j]
			} else {
				matrix[i][j] = math.MaxFloat64
			}
		}
	}

	// Hungarian algorithm implementation
	assignment := make([]int, maxN)
	for i := range assignment {
		assignment[i] = -1
	}

	// Track assigned tasks to prevent duplicates
	assignedTasks := make(map[int]bool)

	// Simplified greedy assignment (for production, use full Hungarian algorithm)
	for i := 0; i < n; i++ {
		minCost := math.MaxFloat64
		minJ := -1
		for j := 0; j < m; j++ {
			// Skip if task already assigned
			if assignedTasks[j] {
				continue
			}
			if matrix[i][j] < minCost {
				minCost = matrix[i][j]
				minJ = j
			}
		}
		if minJ != -1 {
			assignment[i] = minJ
			assignedTasks[minJ] = true
		}
	}

	return assignment[:n]
}

// OptimizeRoute optimizes route using nearest neighbor (simplified TSP)
func (s *Scheduler) OptimizeRoute(startNode *Node, tasks []*Task) []*Node {
	if len(tasks) == 0 {
		return []*Node{startNode}
	}

	// Sort tasks by priority (higher priority first)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority > tasks[j].Priority
	})

	path := []*Node{startNode}
	remainingTasks := make([]*Task, len(tasks))
	copy(remainingTasks, tasks)

	for len(remainingTasks) > 0 {
		currentNode := path[len(path)-1]
		nearestIdx := 0
		nearestDist := math.MaxFloat64

		for i, task := range remainingTasks {
			dist := CalculateDistance(
				currentNode.Lng, currentNode.Lat,
				task.Location.Lng, task.Location.Lat,
			)
			if dist < nearestDist {
				nearestDist = dist
				nearestIdx = i
			}
		}

		path = append(path, &remainingTasks[nearestIdx].Location)
		remainingTasks = append(remainingTasks[:nearestIdx], remainingTasks[nearestIdx+1:]...)
	}

	return path
}

// Dijkstra finds shortest path between two nodes
func (s *Scheduler) Dijkstra(start, end string) ([]string, float64) {
	if !s.graph.HasNode(start) || !s.graph.HasNode(end) {
		return []string{}, 0
	}

	distances := make(map[string]float64)
	previous := make(map[string]string)
	visited := make(map[string]bool)

	for nodeID := range s.graph.Nodes {
		distances[nodeID] = math.MaxFloat64
	}
	distances[start] = 0

	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{nodeID: start, priority: 0})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(*Item)

		if current.nodeID == end {
			break
		}

		if visited[current.nodeID] {
			continue
		}
		visited[current.nodeID] = true

		for _, edge := range s.graph.GetNeighbors(current.nodeID) {
			if visited[edge.To] {
				continue
			}

			newDist := distances[current.nodeID] + edge.Weight
			if newDist < distances[edge.To] {
				distances[edge.To] = newDist
				previous[edge.To] = current.nodeID
				heap.Push(pq, &Item{nodeID: edge.To, priority: newDist})
			}
		}
	}

	// Reconstruct path
	path := []string{}
	current := end
	for current != "" {
		path = append([]string{current}, path...)
		current = previous[current]
	}

	if len(path) == 1 && path[0] != end {
		return []string{}, 0
	}

	return path, distances[end]
}

// Priority queue implementation for Dijkstra
type Item struct {
	nodeID   string
	priority float64
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
