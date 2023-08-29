package core

import (
	"fmt"
	"io/ioutil"
	"routng/models"
	"strconv"

	"github.com/awalterschulze/gographviz"
)

func Print(nodes []models.Node, ants []models.Ant) {
	fmt.Println("--------Nodos---------")
	for _, node := range nodes {
		fmt.Println(node.Id, " Demand:", node.Demand, " ServiceTime:", node.ServiceTime, " Distance", node.Distance, " ServiceTime:", node.ServiceTime)
	}

	fmt.Println("--------Vihicle---------")
	for _, ant := range ants {
		fmt.Println(ant.Id, " Capacity:", ant.Capacity, " FixedCost:", ant.FixedCost, " VariableCost:", ant.VariableCost, " AverageSpeed:", ant.AverageSpeed)
	}
	fmt.Println("------------------------------------------------------------")
}

// Remove duplicate nodes in a route
func RemoveDuplicateNodesInRoute(route []int) []int {
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

func Grafo(nodes []models.Node, ants []models.Ant, iteration int) {
	// Generar la representación gráfica en formato DOT
	graphAst := gographviz.NewGraph()
	graphAst.SetDir(true) // Para un grafo dirigido
	graphAst.SetName("G")
	currentGraph := make(models.Graph)
	for _, ant := range ants {
		route := ant.Route
		for i := 0; i < len(route)-1; i++ {
			from := route[i]
			to := route[i+1]
			if _, ok := currentGraph[from]; !ok {
				currentGraph[from] = make(map[int]int)
			}
			currentGraph[from][to] = nodes[to].Distance[from]
		}
	}
	// Agregar los nodos al grafo
	for node := range currentGraph {
		attrs := map[string]string{
			"label": fmt.Sprintf("%d", node),
		}
		graphAst.AddNode("G", fmt.Sprintf("%d", node), attrs)
	}

	// Agregar las aristas al grafo
	for from, connections := range currentGraph {
		for to, distance := range connections {
			attrs := map[string]string{
				"label": fmt.Sprintf("%d", distance),
			}
			if !(from == to) {
				graphAst.AddEdge(fmt.Sprintf("%d", from), fmt.Sprintf("%d", to), true, attrs)
			}
		}
	}

	dot := graphAst.String()

	// Guardar la representación en un archivo temporal
	dotFilename := "iteration" + strconv.Itoa(iteration) + ".dot"
	err := ioutil.WriteFile("./grafos/"+dotFilename, []byte(dot), 0644)
	if err != nil {
		fmt.Println("Error al guardar la representación DOT:", err)
		return
	}
}
