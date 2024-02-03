package main

import (
	"fmt"
	"math"
	"math/rand"
	"mininearestneighbors/graph"
	"sort"
)

type Vector struct {
	// A vector is a list of integers
	// In mesure case we also have an id
	id     int
	vector []int
}
type HNSW struct {
	layers        map[int]graph.Graph
	entrancePoint graph.Vertex
}

func ConstructHNSW() HNSW {
	layersCount := 3
	layers := map[int]Graph{}
	zeroNode := Vertex{id: 0, vector: []int{0}}
	for i := 0; i < layersCount; i++ {
		layers[i] = Graph{ID(0): &zeroNode}
	}
	return HNSW{
		layers:        layers,
		entrancePoint: zeroNode,
	}
}

func main() {
}

func (v Vertex) String() string {
	return fmt.Sprintf("id %d \n level %d \n neighbors %v \n", int(v.id), v.level, v.neighbors)
}

// func (g *HNSW) Search(q Vector, efSize int, k int) []Vertex {

// 	W := []Vertex{}
// 	queryElement := Vertex{
// 		vector:    q.vector,
// 		id:        q.id,
// 		neighbors: []Vertex{},
// 	}
// 	ep := graph.entrancePoint
// 	enterPointLevel := enterPoint.level
// 	for i := enterPointLevel; i > 0; i-- {
// 		W = searchLevel(queryElement, graph.enterPoint, efSize, i)
// 		enterPoint = getClosest(queryElement, W)[0]
// 	}
// 	W = searchLevel(queryElement, []Vertex{enterPoint}, efSize, 0)
// 	return W[:k]
// }

func insertVector(graph HNSW, queryVector Vector, efSize int) HNSW {
	vertex := Vertex{
		vector:    queryVector.vector,
		id:        queryVector.id,
		neighbors: []Vertex{},
	}
	M := 2 // number of neighbors to add to each vertex on insertion
	M_max := 4

	enterPointHNSW := graph.entrancePoint
	// topLayer := graph.getTopLayer()
	top := len(graph.layers)
	levelMultiplier := 1 / math.Log(float64(M)) // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	nearestElements := []Vertex{}

	level := calculateLevel(levelMultiplier)
	level = min(level, top)
	layerHNSW := graph.layers[level]
	layerHNSW[vertex.id] = append(layerHNSW[vertex.id], vertex)

	//Start at the top level and traverse greedilty to find the epSize closest neighbors to vector
	//These are used as enterPoints in the next step
	for i := top; i > level+1; i-- {
		layer := graph.layers[i]
		nearestElements = searchLevel(vertex, layer, enterPointHNSW, 1)
		enterPointHNSW, _ = getClosest(vertex, nearestElements)
	}
	//searches again from the next layer
	for i := level; i > -1; i-- {
		layer := graph.layers[i]
		nearestElements = searchLevel(vertex, layer, enterPointHNSW, efSize)
		neighbors := selectNeighbors(vertex, nearestElements, M, level)

		for _, n := range neighbors {
			vertex, n = makeConnection(vertex, n)
		}

		for _, n := range neighbors {
			neighbors := n.neighbors
			if len(neighbors) > M_max {
				newNeighbors := selectNeighbors(n, neighbors, M_max, level)

				//remove neighbors
			}
		}
		enterPointHNSW, _ := getClosest(vertex, nearestElements)
	}
	graph.entrancePoint = enterPointHNSW
	return graph
}

func (g HNSW) getTopLayer() Graph {
	top := len(g.layers)
	return g.layers[top-1]
}

// // this can be random
// func getEnterPoint(vertex Vertex, graph HNSW) Vertex {
// 	if len(graph.vertices) == 0 {
// 		return vertex
// 	}
// 	randomIndex := math.Floor(rand.ExpFloat64() * float64(len(graph.vertices)))
// 	randomIndex = 0
// 	return graph.vertices[int(randomIndex)]
// }

// query element q
func searchLevel(vertex Vertex, level Graph, entrancePoint Vertex, efSize int) []Vertex {
	visited := map[ID]Vertex{entrancePoint.id: entrancePoint}
	candidates := map[ID]Vertex{entrancePoint.id: entrancePoint}
	closestNeighbors := make([]Vertex, efSize)
	closestNeighbors = append(closestNeighbors, entrancePoint)

	for len(candidates) > 0 {
		nearest, nearestDist := getClosest(vertex, candidates)
		_, furthestDist := getFurthest(vertex, closestNeighbors)

		if nearestDist > furthestDist {
			break //all elements in W have been evaluated
		}

		//look for more candidates
		fmt.Println("neighbors", nearest)
		for _, neighbor := range nearest.neighbors {
			if _, ok := visited[neighbor.id]; ok {
				continue
			}

			visited[neighbor.id] = neighbor // mark as seen
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

func makeConnection(a Vertex, b Vertex) (Vertex, Vertex) {
	a.neighbors = append(a.neighbors, a)
	b.neighbors = append(b.neighbors, b)
	return a, b
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

func getClosest(vertex Vertex, candidates []Vertex) (Vertex, float64) {
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
