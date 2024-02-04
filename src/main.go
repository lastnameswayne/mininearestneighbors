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
	for i := 0; i < layersCount; i++ {
		zeroNode := g.Vertex{Id: 0, Vector: []int{10000, 10000, 10000, 10000, 10000}, Edges: []g.ID{}}
		layers[i] = g.Graph{g.ID(0): zeroNode}
	}
	zeroNode := g.Vertex{Id: 0, Vector: []int{10000, 10000, 10000, 10000, 10000}, Edges: []g.ID{}}
	return HNSW{
		layers:        layers,
		entrancePoint: zeroNode,
	}
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
	// v3 := Vector{
	// 	id:     3,
	// 	vector: []int{11, 12, 13, 14, 15},
	// }
	// v4 := Vector{
	// 	id:     4,
	// 	vector: []int{1, 1, 1, 1, 1},
	// }
	// v5 := Vector{
	// 	id:     5,
	// 	vector: []int{1, 12, 3, 4, 5},
	// }
	// v6 := Vector{
	// 	id:     6,
	// 	vector: []int{1, 2, 30, 4, 5},
	// }
	// v7 := Vector{
	// 	id:     7,
	// 	vector: []int{2, 2, 3, 4, 5},
	// }
	// v8 := Vector{
	// 	id:     8,
	// 	vector: []int{10, 200, 3, 4, 5},
	// }

	q := Vector{
		id:     9,
		vector: []int{0, 2, 3, 4, 5},
	}
	hnsw := ConstructHNSW()

	hnsw = insertVector(hnsw, v1, 5)
	hnsw = insertVector(hnsw, v2, 5)
	// hnsw = insertVector(hnsw, v3, 5)
	// hnsw = insertVector(hnsw, v4, 5)
	// hnsw = insertVector(hnsw, v5, 5)
	// hnsw = insertVector(hnsw, v6, 5)
	// hnsw = insertVector(hnsw, v7, 5)
	// hnsw = insertVector(hnsw, v8, 5)
	fmt.Println(" . ")
	for idx, layers := range hnsw.layers {
		fmt.Println("layer", idx, "we have", layers)
	}
	fmt.Println(hnsw.Search(q, 5, 3))

}

func (hnsw *HNSW) Search(q Vector, efSize int, k int) s.Set {

	W := s.Set{}
	queryElement := g.Vertex{
		Vector: q.vector,
		Id:     g.ID(q.id),
	}
	ep := hnsw.entrancePoint
	top := len(hnsw.layers) - 1
	for i := top; i > 0; i-- {
		layer := hnsw.layers[i]
		W = searchLayer(queryElement, layer, ep, 1)
		ep = getClosest(queryElement, W, layer)
	}
	W = searchLayer(queryElement, hnsw.layers[0], ep, efSize)
	fmt.Println("finished", W)
	return W
}

func insertVector(graph HNSW, queryVector Vector, efSize int) HNSW {
	M := 2 // number of neighbors to add to each vertex on insertion
	M_max := 4

	enterPointHNSW := graph.entrancePoint
	top := len(graph.layers) - 1
	levelMultiplier := 1 / math.Log(float64(M)) // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	W := s.Set{}

	level := calculateLevel(levelMultiplier)
	level = min(level, top)
	queryVertex := g.Vertex{Id: g.ID(queryVector.id), Vector: queryVector.vector, Edges: []g.ID{}}

	for i := top; i > level+1; i-- {
		layer := graph.layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, 1)
		enterPointHNSW = getClosest(queryVertex, W, layer)
	}

	for i := level; i >= 0; i-- {
		layer := graph.layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, efSize)
		neighbors := selectNeighbors(queryVertex, W, M, layer)
		fmt.Println("neighbors", neighbors)

		for _, n := range neighbors {
			queryVertex, n = layer.AddEdge(queryVertex, n)
		}

		for _, n := range neighbors {
			neighborhood := layer.Neighborhood(n)
			if len(neighborhood) > M_max {
				asSet := verticesToSet(neighborhood)
				newNeighbors := selectNeighbors(n, asSet, M_max, layer)

				setNewNeighborhood(n, newNeighbors, layer)
			}
		}
		enterPointHNSW = getClosest(queryVertex, W, layer)
	}
	graph.entrancePoint = enterPointHNSW
	return graph
}

