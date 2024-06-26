package hnsw

import (
	"os"
	"testing"

	g "github.com/lastnameswayne/mininearestneighbors/src/graph"
	v "github.com/lastnameswayne/mininearestneighbors/src/vector"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("correctly returns three closest", func(t *testing.T) {

		q := v.Vector{Id: "10", Vector: []int{1, 2, 3, 4, 5}}
		hnsw := testHNSW()

		efSize := 5
		vertices := hnsw.Search(q, efSize, 3)

		assert.Equal(t, 3, len(vertices))
		assert.Equal(t, g.ID("4"), vertices[0].Id)
		assert.Equal(t, g.ID("2"), vertices[1].Id)
	})
}

func TestConstructHNSW(t *testing.T) {
	t.Run("construct hnsw with 5 layers with a zero node in each layer", func(t *testing.T) {

		hnsw := ConstructHNSW(5)
		assert.Len(t, hnsw.Layers, 5)

		for _, layer := range hnsw.Layers {
			assert.Equal(t, layer["0"].Id, g.ID("0"))
		}
	})

}

func TestSerialization(t *testing.T) {
	t.Run("hnsw is serialized", func(t *testing.T) {
		hnsw := testHNSW()
		serialized, err := hnsw.Serialize()
		assert.NoError(t, err)
		assert.True(t, len(serialized) > 0)
	})

	t.Run("hnsw deserialized and serialized is the same", func(t *testing.T) {
		hnsw := testHNSW()
		serialized, err := hnsw.Serialize()
		assert.NoError(t, err)

		deserialized, err := Deserialize(serialized)
		assert.NoError(t, err)

		assert.Equal(t, hnsw, deserialized)
	})
}
func TestWriteFile(t *testing.T) {

	hnsw := testHNSW()

	t.Run("write file", func(t *testing.T) {

		bytesAmount, err := hnsw.WriteToFile()
		assert.NoError(t, err)
		assert.True(t, bytesAmount > 0)
		assert.FileExists(t, "/Users/tiarnanswayne/mininearestneighbors/src/hnsw/file.txt")
	})

	t.Run("read file into hnsw", func(t *testing.T) {

		f, err := os.Open("/Users/tiarnanswayne/mininearestneighbors/src/hnsw/file.txt")
		assert.NoError(t, err)

		fi, err := f.Stat()
		assert.NoError(t, err)
		bytes := make([]byte, fi.Size())

		n, err := f.Read(bytes)
		assert.NoError(t, err)
		assert.True(t, n > 0)

		newhnsw, err := Deserialize(bytes)
		assert.NoError(t, err)
		assert.Equal(t, hnsw, newhnsw)

	})
}

func TestSetNewNeighborhood(t *testing.T) {
	t.Run("new neighborhood is set", func(t *testing.T) {
		hnsw := testHNSW()
		bottomLayer := hnsw.Layers[0]
		vertex := bottomLayer[g.ID("1")]

		newNeighborhood := []g.Vertex{
			{
				Id:     "6",
				Vector: []int{400, 400, 400, 400, 400},
				Edges:  []g.ID{"5", "8"},
			},
		}

		setNewNeighborhood(vertex, newNeighborhood, bottomLayer)

		vertex = bottomLayer[g.ID("1")]
		assert.Equal(t, []g.ID{newNeighborhood[0].Id}, vertex.Edges)
	})
}

func TestSelectNeighbors(t *testing.T) {

}

