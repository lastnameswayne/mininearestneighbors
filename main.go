package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
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
	vector    []int
}

type Graph struct {
	vertices []Vertex
}

func main() {

	fmt.Println("hello world")

	// Graph construction is first step
	v1 := Vector{
		id:     1,
		vector: []int{1, 2, 3, 4, 5},
	}
	v2 := Vector{
		id:     2,
		vector: []int{2, 2, 3, 5, 5},
	}
	v3 := Vector{
		id:     3,
		vector: []int{11, 12, 13, 14, 15},
	}
	v4 := Vector{
		id:     4,
		vector: []int{1, 1, 1, 1, 1},
	}
	v5 := Vector{
		id:     5,
		vector: []int{1, 12, 3, 4, 5},
	}
	v6 := Vector{
		id:     6,
		vector: []int{1, 2, 30, 4, 5},
	}
	v7 := Vector{
		id:     7,
		vector: []int{2, 2, 3, 4, 5},
	}
	v8 := Vector{
		id:     8,
		vector: []int{10, 200, 3, 4, 5},
	}

	q := Vector{
		id:     9,
		vector: []int{0, 2, 3, 4, 5},
	}
	Graph := Graph{
		vertices: []Vertex{},
	}

	Graph = insertVector(Graph, v1, 5)
	Graph = insertVector(Graph, v2, 5)
	Graph = insertVector(Graph, v3, 5)
	Graph = insertVector(Graph, v4, 5)
	Graph = insertVector(Graph, v5, 5)
	Graph = insertVector(Graph, v6, 5)
	Graph = insertVector(Graph, v7, 5)
	Graph = insertVector(Graph, v8, 5)

	fmt.Println(Graph)
	fmt.Println("searching")
	fmt.Println(Search(q, Graph, 5))
}

func Search(q Vector, graph Graph, efSize int, k int) []Vertex {
	//searches for q in graph
	//returns the closest vector to q
	//returns the vector and its distance to q
	//search starts at the top layer and traverses greedily to find the efSize closest neighbors to q
	//these are used as enterPoints in the next step
	//s
	W := []Vertex{}
	queryElement := Vertex{
		vector:    q.vector,
		id:        q.id,
		neighbors: []Vertex{},
	}
	enterPoint := getEnterPoint(queryElement, graph)
	enterPointLevel := enterPoint.level
	for i := enterPointLevel; i > 0; i-- {
		W = searchLevel(queryElement, []Vertex{enterPoint}, efSize, i)
		enterPoint = getClosestArr(queryElement, W)[0]
	}
	W = searchLevel(queryElement, []Vertex{enterPoint}, efSize, 0)
	return W[:k]
}

func insertVector(graph Graph, vector Vector, efSize int) Graph {
	vertex := Vertex{
		vector:    vector.vector,
		id:        vector.id,
		neighbors: []Vertex{},
	}
	M := 2 // number of neighbors to add to each vertex on insertion
	M_max := 4

	enterPointHNSW := getEnterPoint(vertex, graph) //get enter point for hnsw
	enterPointLevel := enterPointHNSW.level

	enterPoint := []Vertex{enterPointHNSW}

	levelMultiplier := 1 / math.Log(float64(M)) // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	nearestElements := []Vertex{}

	level := calculateLevel(levelMultiplier)
	vertex.level = level

	//Start at the top level and traverse greedilty to find the epSize closest neighbors to vector
	//These are used as enterPoints in the next step
	for i := enterPointLevel; i > vertex.level; i-- {
		nearestElements = searchLevel(vertex, enterPoint, 1, level)
		enterPoint = getClosestArr(vertex, nearestElements)
	}

	//searches again from the next layer
	for i := min(level, enterPoint[0].level); i > 0; i-- {
		nearestElements = searchLevel(vertex, enterPoint, efSize, i)
		neighbors := selectNeighbors(vertex, nearestElements, M, level)

		//add birectional connections from neighbors to q at layer l_c
		for _, n := range neighbors {
			neighbors := n.neighbors
			if len(neighbors) > M_max {
				newNeighbors := selectNeighbors(n, neighbors, M_max, level)
				n.neighbors = newNeighbors
				vertex.neighbors = newNeighbors
			}
		}
		enterPoint = nearestElements
	}
	if level > enterPointLevel {
		enterPointHNSW = vertex
	}

	graph.vertices = append(graph.vertices, vertex)
	return graph
}