func searchLayer(vertex g.Vertex, layer g.Graph, entrancePoint g.Vertex, efSize int) s.Set {
	visited := s.Set{}
	visited.Add(int(entrancePoint.Id))
	candidates := s.Set{}
	candidates.Add(int(entrancePoint.Id))
	W := s.Set{}
	W.Add(int(entrancePoint.Id))

	for len(candidates) > 0 {
		fmt.Println("candidates", candidates)
		nearest := getClosest(vertex, candidates, layer)
		candidates.Delete(int(nearest.Id))
		furthest := getFurthest(vertex, W, layer)

		fmt.Println("nearest", nearest, distance(nearest, vertex), "for vertex", vertex.Id)
		if distance(nearest, vertex) > distance(furthest, vertex) {
			break //all elements in W have been evaluated
		}

		neighborhood := layer.Neighborhood(nearest)
		for _, neighbor := range neighborhood {
			if visited.Has(int(neighbor)) {
				continue
			}

			visited.Add(int(neighbor))

			furthest := getFurthest(vertex, W, layer)

			neighborVertex := layer[vertex.Id]
			if distance(neighborVertex, vertex) < distance(vertex, furthest) || len(W) < efSize {
				candidates.Add(int(neighbor))
				W.Add(int(neighbor))

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
	if W.Has(int(vertex.Id)) {
		W.Delete(int(vertex.Id))
	}

	vertices := make([]g.Vertex, 0)
	for id, _ := range W {
		vertex := layer[g.ID(id)]
		vertices = append(vertices, vertex)
	}
	fmt.Println("vertices", vertices)

	sort.Slice(vertices, func(i, j int) bool {
		return distance(vertex, vertices[i]) > distance(vertex, vertices[j])
	})

	return vertices[:min(len(vertices), M)]
}

func calculateLevel(levelMultiplier float64) int {
	uniform := rand.Float64()
	prob := -math.Log(uniform * levelMultiplier)
	level := math.Floor(prob)

	if level < 0 {
		level = 0
	}

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

func getClosest(vertex g.Vertex, candidates s.Set, level g.Graph) g.Vertex {
	closestDist := math.Inf(1)
	randomId := candidates.GetRandom()
	closest := level[g.ID(randomId)]
	for id, _ := range candidates {
		candidate, ok := level[g.ID(id)]
		if !ok {
			return g.Vertex{}
		}

		distance := distance(vertex, candidate)
		if distance <= closestDist || distance == math.Inf(1) {
			closest = candidate
			closestDist = distance
		}
	}
	return closest
}

func getFurthest(vertex g.Vertex, candidates s.Set, level g.Graph) g.Vertex {
	furthestDist := math.Inf(-1)
	randomId := candidates.GetRandom()
	furthest := level[g.ID(randomId)]
	for id, _ := range candidates {
		candidate, ok := level[g.ID(id)]
		if !ok {
			return g.Vertex{}
		}

		distance := distance(vertex, candidate)
		if distance > furthestDist || distance == math.Inf(1) {
			furthest = candidate
			furthestDist = distance
		}
	}

	return furthest

}
func (g HNSW) getTopLayer() g.Graph {
	top := len(g.layers)
	return g.layers[top-1]
}

func verticesToSet(vertices []g.ID) s.Set {
	set := s.Set{}
	for _, vID := range vertices {
		set.Add(int(vID))
	}

	return set
}

func setNewNeighborhood(v g.Vertex, neighborhood []g.Vertex, layer g.Graph) {
	currNeighborhood := v.Edges
	for _, n := range currNeighborhood {
		layer.RemoveEdge(v.Id, n)
	}

	for _, n := range neighborhood {
		layer.AddEdge(v, n)
	}
}
