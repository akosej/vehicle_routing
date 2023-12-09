package models

type Graph map[int]map[int]int

type Demand struct {
	Id      int
	Product int
	Demand  int
}
type Node struct {
	Id          int
	Demand      []Demand
	ServiceTime int
	Distance    []float64
}
