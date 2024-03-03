package hnsw

import (
	"math"
	"math/rand"
	"sort"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	q "github.com/lastnameswayne/mininearestneighbors/src/priorityqueue"
	s "github.com/lastnameswayne/mininearestneighbors/src/set"
	v "github.com/lastnameswayne/mininearestneighbors/src/vector"
)

type hnsw struct {
	Layers        map[int]g.Graph
	EntrancePoint g.Vertex
}

var _dummyNode = g.Vertex{Id: 0, Vector: []int{100000000, 10000, 10000, 10000, 10000}}

func ConstructHNSW(layerAmount int) hnsw {
	layers := map[int]g.Graph{}
	for i := 0; i < layerAmount; i++ {
		layers[i] = g.Graph{g.ID(0): _dummyNode}
	}
	return hnsw{
		Layers:        layers,
		EntrancePoint: _dummyNode,
	}
}

func (hnsw *hnsw) Search(q v.Vector, efSize int, k int) []g.Vertex {
	W := s.Set{}
	queryElement := g.Vertex{
		Vector: q.Vector,
		Id:     g.ID(q.Id),
	}
	ep := hnsw.EntrancePoint
	top := len(hnsw.Layers) - 1
	for i := top; i > 0; i-- {
		layer := hnsw.Layers[i]
		W = searchLayer(queryElement, layer, ep, 1)
		ep = getClosest(queryElement, W, layer)
	}
	W = searchLayer(queryElement, hnsw.Layers[0], ep, efSize)
	return getKClosest(W, queryElement, k, hnsw.Layers[0])
}

func (hnsw hnsw) InsertVector(queryVector v.Vector, efSize int, M int, mMax int) hnsw {
	enterPointHNSW := hnsw.EntrancePoint
	top := len(hnsw.Layers) - 1
	levelMultiplier := 1 / math.Log(float64(M)) // rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	W := s.Set{}

	level := calculateLevel(levelMultiplier)
	level = min(level, top)
	queryVertex := g.Vertex{Id: g.ID(queryVector.Id), Vector: queryVector.Vector, Edges: []g.ID{}}

	for i := top; i > level+1; i-- {
		layer := hnsw.Layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, 1)
		enterPointHNSW = getClosest(queryVertex, W, layer)
	}

	for i := level; i >= 0; i-- {
		layer := hnsw.Layers[i]
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
	}
	if level == top {
		enterPointHNSW = getClosest(queryVertex, W, hnsw.getTopLayer())
		hnsw.EntrancePoint = enterPointHNSW
	}
	return hnsw
}

// does a greedy search over a layer, and each layer is a graph by itself
func searchLayer(query g.Vertex, layer g.Graph, entrancePoint g.Vertex, efSize int) s.Set {
	visited := s.Set{} //vertices we have visited
	visited.Add(int(entrancePoint.Id))
	candidates, W := initSearchLayerHeaps(entrancePoint, query)

	for candidates.Size() > 0 {
		nearest := candidates.Pop()
		furthest := getFurthest(query, W, layer)

		if v.Distance(nearest.Vector, query.Vector) > v.Distance(furthest.Vector, query.Vector) {
			break //all elements in a layer have been evaluated
		}

		neighborhood := layer.Neighborhood(nearest)
		for _, neighbor := range neighborhood {
			if visited.Has(int(neighbor)) {
				continue
			}

			visited.Add(int(neighbor))

			furthest := getFurthest(query, W, layer)

			neighborVertex := layer[neighbor]
			neighborIsCloserThanFurthest := v.Distance(query.Vector, neighborVertex.Vector) < v.Distance(query.Vector, furthest.Vector)

			if neighborIsCloserThanFurthest || len(W) < efSize {
				candidates.Add(int(neighbor))
				W.Add(int(neighbor))

				if len(W) > efSize {
					W.Delete(int(getFurthest(query, W, layer).Id))
				}
			}
		}
	}
	return W

}

func initSearchLayerHeaps(entrancePoint g.Vertex, query g.Vertex) (q.PriorityQueue, q.PriorityQueue) {
	candidates := q.New(query) //set of possible found nearest neighbors
	candidates.Push(entrancePoint)

	W := q.New(query) //dynamic list of found nearest neighbors
	W.Push(entrancePoint)

	return candidates, W

}
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
		return v.Distance(vertex.Vector, vertices[i].Vector) > v.Distance(vertex.Vector, vertices[j].Vector)
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

func getClosest(vertex g.Vertex, candidates s.Set, level g.Graph) g.Vertex {
	closestDist := math.Inf(1)
	randomId := candidates.GetRandom()
	closest := level[g.ID(randomId)]
	for id, _ := range candidates {
		candidate, ok := level[g.ID(id)]
		if !ok {
			return g.Vertex{}
		}

		distance := v.Distance(vertex.Vector, candidate.Vector)
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

		distance := v.Distance(vertex.Vector, candidate.Vector)
		if distance > furthestDist || distance == math.Inf(1) {
			furthest = candidate
			furthestDist = distance
		}
	}

	return furthest

}
func (g hnsw) getTopLayer() g.Graph {
	top := len(g.Layers)
	return g.Layers[top-1]
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

func getKClosest(W s.Set, vertex g.Vertex, k int, layer g.Graph) []g.Vertex {
	vertices := make([]g.Vertex, 0)
	for id, _ := range W {
		vertex := layer[g.ID(id)]
		vertices = append(vertices, vertex)
	}

	sort.Slice(vertices, func(i, j int) bool {
		return v.Distance(vertex.Vector, vertices[i].Vector) > v.Distance(vertex.Vector, vertices[j].Vector)
	})

	return vertices[:min(len(vertices), k)]
}
