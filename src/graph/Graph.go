package Graph

type ID int
type Vertex struct {
	Id     ID
	Vector []int
	Edges  map[ID]Vertex
}

type Graph map[ID]*Vertex

func (g Graph) AddVertex(id ID, vector []int) Vertex {
	vertex := Vertex{Id: id, Vector: vector, Edges: map[ID]Vertex{}}
	g[ID(id)] = &vertex
	return vertex
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

func (g Graph) Neighborhood(v Vertex) []Vertex {
	res := make([]Vertex, 0)
	val, ok := g[v.Id]
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

	delete(srcVal.Edges, destVal.Id)
	delete(destVal.Edges, srcVal.Id)

	return true
}
