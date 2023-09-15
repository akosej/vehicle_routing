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
	NumIterations = 10
)

func main() {
	startNode := 0
	// Se define el nodo de inicio como 0

	rand.Seed(time.Now().UnixNano()) // Se establece la semilla para la generación de números aleatorios

	for iteration := 1; iteration <= NumIterations; iteration++ { // Se inicia un ciclo for para iterar un número determinado de veces
		nodes := []models.Node{ // Se crea un slice de estructuras Node
			{Id: 0, Demand: 0, ServiceTime: 30, Distance: []float64{0.00, 98.45, 36.92, 21.97, 80.18, 22.24, 99.45, 65.50, 50.97, 40.00, 41.00, 31.22, 28.44, 128.33, 142.66, 98.45, 49.23, 101.31, 77.10}},
			{Id: 1, Demand: 7182, ServiceTime: 30, Distance: []float64{98.45, 0.00, 128.82, 86.06, 177.07, 93.83, 197.41, 162.62, 141.58, 138.02, 124.61, 84.06, 106.46, 226.02, 240.30, 196.91, 49.23, 198.35, 174.11}},
			{Id: 2, Demand: 7762, ServiceTime: 30, Distance: []float64{36.92, 128.82, 0.00, 42.77, 66.69, 35.27, 79.92, 53.85, 14.06, 20.28, 13.26, 67.65, 57.27, 98.50, 112.52, 74.28, 80.97, 70.32, 63.68}},
			{Id: 3, Demand: 7222, ServiceTime: 30, Distance: []float64{21.97, 86.06, 42.77, 0.00, 98.99, 8.20, 116.69, 84.40, 55.64, 54.35, 39.63, 42.41, 48.58, 140.70, 154.87, 113.77, 39.00, 112.74, 95.85}},
			{Id: 4, Demand: 8996, ServiceTime: 30, Distance: []float64{80.18, 177.07, 66.69, 98.99, 0.00, 94.35, 22.24, 14.68, 65.57, 47.26, 79.56, 95.74, 73.47, 64.84, 77.67, 29.83, 128.33, 49.70, 3.15}},
			{Id: 5, Demand: 6464, ServiceTime: 30, Distance: []float64{22.24, 93.83, 35.27, 8.20, 94.35, 0.00, 111.20, 79.92, 47.76, 48.50, 31.45, 47.76, 50.49, 133.61, 147.71, 107.58, 47.17, 105.53, 91.21}},
			{Id: 6, Demand: 7636, ServiceTime: 30, Distance: []float64{99.45, 197.41, 79.92, 116.69, 22.24, 111.20, 0.00, 35.27, 75.12, 62.73, 91.69, 117.27, 95.27, 45.04, 56.77, 14.06, 148.35, 37.97, 24.56}},
			{Id: 7, Demand: 7902, ServiceTime: 30, Distance: []float64{65.50, 162.62, 53.85, 84.40, 14.68, 79.92, 35.27, 0.00, 54.73, 33.81, 66.99, 82.01, 60.14, 73.81, 87.45, 39.00, 113.77, 53.85, 11.60}},
			{Id: 8, Demand: 6359, ServiceTime: 30, Distance: []float64{50.97, 141.58, 14.06, 55.64, 65.57, 47.76, 75.12, 54.73, 0.00, 25.32, 18.50, 81.55, 69.96, 87.92, 101.62, 67.36, 94.33, 59.68, 62.80}},
			{Id: 9, Demand: 5530, ServiceTime: 30, Distance: []float64{40.00, 138.02, 20.28, 54.35, 47.26, 48.50, 62.73, 33.81, 25.32, 0.00, 33.54, 66.12, 49.73, 88.33, 102.66, 59.49, 88.89, 61.43, 44.16}},
			{Id: 10, Demand: 7774, ServiceTime: 30, Distance: []float64{41.00, 124.61, 13.26, 39.63, 79.56, 31.45, 91.69, 66.99, 18.50, 33.54, 0.00, 72.17, 65.50, 106.41, 120.08, 84.98, 78.62, 78.17, 76.58}},
			{Id: 11, Demand: 4560, ServiceTime: 30, Distance: []float64{31.22, 84.06, 67.65, 42.41, 95.74, 47.76, 117.27, 82.01, 81.55, 66.12, 72.17, 0.00, 22.67, 152.12, 166.39, 119.45, 39.99, 126.69, 93.00}},
			{Id: 12, Demand: 5317, ServiceTime: 30, Distance: []float64{28.44, 106.46, 57.27, 48.58, 73.47, 50.49, 95.27, 60.14, 69.96, 49.73, 65.50, 22.67, 0.00, 131.90, 146.02, 98.34, 60.39, 107.65, 70.83}},
			{Id: 13, Demand: 5053, ServiceTime: 30, Distance: []float64{128.33, 226.02, 98.50, 140.70, 64.84, 133.61, 45.04, 73.81, 87.92, 88.33, 106.41, 152.12, 131.90, 0.00, 14.35, 35.16, 177.07, 28.27, 66.12}},
			{Id: 14, Demand: 7474, ServiceTime: 30, Distance: []float64{142.66, 240.30, 112.52, 154.87, 77.67, 147.71, 56.77, 87.45, 101.62, 102.66, 120.08, 166.39, 146.02, 14.35, 0.00, 48.50, 191.38, 42.19, 79.18}},
			{Id: 15, Demand: 5668, ServiceTime: 30, Distance: []float64{98.45, 196.91, 74.28, 113.77, 29.83, 107.58, 14.06, 39.00, 67.36, 59.49, 84.98, 119.45, 98.34, 35.16, 48.50, 0.00, 147.68, 23.91, 30.97}},
			{Id: 16, Demand: 4322, ServiceTime: 30, Distance: []float64{49.23, 49.23, 80.97, 39.00, 128.33, 47.17, 148.35, 113.77, 94.33, 88.89, 78.62, 39.99, 60.39, 177.07, 191.38, 147.68, 0.00, 149.61, 125.33}},
			{Id: 17, Demand: 6312, ServiceTime: 30, Distance: []float64{101.31, 198.35, 70.32, 112.74, 49.70, 105.53, 37.97, 53.85, 59.68, 61.43, 78.17, 126.69, 107.65, 28.27, 42.19, 23.91, 149.61, 0.00, 49.73}},
			{Id: 18, Demand: 8606, ServiceTime: 30, Distance: []float64{77.10, 174.11, 63.68, 95.85, 3.15, 91.21, 24.56, 11.60, 62.80, 44.16, 76.58, 93.00, 70.83, 66.12, 79.18, 30.97, 125.33, 49.73, 0.00}},
		}
		NumNodes := 18
		ants := []models.Ant{ // Se crea un slice de estructuras Ant
			{Id: 1, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 19990, RemainingCapacity: 19990, FixedCost: 197.26, VariableCost: 19.72},
			{Id: 2, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 23894, RemainingCapacity: 23894, FixedCost: 197.26, VariableCost: 19.29},
			{Id: 3, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 23609, RemainingCapacity: 23609, FixedCost: 197.26, VariableCost: 19.29},
			{Id: 4, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 11070, RemainingCapacity: 11070, FixedCost: 197.26, VariableCost: 22.72},
			{Id: 5, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 11066, RemainingCapacity: 11066, FixedCost: 197.26, VariableCost: 22.72},
		}
		// ----------------------
		// if iteration == 1 {
		// 	fmt.Println(nodes, ants)
		// }
		fmt.Println(" ------------------------------------------------------------Iteration:", iteration, " ------------------------------------------------------------")
		pheromones := make([][]float64, NumNodes) // Se crea una matriz para almacenar los niveles de feromonas
		for i := range pheromones {               // Se itera sobre las filas de la matriz
			pheromones[i] = make([]float64, NumNodes) // Se crea una fila con la longitud adecuada
			for j := range pheromones[i] {            // Se itera sobre las columnas de la fila
				pheromones[i][j] = 1.0 // Se inicializan los niveles de feromonas en 1.0
			}
		}
		var allResult [][]int            // Se declara una variable para almacenar los resultados de todas las iteraciones
		var cost_total_iteration float64 // Se declara una variable para almacenar el costo total de la iteración actual
		var time_total_iteration float64 // Se declara una variable para almacenar el tiempo total de la iteración actual

		containsFalse := true // Se inicializa una variable para verificar si hay nodos no visitados

		remainingVisited := make([]bool, NumNodes) // Se crea un slice para almacenar los nodos no visitados
		for containsFalse {
			vehicles := []models.Ant{ // Se crea un slice de estructuras Ant
				{Id: 1, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 19990, RemainingCapacity: 19990, FixedCost: 197.26, VariableCost: 19.72},
				{Id: 2, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 23894, RemainingCapacity: 23894, FixedCost: 197.26, VariableCost: 19.29},
				{Id: 3, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 23609, RemainingCapacity: 23609, FixedCost: 197.26, VariableCost: 19.29},
				{Id: 4, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 11070, RemainingCapacity: 11070, FixedCost: 197.26, VariableCost: 22.72},
				{Id: 5, AverageSpeed: 50, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 11066, RemainingCapacity: 11066, FixedCost: 197.26, VariableCost: 22.72},
			}
			// ---
			remainingNodes := []int{}                  // Se crea un slice para almacenar los nodos no visitados
			for i, visited := range remainingVisited { // Se itera sobre los nodos visitados
				if visited { // Si el nodo fue visitado
					remainingNodes = append(remainingNodes, i) // Se agrega el nodo al slice de nodos visitados
				}
			}

			for _, node := range remainingNodes { // Se itera sobre los nodos no visitados
				remainingVisited[node] = true // Se marca el nodo como visitado en remainingVisited
			}
			// fmt.Println(remainingNodes)
			// Se llama a la función GenerateRoute del paquete antcolonyalgorithm para generar las rutas de las vehiculoss para los nodos no visitados
			antcolonyalgorithm.GenerateRoute(nodes, vehicles, NumNodes, startNode, pheromones, remainingVisited)
			// Se llama a la función UpdatePheromones del paquete antcolonyalgorithm para actualizar los niveles de feromonas
			antcolonyalgorithm.UpdatePheromones(pheromones, vehicles, nodes, NumNodes)

			for i := range vehicles { // Se itera sobre las vehiculoss para los nodos no visitados
				ants[i].Route = append(ants[i].Route, vehicles[i].Route...) // Se agrega la ruta de cada vehiculos al slice de rutas de todas las vehiculoss
			}

			result, cost_total, time_total := TransportationPeriod(vehicles, nodes) // Se llama a la función TransportationPeriod para realizar el período de transporte para los nodos no visitados

			allResult = append(allResult, result...) // Se agrega el resultado de los nodos no visitados al slice allResult
			cost_total_iteration += cost_total       // Se suma el costo total de los nodos no visitados al costo total acumulado
			time_total_iteration += time_total       // Se suma el tiempo total de los nodos no visitados al tiempo total acumulado

			// Comprobar si se visitaron todos los nodos
			containsFalse = false
			for i, visited := range remainingVisited { // Se itera sobre los nodos visitados
				if !visited && i > 0 { // Si hay un nodo no visitado (excepto el nodo de inicio)
					containsFalse = true // Se establece la variable containsFalse en true
					break                // Se sale del ciclo for
				}
			}

		}

		fmt.Println("-------------")
		fmt.Println("All Routes: ", allResult, " ---COST:", cost_total_iteration, "  ---Total Time:", time_total_iteration) // Se imprime

		core.Grafo(nodes, ants, iteration)
	}
}

