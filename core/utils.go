package core

import (
	"fmt"
	"math/rand"
	"os"
	"routng/models"
	"strconv"
	"strings"
	"time"
)

func Print(nodes []models.Node, ants []models.Ant) {
	fmt.Println("--------Nodos---------")
	for _, node := range nodes {
		fmt.Println(node.Id, " Demand:", node.Demand, " ServiceTime:", node.ServiceTime, " Distance", node.Distance, " ServiceTime:", node.ServiceTime)
	}

	fmt.Println("--------Vihicle---------")
	for _, ant := range ants {
		fmt.Println(ant.Id, " AverageSpeed:", ant.AverageSpeed, " Capacity:", ant.Capacity, " FixedCost:", ant.FixedCost, " VariableCost:", ant.VariableCost)
	}
}

// Remove duplicate nodes in a route
func RemoveDuplicateNodesInRoute(route []int) []int {
	// La función `RemoveDuplicateNodesInRoute` se encarga de eliminar los nodos duplicados en una ruta.

	encountered := map[int]bool{} // Mapa para almacenar los nodos encontrados
	result := []int{}             // Slice para almacenar los nodos únicos

	for _, val := range route {
		if !encountered[val] { // Si el nodo no ha sido encontrado anteriormente
			encountered[val] = true      // Agregar el nodo al mapa como encontrado
			result = append(result, val) // Agregar el nodo al slice de resultados
		}
	}

	result = append(result, 0) // Agregar el nodo de origen al final de la ruta

	return result // Devolver la ruta sin nodos duplicados
}

func Grafo(routes [][]int, iteration int) {

	dotFilename := "iteration" + strconv.Itoa(iteration) + ".dot"
	f, err := os.Create("./grafos/" + dotFilename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString("digraph G {\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	lineColor := "blue"
	if iteration == 1 {
		lineColor = "red"
	}

	for k, route := range routes {
		rand.Seed(time.Now().UnixNano())
		// color := fmt.Sprintf("#%06x", rand.Intn(0xffffff))
		// time.Sleep(100 * time.Millisecond)
		routeString := intArrayToString(route)
		routeString = strings.Trim(fmt.Sprintf("%v", routeString), "[")
		routeString = strings.Trim(fmt.Sprintf("%v", routeString), "]")

		f.WriteString(fmt.Sprintf("\tbeautify=true\n"))
		f.WriteString(fmt.Sprintf("\tcenter=true\n"))
		f.WriteString(fmt.Sprintf("\tconcentrate=true\n"))
		f.WriteString(fmt.Sprintf("\tnode [colorscheme=oranges9] \n"))

		f.WriteString(fmt.Sprintf("\tnode [style=filled, color=%d] %s;\n", k+1, routeString))

		f.WriteString(fmt.Sprintf("\tsubgraph clusterG%d {\n", k))
		f.WriteString(fmt.Sprintf("\tlabel=\"V%d\" \n", k+1))
		f.WriteString(fmt.Sprintf("\tcolor=%d \n", k+1))
		f.WriteString(fmt.Sprintf("\tbgcolor=white\n"))
		routeString = strings.Trim(fmt.Sprintf("%v", routeString), "0")
		f.WriteString(fmt.Sprintf("\t%s", routeString))
		f.WriteString(fmt.Sprintf("}\n\n"))

		for i := 0; i < len(route)-1; i++ {
			from := route[i]
			to := route[i+1]

			_, err := f.WriteString(fmt.Sprintf("\t%d -> %d [label=\"P#%d\"; color=%s];\n", from, to, iteration+1, lineColor))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}

	_, err = f.WriteString("}\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Graph created successfully.")
}

func intArrayToString(arr []int) string {
	strArr := make([]string, len(arr))
	for i, num := range arr {
		strArr[i] = fmt.Sprintf("%d", num)
	}
	return fmt.Sprintf("%s", strArr)
}
