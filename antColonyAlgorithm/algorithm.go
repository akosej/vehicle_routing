package antcolonyalgorithm

import (
	"math/rand"
	"routng/models"
)

func UpdatePheromones(pheromones [][]float64, ants []models.Ant, nodes []models.Node, NumNodes int) {
	/* Esta función `UpdatePheromones` implementa la actualización de los niveles de feromonas en un algoritmo de colonia de hormigas.

	1. `evaporation := 0.5`: Se define la tasa de evaporación de las feromonas.
	Es un valor constante que indica cuánto se evaporan las feromonas existentes en cada iteración del algoritmo.

	2. El primer bucle `for` itera sobre todos los nodos en la matriz de feromonas `pheromones`.
	El número de nodos se especifica mediante el parámetro `NumNodes`.
	3. El segundo bucle `for` itera sobre todos los nodos nuevamente para actualizar las feromonas.
	Si `i` y `j` no son iguales, se multiplica el valor de feromona existente en `pheromones[i][j]` por la tasa de evaporación `evaporation`.
	Esto simula la evaporación de las feromonas en cada iteración.

	4. Después de la evaporación de las feromonas, se recorren todas las hormigas en el slice `ants` utilizando un bucle `for range`.

	5. Se obtiene la ruta de la hormiga actual `ant.Route`.

	6. Se itera sobre la ruta de la hormiga utilizando un bucle `for`.
	En cada iteración, se obtiene el nodo actual `from` y el siguiente nodo `to` en la ruta.

	7. Se incrementa el nivel de feromonas en la matriz `pheromones` en la posición `pheromones[from][to]`.
	La cantidad de incremento está determinada por la fórmula `1 / float64(nodes[to].Distance[from])`.
	Esto significa que cuanto más corta sea la distancia entre los nodos, mayor será el incremento de feromonas.

	En resumen, la función realiza la evaporación de las feromonas existentes y luego actualiza
	los niveles de feromonas en función de las rutas seguidas por las hormigas. Esto es esencial en
	el algoritmo de colonia de hormigas para que las hormigas
	*/
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

func SelectNextNode(ant *models.Ant, pheromones [][]float64, nodes []models.Node, visitedNodes []bool, NumNodes int) int {
	/* La función `selectNextNode` se encarga de seleccionar el siguiente nodo para una hormiga en el algoritmo de colonia de hormigas.

	1. La función toma como argumentos un puntero a una estructura `Ant`, una matriz de feromonas `pheromones`,
	un slice de nodos `nodes` y un slice de nodos visitados `visitedNodes`.

	2. Se declaran variables para almacenar el total de feromonas y el número de nodos disponibles.

	3. Se utiliza un bucle `for` para iterar sobre todos los nodos en la matriz de feromonas.
	Si el nodo no ha sido visitado por la hormiga actual y no ha sido visitado anteriormente,
	se calcula el total de feromonas sumando el producto de las feromonas en la posición `pheromones[ant.CurrentNode][i]`
	y la inversa de la distancia entre los nodos `1 / float64(nodes[i].Distance[ant.CurrentNode])`. Además, se incrementa el
	contador de nodos disponibles.

	4. Si no hay nodos disponibles para visitar, se devuelve el nodo actual de la hormiga.

	5. Se crean slices `probs` y `availableNodes` para almacenar las probabilidades de selección y los nodos disponibles respectivamente.
	Se itera nuevamente sobre todos los nodos para calcular las probabilidades de selección y almacenar los nodos disponibles en los
	slices correspondientes.

	6. Se crean slices `selectedNodes` y `selectedCapacity` para almacenar los nodos seleccionados y la capacidad seleccionada
	respectivamente. Se utiliza un bucle `for` para seleccionar nodos mientras la capacidad seleccionada sea menor que la capacidad
	restante de la hormiga.

	7. Se genera un número aleatorio `r` y se utiliza para seleccionar un nodo basado en las probabilidades de selección acumuladas.
	Si el nodo seleccionado no cumple con la capacidad restante de la hormiga, se rompe el bucle.

	8. Si se han seleccionado nodos, se actualiza la capacidad y el nodo actual de la hormiga. Si solo se selecciona un nodo,
	se devuelve ese nodo. De lo contrario, se devuelve el último nodo seleccionado.

	9. Si no se han seleccionado nodos, se devuelve el nodo actual de la hormiga.

	En resumen, la función `selectNextNode` selecciona el siguiente nodo para una hormiga en función de las feromonas,
	las distancias y las capacidades de los nodos. Esto se hace utilizando un mecanismo de selección basado en las probabilidades
	de las feromonas y la capacidad restante de la hormiga.
	*/
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

		// Comprobación de distancia mínima
		if nodes[selectedNode].Distance[ant.CurrentNode] > ant.RemainingCapacity {
			break
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
