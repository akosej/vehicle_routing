package models

type Graph map[int]map[int]int

type Node struct {
	Id          int
	Demand      int
	ServiceTime int
	Distance    []int
}

type Ant struct {
	Id                int
	Visited           []bool
	CurrentNode       int
	Capacity          int
	AverageSpeed      float64
	FixedCost         float64
	VariableCost      float64
	RemainingCapacity int
	Route             []int
}
