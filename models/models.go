package models

type Graph map[int]map[int]int

type Node struct {
	Id          int
	Demand      []int
	ServiceTime int
	Distance    []float64
}

type Route []int
type Visited []bool

type Ant struct {
	Id                int
	Visited           []Visited
	CurrentNode       []int
	Capacity          []int
	RemainingCapacity []int
	AverageSpeed      float64
	FixedCost         float64
	VariableCost      float64
	Route             []Route
}
