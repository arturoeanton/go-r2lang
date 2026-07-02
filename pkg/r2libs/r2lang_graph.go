package r2libs

import (
	"fmt"
	"sort"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// Graph represents a directed graph for knowledge representation.
type Graph struct {
	Edges      map[string][]string
	Debug      bool
	ReverseMap map[string][]string // Cached reverse edges for efficient descendant queries
}

// NewGraph initializes a new graph.
func NewGraph() *Graph {
	return &Graph{
		Edges:      make(map[string][]string),
		ReverseMap: make(map[string][]string),
	}
}

// EnableDebug enables debug mode.
func (g *Graph) EnableDebug() {
	g.Debug = true
}

// AddEdge adds a directed edge to the graph.
func (g *Graph) AddEdge(from, to string) {
	g.Edges[from] = append(g.Edges[from], to)
	g.ReverseMap[to] = append(g.ReverseMap[to], from)
	if g.Debug {
		fmt.Printf("Added edge: %s -> %s\n", from, to)
	}
}

// AddBidirectionalEdge adds a bidirectional edge.
func (g *Graph) AddBidirectionalEdge(node1, node2 string) {
	g.AddEdge(node1, node2)
	g.AddEdge(node2, node1)
}

// GetAncestors finds all ancestors of a node using DFS.
func (g *Graph) GetAncestors(node string) []string {
	visited := make(map[string]bool)
	var result []string
	g.dfs(node, visited, &result, g.ReverseMap)
	return result
}

// GetDescendants finds all descendants of a node using DFS.
func (g *Graph) GetDescendants(node string) []string {
	visited := make(map[string]bool)
	var result []string
	g.dfs(node, visited, &result, g.Edges)
	return result
}

// GetRelationshipLevel finds the shortest distance between two nodes.
func (g *Graph) GetRelationshipLevel(start, end string) int {
	queue := []string{start}
	visited := make(map[string]bool)
	levels := make(map[string]int)

	levels[start] = 0
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.Edges[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				levels[neighbor] = levels[current] + 1
				if neighbor == end {
					return levels[neighbor]
				}
				queue = append(queue, neighbor)
			}
		}
	}

	return -1 // Return -1 if nodes are not related
}

// GetShortestPath finds the shortest path between two nodes using BFS.
func (g *Graph) GetShortestPath(start, end string) []string {
	queue := [][]string{{start}}
	visited := make(map[string]bool)
	visited[start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		node := path[len(path)-1]
		if node == end {
			return path
		}

		for _, neighbor := range g.Edges[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return nil // Return nil if no path exists
}

// dfs is a helper function for depth-first search.

func (g *Graph) dfs(node string, visited map[string]bool, result *[]string, edges map[string][]string) {
	neighbors := edges[node]
	sort.Strings(neighbors) // Ordenar los vecinos
	for _, neighbor := range neighbors {
		if !visited[neighbor] {
			visited[neighbor] = true
			*result = append(*result, neighbor)
			g.dfs(neighbor, visited, result, edges)
		}
	}
}

// stringsToInterfaceSlice convierte un []string (el tipo que devuelven los
// métodos de Graph) al []interface{} que espera el resto del intérprete
// R2Lang para representar arrays.
func stringsToInterfaceSlice(strs []string) []interface{} {
	arr := make([]interface{}, len(strs))
	for i, s := range strs {
		arr[i] = s
	}
	return arr
}

// GraphObject envuelve un *Graph con una API fluida para scripts R2Lang,
// siguiendo el mismo patrón que PathObject (r2io.go) / CommandObject
// (r2os.go): Eval() devuelve el propio objeto y Getattr() resuelve sus
// métodos bajo demanda.
type GraphObject struct {
	graph *Graph
}

func (g *GraphObject) Eval(env *r2core.Environment) interface{} {
	return g
}

func (g *GraphObject) Getattr(name string) (r2core.Node, bool) {
	switch name {
	case "enableDebug":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			g.graph.EnableDebug()
			return g
		}}, true
	case "addEdge":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			from, to := graphEdgeArgs(args, "graph.addEdge")
			g.graph.AddEdge(from, to)
			return g
		}}, true
	case "addBidirectionalEdge":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			n1, n2 := graphEdgeArgs(args, "graph.addBidirectionalEdge")
			g.graph.AddBidirectionalEdge(n1, n2)
			return g
		}}, true
	case "getAncestors":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			node := graphNodeArg(args, "graph.getAncestors")
			return stringsToInterfaceSlice(g.graph.GetAncestors(node))
		}}, true
	case "getDescendants":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			node := graphNodeArg(args, "graph.getDescendants")
			return stringsToInterfaceSlice(g.graph.GetDescendants(node))
		}}, true
	case "getRelationshipLevel":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			start, end := graphEdgeArgs(args, "graph.getRelationshipLevel")
			return float64(g.graph.GetRelationshipLevel(start, end))
		}}, true
	case "getShortestPath":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			start, end := graphEdgeArgs(args, "graph.getShortestPath")
			path := g.graph.GetShortestPath(start, end)
			if path == nil {
				return nil
			}
			return stringsToInterfaceSlice(path)
		}}, true
	}
	return nil, false
}

