package main

import (
	"fmt"
	"math/rand"
	antcolonyalgorithm "routng/antColonyAlgorithm"
	"routng/core"
	mathematicalmodel "routng/mathematicalModel"
	"routng/models"
	"time"
)

const (
	NumAnts       = 3
	NumNodes      = 10
	NumIterations = 1
)

func main() {
	startNode := 0
	rand.Seed(time.Now().UnixNano())
	// Configuración de condiciones específicas para cada nodo
	nodes := []models.Node{
		{Id: 0, Demand: 0, ServiceTime: 10, Distance: []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}},
		{Id: 1, Demand: 100, ServiceTime: 10, Distance: []int{10, 0, 15, 25, 35, 45, 55, 65, 75, 85}},
		{Id: 2, Demand: 100, ServiceTime: 10, Distance: []int{20, 15, 0, 10, 20, 30, 40, 50, 60, 70}},
		{Id: 3, Demand: 100, ServiceTime: 10, Distance: []int{30, 25, 10, 0, 15, 25, 35, 45, 55, 65}},
		{Id: 4, Demand: 100, ServiceTime: 10, Distance: []int{40, 35, 20, 15, 0, 10, 20, 30, 40, 50}},
		{Id: 5, Demand: 100, ServiceTime: 10, Distance: []int{50, 45, 30, 25, 10, 0, 10, 20, 30, 40}},
		{Id: 6, Demand: 100, ServiceTime: 10, Distance: []int{60, 55, 40, 35, 20, 10, 0, 10, 20, 30}},
		{Id: 7, Demand: 100, ServiceTime: 10, Distance: []int{70, 65, 50, 45, 30, 20, 10, 0, 10, 20}},
		{Id: 8, Demand: 100, ServiceTime: 10, Distance: []int{80, 75, 60, 55, 40, 30, 20, 10, 0, 10}},
		{Id: 9, Demand: 100, ServiceTime: 10, Distance: []int{90, 85, 70, 65, 50, 40, 30, 20, 10, 0}},
	}
	// Configuración de condiciones específicas para cada Vehiculo
	ants := []models.Ant{
		{Id: 1, AverageSpeed: 20, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 200, RemainingCapacity: 200, FixedCost: 10, VariableCost: 10},
		{Id: 2, AverageSpeed: 20, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 200, RemainingCapacity: 200, FixedCost: 10, VariableCost: 10},
		{Id: 3, AverageSpeed: 20, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 200, RemainingCapacity: 200, FixedCost: 10, VariableCost: 10},
	}
	// Print the properties of the nodes and the vehicles
	// core.Print(nodes, ants)

	for iteration := 0; iteration < NumIterations; iteration++ {
		fmt.Println("Iteration:", iteration+1, " ------------------------------------------------------------")
		var allResult [][]int

		pheromones := make([][]float64, NumNodes)
		for i := range pheromones {
			pheromones[i] = make([]float64, NumNodes)
			for j := range pheromones[i] {
				pheromones[i][j] = 1.0
			}
		}

		visitedNodes := make([]bool, NumNodes)

		for step := 0; step < NumNodes-1; step++ {
			for i := range ants {
				ant := &ants[i]

				if ant.RemainingCapacity <= 0 {
					ant.Route = append(ant.Route, startNode)
					continue
				}

				// Select the next node that the vehicle will visit
				ant.CurrentNode = antcolonyalgorithm.SelectNextNode(ant, pheromones, nodes, visitedNodes, NumNodes)

				if ant.CurrentNode == startNode {
					// [0]
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

		antcolonyalgorithm.UpdatePheromones(pheromones, ants, nodes, NumNodes)

		for i := range ants {
			ants[i].Route = append(ants[i].Route, startNode)
		}

		var cost_total_iteration float64
		var time_total_iteration float64

		// vehicleOrder := rand.Perm(len(ants))
		for _, ant := range ants {
			// ant := ants[k]
			route := core.RemoveDuplicateNodesInRoute(ant.Route)
			// Calculate distances and service times for each route
			distance, serviceTime := mathematicalmodel.SumDistanceAndServicesTime(route, nodes)
			// Calculation of the costs of each route
			cost := mathematicalmodel.CostTarget(distance, float64(ant.VariableCost), float64(ant.FixedCost))
			cost_total_iteration += float64(cost)

			// Calculation of total route time:
			totalTime := mathematicalmodel.TimeTarget(ant.AverageSpeed, float64(distance), float64(serviceTime))
			time_total_iteration += totalTime

			fmt.Println("Vehicle Route", ant.Id, "--ROUTE-->", route, "--distance-->", distance, "km --COST-->", cost, " --Total Time:", totalTime)
			// -----------
			allResult = append(allResult, route)
		}
		fmt.Println("-------------")
		fmt.Println("All Routes: ", allResult, " ---COST:", cost_total_iteration, "  ---Total Time:", time_total_iteration)

		// Generate graph for each iteration
		core.Grafo(nodes, ants, iteration)
	}
}
