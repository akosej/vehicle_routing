package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	antcolonyalgorithm "routng/antColonyAlgorithm"
	"routng/core"
	mathematicalmodel "routng/mathematicalModel"
	"routng/models"
	"strconv"
	"strings"
	"time"
)

const (
	NumIterations = 2
)

func main() {
	startNode := 0
	// Se define el nodo de inicio como 0

	rand.Seed(time.Now().UnixNano()) // Se establece la semilla para la generación de números aleatorios

	for iteration := 1; iteration <= NumIterations; iteration++ { // Se inicia un ciclo for para iterar un número determinado de veces
		// nodes := []models.Node{ // Se crea un slice de estructuras Node
		// 	{Id: 0, Demand: 0, ServiceTime: 10, Distance: []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90}},
		// 	{Id: 1, Demand: 100, ServiceTime: 10, Distance: []int{10, 0, 15, 25, 35, 45, 55, 65, 75, 85}},
		// 	{Id: 2, Demand: 100, ServiceTime: 10, Distance: []int{20, 15, 0, 10, 20, 30, 40, 50, 60, 70}},
		// 	{Id: 3, Demand: 100, ServiceTime: 10, Distance: []int{30, 25, 10, 0, 15, 25, 35, 45, 55, 65}},
		// 	{Id: 4, Demand: 100, ServiceTime: 10, Distance: []int{40, 35, 20, 15, 0, 10, 20, 30, 40, 50}},
		// 	{Id: 5, Demand: 100, ServiceTime: 10, Distance: []int{50, 45, 30, 25, 10, 0, 10, 20, 30, 40}},
		// 	{Id: 6, Demand: 100, ServiceTime: 10, Distance: []int{60, 55, 40, 35, 20, 10, 0, 10, 20, 30}},
		// 	{Id: 7, Demand: 100, ServiceTime: 10, Distance: []int{70, 65, 50, 45, 30, 20, 10, 0, 10, 20}},
		// 	{Id: 8, Demand: 100, ServiceTime: 10, Distance: []int{80, 75, 60, 55, 40, 30, 20, 10, 0, 10}},
		// 	{Id: 9, Demand: 100, ServiceTime: 10, Distance: []int{90, 85, 70, 65, 50, 40, 30, 20, 10, 0}},
		// }
		NumNodes := 0
		file2, err := os.Open("nodes.txt")
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
			return
		}
		defer file2.Close()
		scanner1 := bufio.NewScanner(file2)
		var nodes []models.Node

		for scanner1.Scan() {
			line := scanner1.Text()
			if strings.Contains(line, ",") {
				NumNodes++
				values := strings.Split(line, ",")

				id, _ := strconv.Atoi(values[0])
				demand, _ := strconv.Atoi(values[1])
				serviceTime, _ := strconv.Atoi(values[2])

				distances := make([]int, len(values)-3)
				for i := 3; i < len(values); i++ {
					distance, _ := strconv.Atoi(values[i])
					distances[i-3] = distance
				}

				nodes = append(nodes, models.Node{
					Id:          id,
					Demand:      demand,
					ServiceTime: serviceTime,
					Distance:    distances,
				})
			}
		}

		// ants := []models.Ant{ // Se crea un slice de estructuras Ant
		// 	{Id: 1, Route: []int{startNode}},
		// 	{Id: 2, Route: []int{startNode}},
		// 	{Id: 3, Route: []int{startNode}},
		// }
		file, err := os.Open("vehicles.txt")
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var ants []models.Ant

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, ",") {
				values := strings.Split(line, ",")

				id, _ := strconv.Atoi(values[0])
				averageSpeed, _ := strconv.Atoi(values[1])
				capacity, _ := strconv.Atoi(values[2])
				fixedCost, _ := strconv.Atoi(values[3])
				variableCost, _ := strconv.Atoi(values[4])

				ants = append(ants, models.Ant{
					Id:                id,
					AverageSpeed:      float64(averageSpeed),
					Capacity:          capacity,
					RemainingCapacity: capacity,
					FixedCost:         float64(fixedCost),
					VariableCost:      float64(variableCost),
					Visited:           make([]bool, NumNodes),
					CurrentNode:       startNode,
					Route:             []int{startNode},
					// Resto de los campos
				})
			}
		}
		// ----------------------
		core.Print(nodes, ants)
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
			// vehicles := []models.Ant{ // Se crea un nuevo slice de estructuras Ant para atender los nodos no visitados
			// 	{Id: 1, AverageSpeed: 20, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 500, RemainingCapacity: 500, FixedCost: 10, VariableCost: 10},
			// 	{Id: 2, AverageSpeed: 20, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 300, RemainingCapacity: 300, FixedCost: 10, VariableCost: 10},
			// 	{Id: 3, AverageSpeed: 20, Visited: make([]bool, NumNodes), CurrentNode: startNode, Route: []int{startNode}, Capacity: 300, RemainingCapacity: 300, FixedCost: 10, VariableCost: 10},
			// }
			file, err := os.Open("vehicles.txt")
			if err != nil {
				fmt.Println("Error al abrir el archivo:", err)
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			var vehicles []models.Ant

			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, ",") {
					values := strings.Split(line, ",")

					id, _ := strconv.Atoi(values[0])
					averageSpeed, _ := strconv.Atoi(values[1])
					capacity, _ := strconv.Atoi(values[2])
					fixedCost, _ := strconv.Atoi(values[3])
					variableCost, _ := strconv.Atoi(values[4])

					vehicles = append(vehicles, models.Ant{
						Id:                id,
						AverageSpeed:      float64(averageSpeed),
						Capacity:          capacity,
						RemainingCapacity: capacity,
						FixedCost:         float64(fixedCost),
						VariableCost:      float64(variableCost),
						Visited:           make([]bool, NumNodes),
						CurrentNode:       startNode,
						Route:             []int{startNode},
						// Resto de los campos
					})
				}
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
