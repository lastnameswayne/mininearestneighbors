package hnsw

import (
	"math"
	"math/rand"
	"sort"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	"github.com/lastnameswayne/mininearestneighbors/src/heap"
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

func (hnsw *hnsw) Search(query v.Vector, efSize int, k int) []heap.Element {
	queryElement := g.Vertex{
		Vector: query.Vector,
		Id:     g.ID(query.Id),
	}
	ep := hnsw.EntrancePoint
	top := len(hnsw.Layers) - 1

	W := q.New(queryElement, heap.Min)

	for i := top; i > 0; i-- {
		layer := hnsw.Layers[i]
		W = searchLayer(queryElement, layer, ep, 1)
		ep = W.Peek().Vertex
	}
	W = searchLayer(queryElement, hnsw.Layers[0], ep, efSize)
	return getKClosest(W.Elements(), queryElement, k, hnsw.Layers[0])
}

func (hnsw hnsw) InsertVector(queryVector v.Vector, efSize int, M int, mMax int) hnsw {
	enterPointHNSW := hnsw.EntrancePoint
	top := len(hnsw.Layers) - 1
	levelMultiplier := 1 / math.Log(float64(M)) // rule of thumb is mL = 1/ln(M) where M is the number neighbors we add to each vertex on insertion

	level := calculateLevel(levelMultiplier)
	level = min(level, top)
	queryVertex := g.Vertex{Id: g.ID(queryVector.Id), Vector: queryVector.Vector, Edges: []g.ID{}}

	W := q.New(queryVertex, heap.Min)

	for i := top; i > level+1; i-- {
		layer := hnsw.Layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, 1)
		enterPointHNSW = W.Peek().Vertex
	}

	for i := level; i >= 0; i-- {
		layer := hnsw.Layers[i]
		W = searchLayer(queryVertex, layer, enterPointHNSW, efSize)
		WIds := elementsToVertices(W.Elements())
		neighbors := selectNeighbors(queryVertex, WIds, M, layer)

		for _, n := range neighbors {
			layer.AddEdge(queryVertex, n)
			queryVertex = layer[queryVertex.Id]
			n = layer[n.Id]
		}

		for _, n := range neighbors {
			neighborhood := layer.Neighborhood(n.Id)
			if len(neighborhood) > mMax {
				newNeighbors := selectNeighbors(n, neighborhood, mMax, layer)

				setNewNeighborhood(n, newNeighbors, layer)
			}
		}
	}
	if level == top {
		enterPointHNSW = W.Peek().Vertex
		hnsw.EntrancePoint = enterPointHNSW
	}
	return hnsw
}

// does a greedy search over a layer, and each layer is a graph by itself
func searchLayer(query g.Vertex, layer g.Graph, entrancePoint g.Vertex, efSize int) *q.PriorityQueue {
	visited := s.Set{} //vertices we have visited
	visited.Add(int(entrancePoint.Id))
	candidates, W := initSearchLayerHeaps(entrancePoint, query)

	for candidates.Size() > 0 {
		nearest := candidates.Pop()
		furthest := W.Pop()

		if nearest.Weight > furthest.Weight {
			break //all elements in a layer have been evaluated
		}

		neighborhood := layer.Neighborhood(nearest.Vertex.Id)
		for _, neighbor := range neighborhood {
			if visited.Has(int(neighbor)) {
				continue
			}

			visited.Add(int(neighbor))

			furthest := W.Pop()

			neighborVertex := layer[neighbor]
			neighborIsCloserThanFurthest := v.Distance(query.Vector, neighborVertex.Vector) < furthest.Weight

			if neighborIsCloserThanFurthest || W.Size() < efSize {
				candidates.Push(neighborVertex)
				W.Push(neighborVertex)

				if W.Size() > efSize {
					W.Pop()
				}
			}
		}
	}

	return W

}

func initSearchLayerHeaps(entrancePoint g.Vertex, query g.Vertex) (*q.PriorityQueue, *q.PriorityQueue) {
	candidates := q.New(query, heap.Min) //set of possible found nearest neighbors
	candidates.Push(entrancePoint)

	W := q.New(query, heap.Max) //dynamic list of found nearest neighbors
	W.Push(entrancePoint)

	return candidates, W

}
func selectNeighbors(vertex g.Vertex, W []g.ID, M int, layer g.Graph) []g.Vertex { //simple
	// if W.Has(int(vertex.Id)) {
	// 	W.Delete(int(vertex.Id))
	// }

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

func elementsToVertices(elements []heap.Element) []g.ID {
	res := make([]g.ID, len(elements))

	for _, e := range elements {
		res = append(res, e.Vertex.Id)
	}

	return res
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

func getKClosest(W []heap.Element, vertex g.Vertex, k int, layer g.Graph) []heap.Element {
	elements := W

	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Weight > elements[j].Weight
	})

	return elements[:min(len(elements), k)]
}
