package antcolonyalgorithm

import (
	"fmt"
	"math/rand"
	"routng/models"
	"time"
)

func UpdatePheromones(pheromones [][]float64, ants []models.Ant, nodes []models.Node, NumNodes int) {
	// La función `UpdatePheromones` se encarga de actualizar los niveles de feromonas en el grafo del algoritmo de colonia de hormigas.

	evaporation := 0.1

	// Evaporar las feromonas existentes
	for i := 0; i < NumNodes; i++ {
		for j := 0; j < NumNodes; j++ {
			if i != j {
				// Se reduce el nivel de feromonas mediante la evaporación
				pheromones[i][j] *= evaporation
			}
		}
	}

	// Actualizar las feromonas en función de las rutas de las hormigas
	for _, ant := range ants {
		route := ant.Route
		for i := 0; i < len(route)-1; i++ {
			from := route[i]
			to := route[i+1]
			// Se incrementa el nivel de feromonas en la arista correspondiente al movimiento de la hormiga
			pheromones[from][to] += 1 / nodes[to].Distance[from]
		}
	}
}

func SelectNextNode(ant *models.Ant, pheromones [][]float64, nodes []models.Node, visitedNodes []bool, NumNodes int) int {
	// La función `SelectNextNode` se encarga de seleccionar el siguiente nodo para una hormiga en el algoritmo de colonia de hormigas.

	var totalPheromone float64
	var numAvailableNodes int

	// Calcular el total de feromonas y el número de nodos disponibles
	for i := 0; i < NumNodes; i++ {
		if !ant.Visited[i] && !visitedNodes[i] {
			// Calcular el total de feromonas acumulando las feromonas de los nodos disponibles
			totalPheromone += pheromones[ant.CurrentNode][i] * (1 / float64(nodes[i].Distance[ant.CurrentNode]))
			// Contar el número de nodos disponibles
			numAvailableNodes++
		}
	}

	// fmt.Println("Veh: ", ant.Id, " totalPheromone: ", totalPheromone, " numAvailableNodes: ", numAvailableNodes)

	if numAvailableNodes == 0 {
		// Si no hay nodos disponibles para visitar, se devuelve el nodo actual de la hormiga
		return ant.CurrentNode
	}

	// Crear slices para almacenar las probabilidades y los nodos disponibles
	probs := make([]float64, numAvailableNodes)
	availableNodes := make([]int, numAvailableNodes)
	index := 0
	for i := 0; i < NumNodes; i++ {
		if !ant.Visited[i] && !visitedNodes[i] {
			// Calcular la probabilidad de selección para cada nodo disponible
			probs[index] = (pheromones[ant.CurrentNode][i] * (1 / float64(nodes[i].Distance[ant.CurrentNode]))) / totalPheromone
			// Almacenar los nodos disponibles
			availableNodes[index] = i
			index++
		}
	}

	// Crear slices para almacenar los nodos seleccionados y la capacidad seleccionada
	selectedNodes := make([]int, 0)
	selectedCapacity := 0

	for selectedCapacity <= ant.RemainingCapacity[0].Capacity {
		// Generar un número aleatorio entre 0 y 1
		rand.Seed(time.Now().UnixNano())
		SumPro := 0.0
		for i := 0; i < numAvailableNodes; i++ {
			SumPro += probs[i]
		}

		// min := 0.0  // Minimum value
		// Generate a random float64 number within the specified range
		r := rand.Float64() * (SumPro)

		sum := 0.0
		selectedNode := -1
		for i := 0; i < numAvailableNodes; i++ {
			// Acumular las probabilidades de selección
			sum += probs[i]
			if r <= sum {
				// Seleccionar el nodo correspondiente al número aleatorio
				selectedNode = availableNodes[i]
				break
			}
		}

		if selectedNode == -1 {
			// Si no se seleccionó ningún nodo, seleccionar el un nodo aleatorio
			rand.Seed(time.Now().UnixNano())

			// Generar un número aleatorio en el rango del 1 al 9
			randomNumber := rand.Intn(numAvailableNodes)
			selectedNode = availableNodes[randomNumber]
		}

		// Comprobación de distancia mínima
		// if nodes[selectedNode].Distance[ant.CurrentNode] > ant.RemainingCapacity {
		// 	// Si la distancia entre el nodo seleccionado y el nodo actual es mayor que la capacidad restante, se interrumpe la selección
		// 	break
		// }

		// Si la capacidad seleccionada más la demanda del nodo seleccionado es menor o igual a la capacidad restante de la hormiga
		// y la demanda del nodo seleccionado es mayor que 0
		if selectedCapacity+nodes[selectedNode].Demand[0].Demand <= ant.RemainingCapacity[0].Capacity && nodes[selectedNode].Demand[0].Demand > 0 {
			// Se actualiza la capacidad seleccionada, se marca el nodo como visitado y se agrega el nodo a los nodos seleccionados
			selectedCapacity += nodes[selectedNode].Demand[0].Demand
			ant.Visited[selectedNode] = true
			selectedNodes = append(selectedNodes, selectedNode)
		}
		// Si no se cumple la condición anterior, se interrumpe la selección
		break
	}

	if len(selectedNodes) > 0 {

		ant.Capacity[0].Capacity -= nodes[selectedNodes[0]].Demand[0].Demand
		ant.CurrentNode = selectedNodes[0]
		return selectedNodes[0]
	}

	// Si no se seleccionó ningún nodo, se devuelve el nodo actual de la hormiga
	return ant.CurrentNode
}

