package Graph_test

import (
	"testing"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	graph := make(g.Graph)

	t.Run("adds two vertex and edge between two", func(t *testing.T) {
		v1 := g.Vertex{
			Id:     1,
			Vector: []int{1, 2, 3, 4, 5},
			Edges:  make([]g.ID, 0),
		}

		v2 := g.Vertex{
			Id:     2,
			Vector: []int{2, 2, 3, 5, 5},
			Edges:  make([]g.ID, 0),
		}

		graph.AddVertex(v1.Id, v1.Vector)
		graph.AddVertex(v2.Id, v2.Vector)

		assert.NotNil(t, graph[v1.Id])
		assert.NotNil(t, graph[v2.Id])

		v1, v2 = graph.AddEdge(v1, v2)

		assert.Equal(t, len(v1.Edges), 1)
		assert.Equal(t, len(v2.Edges), 1)

		assert.Equal(t, v1.Edges[0], v2.Id)
		assert.Equal(t, v2.Edges[0], v1.Id)

		//get neighborhood

		assert.Equal(t, graph.Neighborhood(v1), graph.Neighborhood(v2))

		//delete edge
		graph.RemoveEdge(v1.Id, v2.Id)

		assert.Len(t, graph.Neighborhood(v1), 0)
		assert.Len(t, graph.Neighborhood(v2), 0)
	})
}
