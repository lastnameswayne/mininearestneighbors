package hnsw

import (
	"fmt"
	"testing"
)

func buildGraphForTest() {

}

func TestSearch(t *testing.T) {
}

// func TestSearchLayer() {

// }

// func TestBuildN() {

// }
func main() {
	fmt.Println("hello world")
	//parameters
	layerCount := 10
	M := 5
	mMax := 2 * M //recommended
	efSize := 5

	// Graph construction is first step
	v1 := Vector{
		Id:     1,
		Vector: []int{1, 2, 3, 4, 5},
	}
	v2 := Vector{
		Id:     2,
		Vector: []int{2, 2, 3, 5, 5},
	}
	v3 := Vector{
		Id:     3,
		Vector: []int{11, 12, 13, 14, 15},
	}
	v4 := Vector{
		Id:     4,
		Vector: []int{10000, 10000, 10000, 10000, 10000},
	}
	v5 := Vector{
		Id:     5,
		Vector: []int{1, 1002, 3, 4, 5},
	}
	v6 := Vector{
		Id:     6,
		Vector: []int{1, 2, 30, 4, 5},
	}
	v7 := Vector{
		Id:     7,
		Vector: []int{2, 2, 3, 4, 5},
	}
	v8 := Vector{
		Id:     8,
		Vector: []int{10, 200, 3, 4, 5},
	}
	v10 := Vector{
		Id:     10,
		Vector: []int{10, 2, 3, 4, 5},
	}
	v11 := Vector{
		Id:     11,
		Vector: []int{0, 2, 300000, 4, 5},
	}

	hnsw := ConstructHNSW(layerCount)

	vs := []Vector{v1, v2, v3, v4, v5, v6, v7, v8, v10, v11}

	for _, vector := range vs {
		hnsw = InsertVector(hnsw, vector, efSize, M, mMax)
	}

	for idx, layers := range hnsw.layers {
		fmt.Println("layer", idx, "we have", layers)
	}
	q := Vector{
		Id:     9,
		Vector: []int{0, 2, 3, 4, 5},
	}
	res := hnsw.Search(q, efSize, 3)

	for _, vertex := range res {
		fmt.Println("id", vertex.Id, "has distance", distance(vertex.Vector, q.Vector))
	}

	fmt.Println("the correct is")
	for _, vector := range vs {
		fmt.Println("id", vector.Id, "has distance", distance(vector.Vector, q.Vector))
	}

}
