package antcolonyalgorithm

import (
	"math/rand"
	"routng/models"
)

func UpdatePheromones(pheromones [][]float64, ants []models.Ant, nodes []models.Node, NumNodes int) {
	// La función `UpdatePheromones` se encarga de actualizar los niveles de feromonas en el grafo del algoritmo de colonia de hormigas.

	evaporation := 0.5

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
			pheromones[from][to] += 1 / float64(nodes[to].Distance[from])
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
	for selectedCapacity < ant.RemainingCapacity {
		// Generar un número aleatorio entre 0 y 1
		r := rand.Float64()
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
			// Si no se seleccionó ningún nodo, seleccionar el último nodo disponible
			selectedNode = availableNodes[numAvailableNodes-1]
		}

		// Comprobación de distancia mínima
		if nodes[selectedNode].Distance[ant.CurrentNode] > ant.RemainingCapacity {
			// Si la distancia entre el nodo seleccionado y el nodo actual es mayor que la capacidad restante, se interrumpe la selección
			break
		}

		if selectedCapacity+nodes[selectedNode].Demand <= ant.RemainingCapacity && nodes[selectedNode].Demand > 0 {
			// Si la capacidad seleccionada más la demanda del nodo seleccionado es menor o igual a la capacidad restante de la hormiga y la demanda del nodo seleccionado es mayor que 0
			// Se actualiza la capacidad seleccionada, se marca el nodo como visitado y se agrega el nodo a los nodos seleccionados
			selectedCapacity += nodes[selectedNode].Demand
			ant.Visited[selectedNode] = true
			selectedNodes = append(selectedNodes, selectedNode)
		} else {
			// Si no se cumple la condición anterior, se interrumpe la selección
			break
		}
	}

	if len(selectedNodes) > 0 {
		// Si se seleccionaron nodos
		for _, node := range selectedNodes {
			// Se actualiza la capacidad y el nodo actual de la hormiga según los nodos seleccionados
			ant.Capacity -= nodes[node].Demand
			ant.CurrentNode = node
		}
		if len(selectedNodes) == 1 {
			// Si solo se seleccionó un nodo, se devuelve ese nodo
			return selectedNodes[0]
		}
		// Si se seleccionaron varios nodos, se devuelve el último nodo seleccionado
		return selectedNodes[len(selectedNodes)-1]
	}

	// Si no se seleccionó ningún nodo, se devuelve el nodo actual de la hormiga
	return ant.CurrentNode
}

func GenerateRoute(nodes []models.Node, ants []models.Ant, NumNodes, startNode int, pheromones [][]float64, remainingVisited []bool) {
	// Generar nuevas rutas con los nodos no visitados
	for step := 0; step < NumNodes-1; step++ {
		for i := range ants {
			ant := &ants[i]

			if ant.RemainingCapacity <= 0 {
				// Si la capacidad restante de la hormiga es menor o igual a cero,
				// se agrega el nodo inicial a la ruta de la hormiga y se continúa con la siguiente iteración.
				ant.Route = append(ant.Route, startNode)
				continue
			}

			// Seleccionar el siguiente nodo que aún no se haya visitado
			ant.CurrentNode = SelectNextNode(ant, pheromones, nodes, remainingVisited, NumNodes)

			if ant.CurrentNode == startNode {
				// Si el siguiente nodo seleccionado es el nodo inicial,
				// se agrega el nodo inicial a la ruta de la hormiga.
				ant.Route = append(ant.Route, startNode)
			} else {
				// Si el siguiente nodo seleccionado no es el nodo inicial,
				// se marca como visitado tanto en la hormiga como en la lista de nodos no visitados.
				ant.Visited[ant.CurrentNode] = true
				remainingVisited[ant.CurrentNode] = true

				demand := nodes[ant.CurrentNode].Demand
				if ant.RemainingCapacity >= demand {
					// Se actualiza la capacidad restante de la hormiga restando la demanda del nodo seleccionado.
					ant.Route = append(ant.Route, ant.CurrentNode)
					ant.RemainingCapacity -= demand
				}
			}
		}
	}
}
