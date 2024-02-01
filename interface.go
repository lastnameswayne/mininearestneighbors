package main

type NN interface {
	Insert(v Vertex, g HNSW) error
	Search(v Vertex, g HNSW) Vertex
}
