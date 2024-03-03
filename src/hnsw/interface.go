package hnsw

import (
	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	v "github.com/lastnameswayne/mininearestneighbors/src/vector"
)

type HNSW interface {
	InsertVector(v.Vector, int, int, int) hnsw
	Search(v.Vector, int, int) []g.Vertex
}
