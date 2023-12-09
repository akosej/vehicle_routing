package models

type Compartment struct {
	Id       int
	Capacity int
}

type Ant struct {
	Id                int
	Visited           []bool
	CurrentNode       int
	Capacity          []Compartment
	AverageSpeed      float64
	FixedCost         float64
	VariableCost      float64
	RemainingCapacity []Compartment
	Route             []int
}
