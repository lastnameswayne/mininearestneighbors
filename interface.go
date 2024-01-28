package main

type NN interface {
	Insert(v Vertex, g Graph) error
	Search(v Vertex, g Graph) Vertex
}
