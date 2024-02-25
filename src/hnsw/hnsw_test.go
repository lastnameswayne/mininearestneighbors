package hnsw

import (
	"testing"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	"github.com/stretchr/testify/assert"
)

func buildGraphForTest() {

}

func TestSearch(t *testing.T) {
}

// func TestSearchLayer() {

// }

// func TestBuildN() {

// }

func TestInsertPoint(t *testing.T) {
	layerCount := 3
	efSize := 3
	M := 2
	M_max := 4
	hnsw := ConstructHNSW(layerCount)
	t.Run("vector is inserted into bottom layer", func(t *testing.T) {
		//All inserted vectors will be in the bottom layer

		v1 := Vector{
			Id:     1,
			Vector: []int{1, 2, 3, 4, 5},
		}

		hnsw := hnsw.InsertVector(v1, efSize, M, M_max)

		found := false
		for _, item := range hnsw.layers[0] {
			if item.Id == g.ID(v1.Id) {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("vect", func(t *testing.T) {

	})
}

func TestSearchLayer(t *testing.T) {

	hnsw := constructTestHNSW()

	q := g.Vertex{
		Id:     10,
		Vector: []int{1, 2, 3, 4, 5},
	}
	t.Run("should output an amount the ef=2 closest neighbors", func(t *testing.T) {
		layer0 := hnsw.layers[0]

		nearestInLayer := searchLayer(q, layer0, hnsw.entrancePoint, 2)

		asList := nearestInLayer.UnsortedList()

		assert.Contains(t, asList, 2)
		assert.Contains(t, asList, 4)
	})

}

func constructTestHNSW() hnsw {
	// Graph construction is first step
	v1 := Vector{
		Id:     1,
		Vector: []int{100, 1002, 313, 314, 580},
	}
	v2 := Vector{
		Id:     2,
		Vector: []int{2, 2, 3, 5, 5},
	}
	v3 := Vector{
		Id:     3,
		Vector: []int{10000, 10000, 10000, 10000, 10000},
	}
	v4 := Vector{
		Id:     4,
		Vector: []int{1, 2, 3, 4, 5}}
	layerCount := 3
	M := 2
	mMax := 2 * M
	efSize := 5

	hnsw := ConstructHNSW(layerCount)

	vs := []Vector{v1, v2, v3, v4}

	for _, vector := range vs {
		hnsw = hnsw.InsertVector(vector, efSize, M, mMax)
	}

	return hnsw
}

func main() {
	// fmt.Println("hello world")
	// //parameters
	// }
	// res := hnsw.Search(q, efSize, 3)

	// for _, vertex := range res {
	// 	fmt.Println("id", vertex.Id, "has distance", distance(vertex.Vector, q.Vector))
	// }

	// fmt.Println("the correct is")
	// for _, vector := range vs {
	// 	fmt.Println("id", vector.Id, "has distance", distance(vector.Vector, q.Vector))
	// }

}
