package hnsw

import g "github.com/lastnameswayne/mininearestneighbors/src/graph"

type HNSW interface {
	InsertVector(queryVector Vector, efSize int, M int, mMax int) hnsw
	Search(q Vector, efSize int, k int) []g.Vertex
}
