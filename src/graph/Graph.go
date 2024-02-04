package Graph

type ID int
type Vertex struct {
	Id     ID
	Vector []int
	Edges  []ID
}

type Graph map[ID]Vertex

func (g *Graph) AddVertex(id ID, vector []int) Vertex {
	vertex := Vertex{Id: id, Vector: vector, Edges: []ID{}}
	(*g)[ID(id)] = vertex
	return vertex
}

func (g *Graph) AddEdge(src Vertex, dest Vertex) (Vertex, Vertex) {
	srcVal, ok := (*g)[src.Id]
	if !ok {
		src = g.AddVertex(src.Id, src.Vector)
	}
	destVal, ok := (*g)[dest.Id]
	if !ok {
		dest = g.AddVertex(dest.Id, dest.Vector)

	}
	srcVal.Edges = append(srcVal.Edges, destVal.Id)
	destVal.Edges = append(destVal.Edges, srcVal.Id)
	return srcVal, destVal
}

func (g *Graph) Neighborhood(v Vertex) []ID {
	val, ok := (*g)[v.Id]
	if !ok {
		return []ID{}
	}
	return val.Edges
}

func (g *Graph) RemoveEdge(srcKey ID, destKey ID) {
	srcVal, ok := (*g)[srcKey]
	if !ok {
		panic("remove edges")
	}
	destVal, ok := (*g)[destKey]
	if !ok {
		panic("remove edges")
	}
	newsrcvalEdges := make([]ID, 0)

	for _, edge := range srcVal.Edges {
		if destVal.Id == edge {
			continue
		}
		newsrcvalEdges = append(newsrcvalEdges, edge)
	}

	newdestvalEdges := make([]ID, 0)
	for _, edge := range destVal.Edges {
		if srcVal.Id == edge {
			continue
		}
		newdestvalEdges = append(newdestvalEdges, edge)
	}
	srcVal.Edges = newsrcvalEdges
	destVal.Edges = newdestvalEdges
}
