package r2lang

import (
	"fmt"
	"sort"
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
