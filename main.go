package main

import (
	"fmt"
	"math/rand"
)

const (
	NumAnts       = 3
	NumNodes      = 10
	NumIterations = 1
	// VehicleCapacity = 200

)

type Node struct {
	Demand   int
	Distance []int
}

type Ant struct {
	Visited           []bool
	CurrentNode       int
	Capacity          int
	RemainingCapacity int
	Route             []int
}

func main() {
	startNode := 0
	// rand.Seed(time.Now().UnixNano())
	// Configuración de condiciones específicas para cada nodo
	nodes := []Node{
		{Demand: 0, Distance: []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}},
		{Demand: 100, Distance: []int{10, 0, 15, 25, 35, 45, 55, 65, 75, 85}},
		{Demand: 100, Distance: []int{20, 15, 0, 10, 20, 30, 40, 50, 60, 70}},
		{Demand: 100, Distance: []int{30, 25, 10, 0, 15, 25, 35, 45, 55, 65}},
		{Demand: 100, Distance: []int{40, 35, 20, 15, 0, 10, 20, 30, 40, 50}},
		{Demand: 100, Distance: []int{50, 45, 30, 25, 10, 0, 10, 20, 30, 40}},
		{Demand: 50, Distance: []int{60, 55, 40, 35, 20, 10, 0, 10, 20, 30}},
		{Demand: 50, Distance: []int{70, 65, 50, 45, 30, 20, 10, 0, 10, 20}},
		{Demand: 50, Distance: []int{80, 75, 60, 55, 40, 30, 20, 10, 0, 10}},
		{Demand: 50, Distance: []int{90, 85, 70, 65, 50, 40, 30, 20, 10, 0}},
	}
	// Configuración de condiciones específicas para cada Vehiculo
	ants := []Ant{
		{Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 100, RemainingCapacity: 100},
		{Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 50, RemainingCapacity: 50},
		{Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 300, RemainingCapacity: 300},
	}

	pheromones := make([][]float64, NumNodes)
	for i := range pheromones {
		pheromones[i] = make([]float64, NumNodes)
		for j := range pheromones[i] {
			pheromones[i][j] = 1.0
		}
	}

	visitedNodes := make([]bool, NumNodes)

	for iteration := 0; iteration < NumIterations; iteration++ {
		for step := 0; step < NumNodes-1; step++ {
			for i := range ants {
				ant := &ants[i]
				if ant.RemainingCapacity <= 0 {
					ant.Route = append(ant.Route, startNode)
					continue
				}
				ant.CurrentNode = selectNextNode(ant, pheromones, nodes, visitedNodes)
				if ant.CurrentNode == startNode {
					ant.Route = append(ant.Route, startNode)
				} else {
					ant.Visited[ant.CurrentNode] = true
					visitedNodes[ant.CurrentNode] = true
					demand := nodes[ant.CurrentNode].Demand
					if ant.RemainingCapacity >= demand {
						ant.Route = append(ant.Route, ant.CurrentNode)
						ant.RemainingCapacity -= demand
					}
				}
			}
		}
		updatePheromones(pheromones, ants, nodes)
	}

	for i := range ants {
		ants[i].Route = append(ants[i].Route, startNode)
	}

	for k, ant := range ants {
		fmt.Println("Vehicle Route", k+1)
		fmt.Println(removeDuplicates(ant.Route))
	}
}

func selectNextNode(ant *Ant, pheromones [][]float64, nodes []Node, visitedNodes []bool) int {
	var totalPheromone float64
	var numAvailableNodes int

	for i := 0; i < NumNodes; i++ {
		if !ant.Visited[i] && !visitedNodes[i] {
			totalPheromone += pheromones[ant.CurrentNode][i] * (1 / float64(nodes[i].Distance[ant.CurrentNode]))
			numAvailableNodes++
		}
	}

	if numAvailableNodes == 0 {
		return ant.CurrentNode
	}

	probs := make([]float64, numAvailableNodes)
	availableNodes := make([]int, numAvailableNodes)
	index := 0
	for i := 0; i < NumNodes; i++ {
		if !ant.Visited[i] && !visitedNodes[i] {
			probs[index] = (pheromones[ant.CurrentNode][i] * (1 / float64(nodes[i].Distance[ant.CurrentNode]))) / totalPheromone
			availableNodes[index] = i
			index++
		}
	}

	selectedNodes := make([]int, 0)
	selectedCapacity := 0
	for selectedCapacity < ant.RemainingCapacity {
		r := rand.Float64()
		sum := 0.0
		selectedNode := -1
		for i := 0; i < numAvailableNodes; i++ {
			sum += probs[i]
			if r <= sum {
				selectedNode = availableNodes[i]
				break
			}
		}

		if selectedNode == -1 {
			selectedNode = availableNodes[numAvailableNodes-1]
		}

		if selectedCapacity+nodes[selectedNode].Demand <= ant.RemainingCapacity && nodes[selectedNode].Demand > 0 {
			selectedCapacity += nodes[selectedNode].Demand
			ant.Visited[selectedNode] = true
			selectedNodes = append(selectedNodes, selectedNode)
		} else {
			break
		}
	}

	if len(selectedNodes) > 0 {
		for _, node := range selectedNodes {
			ant.Capacity -= nodes[node].Demand
			ant.CurrentNode = node
		}
		if len(selectedNodes) == 1 {
			return selectedNodes[0]
		}
		return selectedNodes[len(selectedNodes)-1]
	}

	return ant.CurrentNode
}

func updatePheromones(pheromones [][]float64, ants []Ant, nodes []Node) {
	evaporation := 0.5

	for i := 0; i < NumNodes; i++ {
		for j := 0; j < NumNodes; j++ {
			if i != j {
				pheromones[i][j] *= evaporation
			}
		}
	}

	for _, ant := range ants {
		route := ant.Route
		for i := 0; i < len(route)-1; i++ {
			from := route[i]
			to := route[i+1]
			pheromones[from][to] += 1 / float64(nodes[to].Distance[from])
		}
	}
}

func removeDuplicates(route []int) []int {
	encountered := map[int]bool{}
	result := []int{}

	for _, val := range route {
		if !encountered[val] {
			encountered[val] = true
			result = append(result, val)
		}
	}

	result = append(result, 0)

	return result
}