func TransportationPeriod(ants []models.Ant, nodes []models.Node) ([][]int, float64, float64) {
	var result [][]int
	var cost_total float64
	var time_total float64

	for _, ant := range ants {
		// Elimina nodos duplicados en la ruta
		route := core.RemoveDuplicateNodesInRoute(ant.Route)
		// Calcula la distancia y el tiempo de servicio de la ruta
		distance, serviceTime := mathematicalmodel.SumDistanceAndServicesTime(route, nodes)
		// Calcula el costo de la ruta
		cost := mathematicalmodel.CostTarget(distance, float64(ant.VariableCost), float64(ant.FixedCost))
		// Calcula el tiempo total de la ruta
		totalTime := mathematicalmodel.TimeTarget(ant.AverageSpeed, float64(distance), float64(serviceTime))

		cost_total += cost      // Suma el costo al costo total
		time_total += totalTime // Suma el tiempo total al tiempo total

		// Si se crea una ruta para el vehiculo
		if len(route) > 2 {
			fmt.Println("Vehicle:", ant.Id, "--ROUTE-->", route, "--distance-->", distance, "km --COST-->", cost, " --Total Time:", totalTime, " ---RemainingCapacity--", ant.RemainingCapacity, " ---Capacity--", ant.Capacity)
			result = append(result, route) // Agrega la ruta al resultado
		}
	}

	return result, cost_total, time_total
}