func TestInsertPoint(t *testing.T) {
	layerCount := 3
	efSize := 3
	M := 2
	M_max := 4
	t.Run("vector is inserted into bottom layer", func(t *testing.T) {
		//All inserted vectors will be in the bottom layer
		hnsw := ConstructHNSW(layerCount)
		v1 := v.Vector{
			Id:     "1",
			Vector: []int{1, 2, 3, 4, 5},
		}

		hnsw = hnsw.InsertVector(v1, efSize, M, M_max)

		found := false
		for _, item := range hnsw.Layers[0] {
			if item.Id == g.ID(v1.Id) {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("inserting five vectors has five in the bottom layer", func(t *testing.T) {
		hnsw := ConstructHNSW(layerCount)
		v1 := v.Vector{
			Id:     "1",
			Vector: []int{100, 1002, 313, 314, 580},
		}
		v2 := v.Vector{
			Id:     "2",
			Vector: []int{2, 2, 3, 5, 5},
		}
		v3 := v.Vector{
			Id:     "3",
			Vector: []int{10000, 10000, 10000, 10000, 10000},
		}
		v4 := v.Vector{
			Id:     "4",
			Vector: []int{1, 2, 3, 4, 5}}

		vs := []v.Vector{v1, v2, v3, v4}

		for _, vector := range vs {
			hnsw = hnsw.InsertVector(vector, efSize, M, M_max)
		}

		found := false
		assert.Len(t, hnsw.Layers[0], 5) //4 plus the dummy node
		for _, vertex := range hnsw.Layers[0] {
			assert.True(t, len(vertex.Edges) <= M_max)
			for _, neighbor := range vertex.Edges {
				if neighbor != g.ID("0") {
					found = true
					break
				}
			}

		}
		assert.True(t, found)
	})
}

func TestSearchLayer(t *testing.T) {
	hnsw := testHNSW()

	q := g.Vertex{
		Id:     "10",
		Vector: []int{1, 2, 3, 4, 5},
	}

	t.Run("should output an amount the ef=2 closest neighbors", func(t *testing.T) {
		layer0 := hnsw.Layers[0]

		nearestInLayer := searchLayer(q, layer0, hnsw.EntrancePoint, 2)

		asList := nearestInLayer.Elements()
		t.Log(asList)

		found2 := false
		found1 := false
		for _, elem := range asList {
			if elem.Vertex.Id == "2" {
				found2 = true
			}
			if elem.Vertex.Id == "4" {
				found1 = true
			}
		}

		assert.True(t, found1)
		assert.True(t, found2)

	})

	t.Run("entry point in top layer should be node 1", func(t *testing.T) {
		topLayer := hnsw.getTopLayer()
		nearestInLayer := searchLayer(q, topLayer, hnsw.EntrancePoint, 1)

		asList := nearestInLayer.Elements()

		assert.Equal(t, asList[0].Vertex.Id, g.ID("1"))

	})

	t.Run("entry point in middle layer should be node 1 and 2", func(t *testing.T) {
		layer1 := hnsw.Layers[1]

		nearestInLayer := searchLayer(q, layer1, hnsw.EntrancePoint, 2)

		asList := nearestInLayer.Elements()
		found2 := false
		found1 := false
		for _, elem := range asList {
			if elem.Vertex.Id == "1" {
				found2 = true
			}
			if elem.Vertex.Id == "2" {
				found1 = true
			}
		}

		assert.True(t, found1)
		assert.True(t, found2)
	})

}

func testHNSW() hnsw {
	layerAmount := 3
	//max 3 connections per node
	// Graph construction is first step
	v1 := g.Vertex{
		Id:     "1",
		Vector: []int{100, 100, 100, 100, 100},
		Edges:  []g.ID{"2", "3", "4"},
	}
	v2 := g.Vertex{
		Id:     "2",
		Vector: []int{2, 2, 3, 5, 5},
		Edges:  []g.ID{"1", "4"},
	}
	v3 := g.Vertex{
		Id:     "3",
		Vector: []int{10000, 10000, 10000, 10000, 10000},
		Edges:  []g.ID{"1", "4", "5"},
	}
	v4 := g.Vertex{
		Id:     "4",
		Vector: []int{1, 2, 3, 4, 5},
		Edges:  []g.ID{"1", "2", "3"},
	}
	v5 := g.Vertex{
		Id:     "5",
		Vector: []int{300, 300, 300, 300, 300},
		Edges:  []g.ID{"3", "6", "7"},
	}
	v6 := g.Vertex{
		Id:     "6",
		Vector: []int{400, 400, 400, 400, 400},
		Edges:  []g.ID{"5", "8"},
	}
	v7 := g.Vertex{
		Id:     "7",
		Vector: []int{500, 500, 500, 500, 500},
		Edges:  []g.ID{"5", "8"},
	}
	v8 := g.Vertex{
		Id:     "8",
		Vector: []int{600, 600, 600, 600, 600},
		Edges:  []g.ID{"6", "7"},
	}
	layer0 := []g.Vertex{v1, v2, v3, v4, v5, v6, v7, v8}
	layers := map[int]g.Graph{}
	layers[0] = createAndAddLayer(layer0, layerAmount)

	v1 = g.Vertex{
		Id:     "1",
		Vector: []int{100, 100, 100, 100, 100},
		Edges:  []g.ID{"2", "5"},
	}
	v2 = g.Vertex{
		Id:     "2",
		Vector: []int{2, 2, 3, 5, 5},
		Edges:  []g.ID{"1", "8"},
	}
	v5 = g.Vertex{
		Id:     "5",
		Vector: []int{300, 300, 300, 300, 300},
		Edges:  []g.ID{"1", "8"},
	}
	v7 = g.Vertex{
		Id:     "7",
		Vector: []int{500, 500, 500, 500, 500},
		Edges:  []g.ID{"5", "8"},
	}
	v8 = g.Vertex{
		Id:     "8",
		Vector: []int{600, 600, 600, 600, 600},
		Edges:  []g.ID{"2", "5", "7"},
	}

	layer1 := []g.Vertex{v1, v2, v3, v4, v5, v6, v7, v8}
	layers[1] = createAndAddLayer(layer1, layerAmount)

	v1 = g.Vertex{
		Id:     "1",
		Vector: []int{100, 100, 100, 100, 100}, Edges: []g.ID{"7"},
	}
	v7 = g.Vertex{
		Id:     "7",
		Vector: []int{500, 500, 500, 500, 500},
		Edges:  []g.ID{"5", "8"},
	}

	layer2 := []g.Vertex{v1, v7}
	layers[2] = createAndAddLayer(layer2, layerAmount)

	return hnsw{
		Layers:        layers,
		EntrancePoint: v1,
	}
}

func createAndAddLayer(layer []g.Vertex, layerAmount int) g.Graph {
	layer0g := g.Graph{}
	for _, v := range layer {
		layer0g[v.Id] = v
	}

	layers := map[int]g.Graph{}
	for i := 0; i < layerAmount; i++ {
		layers[i] = g.Graph{g.ID("0"): _dummyNode}
	}

	return layer0g
}
