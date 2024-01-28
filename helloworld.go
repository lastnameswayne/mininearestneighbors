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
	level     int
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

	entryPointHNSW := getEntryPoint() //get enter point for hnsw
	entryPointLevel := entryPoint.level

	levelMultiplier := 0.0 // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	nearestElements := []Vector{}

	level := calculateLevel(vector, levelMultiplier)

	for i := 0; i < level; i++ {
		W := searchLayer(vector, entryPoint, 1, level)
		entryPoint := W[0]
	}
	for i := min(level, entryPoint.level); i > 0; i-- {
		W = searchLayer(vector, entryPoint, efConstruction, i)
		neighbors := selectNeighbors(vector, W, M, level)

		//add birectional connections from neighbors to q at layer l_c
		for _, n := range neighbors {
			neighbors := n.neighbors
			if len(neighbors) > M_max {
				newNeighbors := selectNeighbors(n, neighbors, M_max, level)
				n.neighbors = newNeighbors
			}
		}
		entryPoint = W
	}
	if level > entryPointLevel {
		entryPointHNSW = vector
	}

}

func getEntryPoint(vertex Vertex) Vertex {
	return vertex
}

func searchLayer(vertex Vertex, entryPoint Vertex, ef int, level int) []Vertex {

}

func selectNeighbors(vertex Vertex, W []Vertex, M int, level int) []Vertex {

}

func calculateLevel(vector Vector, levelMultiplier float64) int {
	//floor of -ln(unif(0,1)*mL)
	uniform := rand.Float64()
	prob := math.Log(-uniform * levelMultiplier)
	level := math.Floor(prob)

	return int(level)
}
