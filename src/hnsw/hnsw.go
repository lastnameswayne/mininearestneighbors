package hnsw

import (
	"math"
	"math/rand"
	"sort"

	"encoding/json"
	"os"

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

var _dummyNode = g.Vertex{Id: "0", Vector: []int{100000000, 10000, 10000, 10000, 10000}}

func ConstructHNSW(layerAmount int) hnsw {
	layers := map[int]g.Graph{}
	for i := 0; i < layerAmount; i++ {
		layers[i] = g.Graph{g.ID("0"): _dummyNode}
	}
	return hnsw{
		Layers:        layers,
		EntrancePoint: _dummyNode,
	}
}

func (hnsw *hnsw) Serialize() ([]byte, error) {
	res, err := json.Marshal(hnsw)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Deserialize(hnswString []byte) (hnsw, error) {
	var hnsw hnsw
	err := json.Unmarshal(hnswString, &hnsw)
	if err != nil {
		return hnsw, err
	}
	return hnsw, nil
}

func (hnsw *hnsw) WriteToFile() (int, error) {
	file, err := os.Create("file.txt")

	if err != nil {
		return 0, err
	}
	defer file.Close()

	bytes, err := hnsw.Serialize()
	if err != nil {
		return 0, err
	}

	return file.Write(bytes)
}

func (hnsw *hnsw) Search(query v.Vector, efSize int, k int) []g.Vertex {
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

	closest := getKClosest(W.Elements(), queryElement, k) //this should be heap sort

	res := []g.Vertex{}
	for _, elem := range closest {
		res = append(res, elem.Vertex)
	}
	return res
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
	visited := s.New[g.ID]() //vertices we have visited
	visited.Add(entrancePoint.Id)
	candidates, W := initSearchLayerHeaps(entrancePoint, query)

	for candidates.Size() > 0 {
		nearest := candidates.Pop()
		furthest := W.Peek()

		if nearest.Weight > furthest.GetWeight() {
			break //all elements in a layer have been evaluated
		}

		neighborhood := layer.Neighborhood(nearest.Vertex.Id)
		for _, neighbor := range neighborhood {
			if visited.Has(neighbor) {
				continue
			}

			visited.Add(neighbor)

			neighborVertex := layer[neighbor]
			neighborIsCloserThanFurthest := v.Distance(query.Vector, neighborVertex.Vector) < furthest.GetWeight()

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
	vertices := make([]g.Vertex, 0)
	for _, id := range W {
		vertex := layer[g.ID(id)]
		vertices = append(vertices, vertex)
	}

	sort.Slice(vertices, func(i, j int) bool {
		return v.Distance(vertex.Vector, vertices[i].Vector) > v.Distance(vertex.Vector, vertices[j].Vector)
	})

	return vertices[:min(len(vertices), M)]
}

// Input: base element q, candidate elements C, number of neighbors to
// return M, layer number lc, flag indicating whether or not to extend
// candidate list extendCandidates, flag indicating
func selectNeighborsHeuristic(vertex g.Vertex, W []g.ID, M int, layer g.Graph, extendCandidates, keepPrunedConnections bool) []g.Vertex { //simple
	R := s.New[g.ID]()
	W_queue := q.New(vertex, heap.Min)
	W_seen := s.New[g.ID]()
	for _, elem := range W {
		W_queue.Push(layer[elem])
	}

	if extendCandidates {
		for _, elem := range W {
			for _, neighbor := range layer.Neighborhood(elem) {
				if !W_seen.Has(neighbor) {
					W_seen.Add(neighbor)
					W_queue.Push(layer[neighbor])
				}
			}
		}
	}

	discarded := q.New(vertex, heap.Min)
	for W_queue.Size() > 0 && len(R) < M {
		closest := W_queue.Pop()

		closestToQInR := math.Inf(1)
		for elem, _ := range R {
			distance := v.Distance(layer[g.ID(elem)].Vector, vertex.Vector)
			closestToQInR = min(closestToQInR, distance)
		}

		if closest.Weight < closestToQInR {
			R.Add(closest.Vertex.Id)
			break
		}

		discarded.Push(closest.Vertex)
	}

	if keepPrunedConnections {
		for discarded.Size() > 0 && len(R) < M {
			closest := discarded.Pop()
			R.Add(closest.Vertex.Id)
		}
	}

	result := []g.Vertex{}
	for elem, _ := range R {
		vertex, ok := layer[elem]
		if !ok {
			continue
		}

		result = append(result, vertex)
	}

	return result
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

func (g hnsw) getTopLayer() g.Graph {
	top := len(g.Layers)
	return g.Layers[top-1]
}

func elementsToVertices(elements []heap.Element) []g.ID {
	res := make([]g.ID, 0)

	for _, e := range elements {
		res = append(res, e.Vertex.Id)
	}

	return res
}

func setNewNeighborhood(v g.Vertex, newNeighborhood []g.Vertex, layer g.Graph) {
	currNeighborhood := v.Edges
	for _, n := range currNeighborhood {
		layer.RemoveEdge(v.Id, n)
	}

	for _, n := range newNeighborhood {
		layer.AddEdge(v, n)
	}
}

func getKClosest(W []heap.Element, vertex g.Vertex, k int) []heap.Element {
	sort.Slice(W, func(i, j int) bool {
		return W[i].Weight < W[j].Weight
	})

	return W[:min(len(W), k)]
}
