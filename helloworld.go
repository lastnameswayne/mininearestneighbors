package main

import (
	"fmt"
	"math"
	"math/rand"
	"rand"
)

type Vector struct {
	// A vector is a list of integers
	// In mesure case we also have an id
	id     ID
	vector []int
}
type ID int

type Vertex struct {
	id        ID
	neighbors []Vertex
}

type Graph struct {
	vertices []Vertex
}

func main() {

	fmt.Println("hello world")

	// Graph construction is first step

	//
}

func insertVector(vector Vector, efConstruction int) {
	M := 2 // number of neighbors to add to each vertex on insertion
	M_max := 4
	efConstruction := 100 // size of the dynamic list for the nearest neighbors

	levelMultiplier := 0 // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	nearestElements := []Vector{}

	level := calculateLevel(vector, levelMultiplier)

}

func calculateLevel(vector Vector, levelMultiplier float64) int {
	uniform := rand.Float64()
	prob := math.Log(-uniform * levelMultiplier)
	level := math.Floor(prob)

	return int(level)
}
