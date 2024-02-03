package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	s "github.com/lastnameswayne/mininearestneighbors/src/set"
)

type Vector struct {
	// A vector is a list of integers
	// In mesure case we also have an id
	id     int
	vector []int
}
type HNSW struct {
	layers        map[int]g.Graph
	entrancePoint g.Vertex
}

func ConstructHNSW() HNSW {
	layersCount := 3
	layers := map[int]g.Graph{}
	zeroNode := g.Vertex{Id: 0, Vector: []int{0}, Edges: map[g.ID]g.Vertex{}}
	for i := 0; i < layersCount; i++ {
		layers[i] = g.Graph{g.ID(0): &zeroNode}
	}
	return HNSW{
		layers:        layers,
		entrancePoint: zeroNode,
	}
}

func main() {
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
	M := 2 // number of neighbors to add to each vertex on insertion
	M_max := 4

	enterPointHNSW := graph.entrancePoint
	// topLayer := graph.getTopLayer()
	top := len(graph.layers)
	levelMultiplier := 1 / math.Log(float64(M)) // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	W := s.Set{}

	level := calculateLevel(levelMultiplier)
	level = min(level, top)
	layerHNSW := graph.layers[level]
	queryVertex := layerHNSW.AddVertex(g.ID(queryVector.id), queryVector.vector)

	//Start at the top level and traverse greedilty to find the epSize closest neighbors to vector
	//These are used as enterPoints in the next step
	for i := top; i > level+1; i-- {
		layer := graph.layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, 1)
		enterPointHNSW, _ = getClosest(queryVertex, W, layer)
	}
	//searches again from the next layer
	for i := level; i >= 0; i-- {
		layer := graph.layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, efSize)
		neighbors := selectNeighbors(queryVertex, W, M, layer)

		for _, n := range neighbors {
			layer.AddEdge(queryVertex.Id, n.Id)
		}

		for _, n := range neighbors {
			neighborhood := layer.Neighborhood(n)
			if len(neighborhood) > M_max {
				asSet := verticesToSet(neighborhood)
				newNeighbors := selectNeighbors(n, asSet, M_max, layer)

				setNewNeighborhood(&n, newNeighbors, layer)
			}
		}
		enterPointHNSW, _ = getClosest(queryVertex, W, layer)
	}
	graph.entrancePoint = enterPointHNSW
	return graph
}

func (g HNSW) getTopLayer() g.Graph {
	top := len(g.layers)
	return g.layers[top-1]
}

func verticesToSet(vertices []g.Vertex) s.Set {
	set := s.Set{}
	for _, v := range vertices {
		set.Add(int(v.Id))
	}

	return set
}

func setNewNeighborhood(v *g.Vertex, neighborhood []g.Vertex, layer g.Graph) {
	currNeighborhood := v.Edges
	for _, n := range currNeighborhood {
		layer.RemoveEdge(v.Id, n.Id)
	}

	for _, n := range neighborhood {
		layer.AddEdge(v.Id, n.Id)
	}
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
func searchLayer(vertex g.Vertex, layer g.Graph, entrancePoint g.Vertex, efSize int) s.Set {
	visited := s.Set{}
	visited.Add(int(entrancePoint.Id))
	candidates := s.Set{}
	candidates.Add(int(entrancePoint.Id))
	W := s.Set{}
	W.Add(int(entrancePoint.Id))

	for len(candidates) > 0 {
		fmt.Println("candidates", candidates)
		nearest, nearestDist := getClosest(vertex, candidates, layer)
		_, furthestDist := getFurthest(vertex, W, layer)

		if nearestDist > furthestDist {
			break //all elements in W have been evaluated
		}

		neighborhood := layer.Neighborhood(nearest)
		for _, neighbor := range neighborhood {
			if visited.Has(int(neighbor.Id)) {
				continue
			}

			visited.Add(int(neighbor.Id))

			furthest, furthestDist := getFurthest(vertex, W, layer)

			if distance(neighbor, vertex) < furthestDist || len(W) < efSize {
				candidates.Add(int(neighbor.Id))
				W.Add(int(neighbor.Id))

				if len(W) > efSize {
					layer.RemoveEdge(entrancePoint.Id, furthest.Id)
				}
			}
		}
	}
	return W

}

// selects M nearest neighbors
func selectNeighbors(vertex g.Vertex, W s.Set, M int, layer g.Graph) []g.Vertex { //simple
	vertices := make([]g.Vertex, M)
	for id, _ := range W {
		vertex := layer[g.ID(id)]
		vertices = append(vertices, *vertex)
	}

	sort.Slice(vertices, func(i, j int) bool {
		return distance(vertex, vertices[i]) > distance(vertex, vertices[j])
	})
	return vertices[:M]
}

func calculateLevel(levelMultiplier float64) int {
	//floor of -ln(unif(0,1)*mL)
	uniform := rand.Float64()
	prob := math.Log(-uniform * levelMultiplier)
	level := math.Floor(prob)

	return int(level)
}

func distance(v1 g.Vertex, v2 g.Vertex) float64 {
	if len(v1.Vector) != len(v2.Vector) {
		return math.Inf(1) // or any other error handling
	}
	sum := 0.0
	for i := 0; i < len(v1.Vector); i++ {
		diff := float64(v1.Vector[i] - v2.Vector[i])
		sum += diff * diff
	}
	return math.Sqrt(sum)
}

func getClosest(vertex g.Vertex, candidates s.Set, level g.Graph) (g.Vertex, float64) {
	closestDist := math.Inf(1)
	var closest g.Vertex
	for id, _ := range candidates {
		candidate, ok := level[g.ID(id)]
		if !ok {
			return g.Vertex{}, math.Inf(1)
		}

		distance := distance(closest, *candidate)
		if distance <= closestDist || distance == math.Inf(1) {
			closest = *candidate
			closestDist = distance
		}
	}
	return closest, closestDist
}

func getFurthest(vertex g.Vertex, candidates s.Set, level g.Graph) (g.Vertex, float64) {
	furthestDist := math.Inf(-1)
	var furthest g.Vertex
	for id, _ := range candidates {
		candidate, ok := level[g.ID(id)]
		if !ok {
			return g.Vertex{}, math.Inf(-1)
		}

		distance := distance(furthest, *candidate)
		if distance >= furthestDist || distance == math.Inf(1) {
			furthest = vertex
			furthestDist = distance
		}
	}

	return furthest, furthestDist

}
