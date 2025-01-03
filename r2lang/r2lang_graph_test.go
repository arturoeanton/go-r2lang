package r2lang

import (
	"reflect"
	"sort"
	"testing"
)

func TestGraph_AddEdge(t *testing.T) {
	graph := NewGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")

	if !reflect.DeepEqual(graph.Edges["A"], []string{"B"}) {
		t.Errorf("Expected edge A -> B, but got %v", graph.Edges["A"])
	}

	if !reflect.DeepEqual(graph.Edges["B"], []string{"C"}) {
		t.Errorf("Expected edge B -> C, but got %v", graph.Edges["B"])
	}
}

func TestGraph_GetAncestors(t *testing.T) {
	graph := NewGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("C", "D")

	ancestors := graph.GetAncestors("D")
	expected := []string{"C", "B", "A"}

	if !reflect.DeepEqual(ancestors, expected) {
		t.Errorf("Expected ancestors %v, but got %v", expected, ancestors)
	}
}

func TestGraph_GetDescendants(t *testing.T) {
	graph := NewGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("A", "C")
	graph.AddEdge("B", "D")

	descendants := graph.GetDescendants("A")
	expected := []string{"B", "C", "D"}

	sort.Strings(descendants) // Ordenar para comparación
	sort.Strings(expected)

	if !reflect.DeepEqual(descendants, expected) {
		t.Errorf("Expected descendants %v, but got %v", expected, descendants)
	}
}

func TestGraph_GetRelationshipLevel(t *testing.T) {
	graph := NewGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("C", "D")

	level := graph.GetRelationshipLevel("A", "D")
	if level != 3 {
		t.Errorf("Expected relationship level 3, but got %d", level)
	}

	level = graph.GetRelationshipLevel("D", "A")
	if level != -1 {
		t.Errorf("Expected relationship level -1 for unconnected nodes, but got %d", level)
	}
}

func TestGraph_GetShortestPath(t *testing.T) {
	graph := NewGraph()
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("A", "D")
	graph.AddEdge("D", "C")

	path := graph.GetShortestPath("A", "C")
	expected := []string{"A", "B", "C"}

	if !reflect.DeepEqual(path, expected) {
		t.Errorf("Expected shortest path %v, but got %v", expected, path)
	}
}