// this can be random
func getEnterPoint(vertex Vertex, graph Graph) Vertex {
	if len(graph.vertices) == 0 {
		return vertex
	}
	randomIndex := math.Floor(rand.ExpFloat64() * float64(len(graph.vertices)))

	return graph.vertices[int(randomIndex)]
}

// query element q
func searchLevel(vertex Vertex, enterPoints []Vertex, efSize int, level int) []Vertex {
	visited := map[ID]bool{}
	candidates := map[ID]Vertex{}
	closestNeighbors := make([]Vertex, efSize)

	for _, elem := range enterPoints {
		visited[elem.id] = true
		candidates[elem.id] = elem
		closestNeighbors = append(closestNeighbors, elem)
	}

	for len(candidates) > 0 {
		nearest, nearestDist := getClosest(vertex, candidates)
		_, furthestDist := getFurthest(vertex, closestNeighbors)

		if nearestDist > furthestDist {
			break //all elements in W have been evaluated
		}

		//look for more candidates
		for _, neighbor := range nearest.neighbors {
			if _, ok := visited[neighbor.id]; ok {
				continue
			}

			visited[neighbor.id] = true
			furthest, furthestDist := getFurthest(vertex, closestNeighbors)
			if distance(neighbor, vertex) < furthestDist || len(closestNeighbors) < efSize {
				candidates[neighbor.id] = neighbor
				closestNeighbors = append(closestNeighbors, neighbor)
				if len(closestNeighbors) > efSize {
					closestNeighbors = removeVertex(furthest, closestNeighbors)
				}
			}
		}

	}
	return closestNeighbors

}

// selects M nearest neighbors
func selectNeighbors(vertex Vertex, W []Vertex, M int, level int) []Vertex { //simple
	sort.Slice(W, func(i, j int) bool {
		return distance(vertex, W[i]) > distance(vertex, W[j])
	})
	return W[:M]
}

func calculateLevel(levelMultiplier float64) int {
	//floor of -ln(unif(0,1)*mL)
	uniform := rand.Float64()
	prob := math.Log(-uniform * levelMultiplier)
	level := math.Floor(prob)

	return int(level)
}

func distance(v1 Vertex, v2 Vertex) float64 {
	if len(v1.vector) != len(v2.vector) {
		return math.Inf(1) // or any other error handling
	}
	sum := 0.0
	for i := 0; i < len(v1.vector); i++ {
		diff := float64(v1.vector[i] - v2.vector[i])
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func getClosestArr(vertex Vertex, candidates []Vertex) []Vertex {
	closest := candidates[0]
	closestDist := distance(vertex, closest)
	for _, candidate := range candidates {
		distance := distance(closest, candidate)
		if distance < closestDist {
			closest = candidate
			closestDist = distance
		}
	}
	return []Vertex{closest}
}

func getClosest(vertex Vertex, candidates map[ID]Vertex) (Vertex, float64) {
	closest := candidates[0]
	closestDist := distance(vertex, closest)
	for _, candidate := range candidates {
		distance := distance(closest, candidate)
		if distance < closestDist {
			closest = candidate
			closestDist = distance
		}
	}
	return closest, closestDist
}

func getFurthest(vertex Vertex, W []Vertex) (Vertex, float64) {
	furthest := W[0]
	furthestDist := distance(vertex, furthest)
	for _, vertex := range W {
		distance := distance(vertex, furthest)
		if distance > furthestDist {
			furthest = vertex
			furthestDist = distance
		}
	}

	return furthest, furthestDist

}

func removeVertex(vertex Vertex, vertices []Vertex) []Vertex {
	newVertices := make([]Vertex, len(vertices)-1)
	for _, v := range vertices {
		if v.id == vertex.id {
			continue
		}
		newVertices = append(newVertices, v)
	}
	return newVertices
}
