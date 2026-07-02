package r2libs

import (
	"reflect"
	"sort"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
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

// TestRegisterGraph_EndToEnd exercises the R2Lang-facing "graph" module
// (graph.new()/addEdge/getAncestors/getShortestPath/...) end to end through
// a real script, not just the underlying *Graph struct — this is the API
// surface actually reachable from a .r2 script now that RegisterGraph is
// wired into pkg/r2lang/r2lang.go and pkg/r2repl/r2repl.go.
func TestRegisterGraph_EndToEnd(t *testing.T) {
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	RegisterStd(env)
	RegisterGraph(env)

	code := `
let g = graph.new()
g.addEdge("A", "B")
g.addEdge("B", "C")
g.addEdge("A", "D")
g.addEdge("D", "C")

let ancestors = g.getAncestors("C")
let descendants = g.getDescendants("A")
let level = g.getRelationshipLevel("A", "C")
let path = g.getShortestPath("A", "C")
let noPath = g.getShortestPath("C", "A")
`
	parser := r2core.NewParser(code)
	program := parser.ParseProgram()
	program.Eval(env)

	ancestorsVal, ok := env.Get("ancestors")
	if !ok {
		t.Fatal("ancestors not found in environment")
	}
	ancestors, ok := ancestorsVal.([]interface{})
	if !ok || len(ancestors) != 3 {
		t.Fatalf("expected ancestors to be a 3-element array, got %v (%T)", ancestorsVal, ancestorsVal)
	}

	descendantsVal, _ := env.Get("descendants")
	descendants, ok := descendantsVal.([]interface{})
	if !ok || len(descendants) != 3 {
		t.Fatalf("expected descendants to be a 3-element array, got %v (%T)", descendantsVal, descendantsVal)
	}

	levelVal, _ := env.Get("level")
	if levelVal != float64(2) {
		t.Fatalf("expected level == 2 (float64), got %v (%T)", levelVal, levelVal)
	}

	pathVal, _ := env.Get("path")
	path, ok := pathVal.([]interface{})
	if !ok || len(path) != 3 {
		t.Fatalf("expected path to be a 3-element array, got %v (%T)", pathVal, pathVal)
	}
	if path[0] != "A" || path[len(path)-1] != "C" {
		t.Fatalf("expected path to start at A and end at C, got %v", path)
	}

	noPathVal, _ := env.Get("noPath")
	if noPathVal != nil {
		t.Fatalf("expected noPath to be nil (no route C->A), got %v (%T)", noPathVal, noPathVal)
	}
}
