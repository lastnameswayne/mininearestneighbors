package main

import (
	"fmt"
	"testing"
)

func buildGraphForTest() {

}

func test_package(*testing testing) {

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
	hnsw := ConstructHNSW()

	hnsw = insertVector(hnsw, v1, 5)
	hnsw = insertVector(hnsw, v2, 5)
	hnsw = insertVector(hnsw, v3, 5)
	hnsw = insertVector(hnsw, v4, 5)
	hnsw = insertVector(hnsw, v5, 5)
	hnsw = insertVector(hnsw, v6, 5)
	hnsw = insertVector(hnsw, v7, 5)
	hnsw = insertVector(hnsw, v8, 5)
	fmt.Println(" . ")
	for idx, layers := range hnsw.layers {
		fmt.Println("layer", idx, "we have", layers)
	}
	fmt.Println(hnsw.Search(q, 5, 3))

}

func test_search() {

}

func test_buildNeighbors() {

}