// graphNodeArg valida y extrae un único argumento string (nombre de nodo).
func graphNodeArg(args []interface{}, label string) string {
	if len(args) < 1 {
		panic(fmt.Sprintf("%s: necesita (node)", label))
	}
	node, ok := args[0].(string)
	if !ok {
		panic(fmt.Sprintf("%s: el argumento debe ser un string (nombre de nodo)", label))
	}
	return node
}

// graphEdgeArgs valida y extrae dos argumentos string (par de nodos).
func graphEdgeArgs(args []interface{}, label string) (string, string) {
	if len(args) < 2 {
		panic(fmt.Sprintf("%s: necesita (node1, node2)", label))
	}
	n1, ok1 := args[0].(string)
	n2, ok2 := args[1].(string)
	if !ok1 || !ok2 {
		panic(fmt.Sprintf("%s: ambos argumentos deben ser strings (nombres de nodo)", label))
	}
	return n1, n2
}

// RegisterGraph expone bajo el namespace "graph" un grafo dirigido simple
// (útil para modelar relaciones/jerarquías desde un script R2Lang):
// graph.new() crea una instancia; sus métodos (addEdge, getAncestors,
// getShortestPath, ...) se acceden de forma fluida sobre el objeto
// devuelto.
func RegisterGraph(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"new": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return &GraphObject{graph: NewGraph()}
		}),
	}
	RegisterModule(env, "graph", functions)
}

/*
// Example usage
func init() {
	graph := NewGraph()
	graph.EnableDebug()

	// Add facts as edges
	graph.AddEdge("Federico", "Elias")
	graph.AddEdge("Federico", "Eugenia")
	graph.AddEdge("Elias", "Sara")
	graph.AddEdge("Elias", "Arturo")
	graph.AddEdge("Sara", "Telma")

	// Query ancestors
	queryNode := "Sara"
	fmt.Printf("\nAncestors of %s:\n", queryNode)
	ancestors := graph.GetAncestors(queryNode)
	fmt.Println(ancestors)

	// Query descendants
	fmt.Printf("\nDescendants of %s:\n", queryNode)
	descendants := graph.GetDescendants(queryNode)
	fmt.Println(descendants)

	// Check relationship level
	fmt.Printf("\nRelationship level between Federico and Telma:\n")
	level := graph.GetRelationshipLevel("Federico", "Telma")
	fmt.Println(level)

	// Find shortest path
	fmt.Printf("\nShortest path between Federico and Telma:\n")
	path := graph.GetShortestPath("Federico", "Telma")
	fmt.Println(path)
}

//*/
