package main

func buildGraphForTest() {

}

func test_package() {

}
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
	Graph := Graph{
		vertices: []Vertex{},
	}

	Graph = insertVector(Graph, v1, 5)
	Graph = insertVector(Graph, v2, 5)
	Graph = insertVector(Graph, v3, 5)
	Graph = insertVector(Graph, v4, 5)
	Graph = insertVector(Graph, v5, 5)
	Graph = insertVector(Graph, v6, 5)
	Graph = insertVector(Graph, v7, 5)
	Graph = insertVector(Graph, v8, 5)

	Graph.PrintLayers()
	fmt.Println("searching")
	fmt.Println(Search(q, Graph, 5, 3))

}

func test_search() {

}

func test_buildNeighbors() {

}
