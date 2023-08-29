package mathematicalmodel

import "routng/models"

// Sum of distances and service times
func SumDistanceAndServicesTime(route []int, nodes []models.Node) (float64, float64) {

	var distance, serviceTime float64
	if len(route) == 2 {
		return 0, 0
	}
	for i := 0; i < len(route)-1; i++ {
		from := route[i]
		to := route[i+1]
		distance += float64(nodes[from].Distance[to])
		// Sum of the ServicesTime except the return to the deposit
		serviceTime += float64(nodes[i].ServiceTime)
	}
	return distance, serviceTime
}

// Calculation of the costs of each route
// Tienes el código de python? En la función objetivo de costo:
// Costo total = costos fijos + costos variables
// El costo fijo de cada vehículo es un valor fijo se se utiliza el vehículo, si no se utiliza 0
// El costo variable depende de la distancia recorrida. O sea, costo variable por unidad de distancia * distancia recorrida
func CostTarget(distance, variable_cost, fixed_cost float64) float64 {
	return (distance * variable_cost) + fixed_cost
}

// En la función objetivo de tiempo:
// Tiempo total = tiempo de viaje + tiempo de servicio
// El tiempo de viaje es igual al tiempo por unidad de distancia * distancia. El tiempo de servicio es un valor fijo en cada
// nodo que establece el tiempo que se tarda cualquier vehículo en serviciar al cliente
func TimeTarget(averageSpeed, distance, serviceTime float64) float64 {
	return ((1 / averageSpeed) * distance) + serviceTime
}
