package Graph

type ID int
type Vertex struct {
	id     ID
	vector []int
	Edges  map[ID]Vertex
}

type Graph map[ID]*Vertex

func (g Graph) AddVertex(id int, vector []int) {
	g[ID(id)] = &Vertex{id: ID(id), vector: vector, Edges: map[ID]Vertex{}}
}

func (g Graph) AddEdge(srcKey ID, destKey ID) bool {
	srcVal, ok := g[srcKey]
	if !ok {
		return false
	}
	destVal, ok := g[destKey]
	if !ok {
		return false
	}

	srcVal.Edges[destKey] = *destVal
	destVal.Edges[srcKey] = *srcVal
	return true
}

func (g Graph) Neighbors(v Vertex) []Vertex {
	res := make([]Vertex, 0)
	val, ok := g[v.id]
	if !ok {
		return []Vertex{}
	}
	for _, v := range val.Edges {
		res = append(res, v)

	}
	return res
}

func (g Graph) RemoveEdge(srcKey ID, destKey ID) bool {
	srcVal, ok := g[srcKey]
	if !ok {
		return false
	}
	destVal, ok := g[destKey]
	if !ok {
		return false
	}

	delete(srcVal.Edges, destVal.id)
	delete(destVal.Edges, srcVal.id)

	return true
}