func GenerateRoute(nodes []models.Node, ants []models.Ant, NumNodes, startNode int, pheromones [][]float64, remainingVisited []bool) {
	// Generar nuevas rutas con los nodos no visitados
	for i := range ants {
		ant := &ants[i]
		// TODO Hacer el ciclo para iterar sobre los compartimientos
		if ant.RemainingCapacity[0].Capacity <= 0 {
			// Si la capacidad restante de la hormiga es menor o igual a cero,
			// se agrega el nodo inicial a la ruta de la hormiga y se continúa con la siguiente iteración.
			ant.Route = append(ant.Route, startNode)
			continue
		}

		// Seleccionar el siguiente nodo que aún no se haya visitado
		totalDemand := 0
		for i, value := range remainingVisited {
			if !value {
				totalDemand += nodes[i].Demand[0].Demand
			}
		}
		//
		for ant.RemainingCapacity[0].Capacity > 0 && totalDemand > 0 {

			ant.CurrentNode = SelectNextNode(ant, pheromones, nodes, remainingVisited, NumNodes)
			if ant.CurrentNode == startNode {
				// Si el siguiente nodo seleccionado es el nodo inicial,
				// se agrega el nodo inicial a la ruta de la hormiga.
				ant.Route = append(ant.Route, startNode)
			} else {
				// Si el siguiente nodo seleccionado no es el nodo inicial,
				// se marca como visitado tanto en la hormiga como en la lista de nodos no visitados.
				if !remainingVisited[ant.CurrentNode] {
					ant.Visited[ant.CurrentNode] = true
					remainingVisited[ant.CurrentNode] = true

					demand := nodes[ant.CurrentNode].Demand[0].Demand
					product := nodes[ant.CurrentNode].Demand[0].Product
					fmt.Println(product)
					// fmt.Println(ant.RemainingCapacity, demand, nodes[ant.CurrentNode].Id, nodes[ant.CurrentNode].Demand)

					if ant.RemainingCapacity[0].Capacity >= demand {
						// Se actualiza la capacidad restante de la hormiga restando la demanda del nodo seleccionado.
						ant.Route = append(ant.Route, ant.CurrentNode)
						ant.RemainingCapacity[0].Capacity -= demand
						totalDemand -= demand
					}
					// TODO Aqui hay que ver lo de los remanentes para ver la capacidad que le queda al vehiculo
					if ant.RemainingCapacity[0].Capacity < demand {
						ant.RemainingCapacity[0].Capacity = 0
					}
				}
			}

		}

	}
}
