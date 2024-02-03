package Graph

type ID int
type Vertex struct {
	Id     ID
	Vector []int
	Edges  map[ID]Vertex
}

type Graph map[ID]Vertex

func (g Graph) AddVertex(id ID, vector []int) Vertex {
	vertex := Vertex{Id: id, Vector: vector, Edges: map[ID]Vertex{}}
	g[ID(id)] = vertex
	return vertex
}

func (g Graph) AddEdge(src Vertex, dest Vertex) {
	srcVal, ok := g[src.Id]
	if !ok {
		src = g.AddVertex(src.Id, src.Vector)
	}
	destVal, ok := g[dest.Id]
	if !ok {
		dest = g.AddVertex(dest.Id, dest.Vector)

	}
	src.Edges[dest.Id] = destVal
	dest.Edges[src.Id] = srcVal
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
		panic("remove edges")
		return false
	}
	destVal, ok := g[destKey]
	if !ok {
		panic("remove edges")
		return false
	}

	delete(srcVal.Edges, destVal.Id)
	delete(destVal.Edges, srcVal.Id)

	return true
}
