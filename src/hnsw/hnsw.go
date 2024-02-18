package hnsw

import (
	"math"
	"math/rand"
	"sort"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	s "github.com/lastnameswayne/mininearestneighbors/src/set"
)

type Vector struct {
	// A vector is a list of integers
	// In mesure case we also have an id
	Id     int
	Size   string
	Vector []int
}
type HNSW struct {
	layers        map[int]g.Graph
	entrancePoint g.Vertex
}

func ConstructHNSW(layerAmount int) HNSW {
	layers := map[int]g.Graph{}
	for i := 0; i < layerAmount; i++ {
		zeroNode := g.Vertex{Id: 0, Vector: []int{10000, 10000, 10000, 10000, 10000}, Edges: []g.ID{}}
		layers[i] = g.Graph{g.ID(0): zeroNode}
	}
	zeroNode := g.Vertex{Id: 0, Vector: []int{10000, 10000, 10000, 10000, 10000}, Edges: []g.ID{}}
	return HNSW{
		layers:        layers,
		entrancePoint: zeroNode,
	}
}

func (hnsw *HNSW) Search(q Vector, efSize int, k int) []g.Vertex {

	W := s.Set{}
	queryElement := g.Vertex{
		Vector: q.Vector,
		Id:     g.ID(q.Id),
	}
	ep := hnsw.entrancePoint
	top := len(hnsw.layers) - 1
	for i := top; i > 0; i-- {
		layer := hnsw.layers[i]
		W = searchLayer(queryElement, layer, ep, 1)
		ep = getClosest(queryElement, W, layer)
	}
	W = searchLayer(queryElement, hnsw.layers[0], ep, efSize)
	return getKClosest(W, queryElement, k, hnsw.layers[0])
}

func getKClosest(W s.Set, vertex g.Vertex, k int, layer g.Graph) []g.Vertex {
	vertices := make([]g.Vertex, 0)
	for id, _ := range W {
		vertex := layer[g.ID(id)]
		vertices = append(vertices, vertex)
	}

	sort.Slice(vertices, func(i, j int) bool {
		return distance(vertex.Vector, vertices[i].Vector) > distance(vertex.Vector, vertices[j].Vector)
	})

	return vertices[:min(len(vertices), k)]
}

func InsertVector(graph HNSW, queryVector Vector, efSize int, M int, mMax int) HNSW {
	enterPointHNSW := graph.entrancePoint
	top := len(graph.layers) - 1
	levelMultiplier := 1 / math.Log(float64(M)) // m_L = rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	// A vector is added to insertion layer and every layer below it
	W := s.Set{}

	level := calculateLevel(levelMultiplier)
	level = min(level, top)
	queryVertex := g.Vertex{Id: g.ID(queryVector.Id), Vector: queryVector.Vector, Edges: []g.ID{}}

	for i := top; i > level+1; i-- {
		layer := graph.layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, 1)
		enterPointHNSW = getClosest(queryVertex, W, layer)
	}

	for i := level; i >= 0; i-- {
		layer := graph.layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, efSize)
		neighbors := selectNeighbors(queryVertex, W, M, layer)

		for _, n := range neighbors {
			layer.AddEdge(queryVertex, n)
			queryVertex = layer[queryVertex.Id]
			n = layer[n.Id]
		}

		for _, n := range neighbors {
			neighborhood := layer.Neighborhood(n)
			if len(neighborhood) > mMax {
				asSet := verticesToSet(neighborhood)
				newNeighbors := selectNeighbors(n, asSet, mMax, layer)

				setNewNeighborhood(n, newNeighbors, layer)
			}
		}
		enterPointHNSW = getClosest(queryVertex, W, layer)
	}
	if level > top {
		graph.entrancePoint = enterPointHNSW
	}
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
		nearest := getClosest(vertex, candidates, layer)
		candidates.Delete(int(nearest.Id))
		furthest := getFurthest(vertex, W, layer)

		if distance(nearest.Vector, vertex.Vector) > distance(furthest.Vector, vertex.Vector) {
			break //all elements in W have been evaluated
		}

		neighborhood := layer.Neighborhood(nearest)
		for _, neighbor := range neighborhood {
			if visited.Has(int(neighbor)) {
				continue
			}

			visited.Add(int(neighbor))

			furthest := getFurthest(vertex, W, layer)

			neighborVertex, ok := layer[neighbor]
			if !ok {
				panic("errro no neighbor found")
			}
			if distance(neighborVertex.Vector, vertex.Vector) < distance(vertex.Vector, furthest.Vector) || len(W) < efSize {
				candidates.Add(int(neighbor))
				W.Add(int(neighbor))

				if len(W) > efSize {
					W.Delete(int(getFurthest(vertex, W, layer).Id))
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

	sort.Slice(vertices, func(i, j int) bool {
		return distance(vertex.Vector, vertices[i].Vector) > distance(vertex.Vector, vertices[j].Vector)
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

func distance(v1 []int, v2 []int) float64 {
	if len(v1) != len(v2) {
		return math.Inf(1) // or any other error handling
	}
	sum := 0.0
	for i := 0; i < len(v1); i++ {
		diff := float64(v1[i] - v2[i])
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

		distance := distance(vertex.Vector, candidate.Vector)
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

		distance := distance(vertex.Vector, candidate.Vector)
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
