package r2core

import (
	"fmt"
	"testing"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment()

	if env == nil {
		t.Fatal("NewEnvironment() returned nil")
	}

	if env.store == nil {
		t.Error("Expected store to be initialized")
	}

	if env.outer != nil {
		t.Error("Expected outer to be nil for new environment")
	}

	if env.imported == nil {
		t.Error("Expected imported to be initialized")
	}

	if len(env.store) != 0 {
		t.Error("Expected store to be empty initially")
	}

	if len(env.imported) != 0 {
		t.Error("Expected imported to be empty initially")
	}
}

func TestNewInnerEnv(t *testing.T) {
	outer := NewEnvironment()
	outer.Dir = "/test/path"

	inner := NewInnerEnv(outer)

	if inner == nil {
		t.Fatal("NewInnerEnv() returned nil")
	}

	if inner.outer != outer {
		t.Error("Expected inner environment to reference outer environment")
	}

	if inner.Dir != outer.Dir {
		t.Error("Expected inner environment to inherit Dir from outer")
	}

	if inner.store == nil {
		t.Error("Expected store to be initialized in inner environment")
	}

	if len(inner.store) != 0 {
		t.Error("Expected inner store to be empty initially")
	}
}

func TestEnvironment_Set(t *testing.T) {
	env := NewEnvironment()

	// Test setting a string value
	env.Set("name", "test")
	if len(env.store) != 1 {
		t.Error("Expected store to have 1 item")
	}

	val, exists := env.store["name"]
	if !exists {
		t.Error("Expected 'name' to exist in store")
	}
	if val.Value != "test" {
		t.Errorf("Expected 'test', got %v", val.Value)
	}

	// Test setting a number value
	env.Set("age", 42)
	if len(env.store) != 2 {
		t.Error("Expected store to have 2 items")
	}

	val, exists = env.store["age"]
	if !exists {
		t.Error("Expected 'age' to exist in store")
	}
	if val.Value != 42 {
		t.Errorf("Expected 42, got %v", val.Value)
	}

	// Test overwriting a value
	env.Set("name", "updated")
	if len(env.store) != 2 {
		t.Error("Expected store to still have 2 items")
	}

	val, exists = env.store["name"]
	if !exists {
		t.Error("Expected 'name' to exist in store")
	}
	if val.Value != "updated" {
		t.Errorf("Expected 'updated', got %v", val.Value)
	}
}

func TestEnvironment_Get_LocalScope(t *testing.T) {
	env := NewEnvironment()

	// Test getting non-existent variable
	val, exists := env.Get("nonexistent")
	if exists {
		t.Error("Expected 'nonexistent' to not exist")
	}
	if val != nil {
		t.Errorf("Expected nil, got %v", val)
	}

	// Test getting existing variable
	env.Set("test", "value")
	val, exists = env.Get("test")
	if !exists {
		t.Error("Expected 'test' to exist")
	}
	if val != "value" {
		t.Errorf("Expected 'value', got %v", val)
	}
}

func TestEnvironment_Get_NestedScopes(t *testing.T) {
	outer := NewEnvironment()
	inner := NewInnerEnv(outer)

	// Set variable in outer scope
	outer.Set("outer_var", "outer_value")

	// Test accessing outer variable from inner scope
	val, exists := inner.Get("outer_var")
	if !exists {
		t.Error("Expected 'outer_var' to be accessible from inner scope")
	}
	if val != "outer_value" {
		t.Errorf("Expected 'outer_value', got %v", val)
	}

	// Set variable with same name in inner scope (shadowing)
	inner.Set("outer_var", "inner_value")

	// Test that inner scope shadows outer scope
	val, exists = inner.Get("outer_var")
	if !exists {
		t.Error("Expected 'outer_var' to exist in inner scope")
	}
	if val != "inner_value" {
		t.Errorf("Expected 'inner_value' (shadowed), got %v", val)
	}

	// Test that outer scope is unchanged
	val, exists = outer.Get("outer_var")
	if !exists {
		t.Error("Expected 'outer_var' to exist in outer scope")
	}
	if val != "outer_value" {
		t.Errorf("Expected 'outer_value', got %v", val)
	}
}

func TestEnvironment_Get_MultipleNestedScopes(t *testing.T) {
	// Create chain: global -> middle -> inner
	global := NewEnvironment()
	middle := NewInnerEnv(global)
	inner := NewInnerEnv(middle)

	// Set variables at different levels
	global.Set("global_var", "global")
	middle.Set("middle_var", "middle")
	inner.Set("inner_var", "inner")

	// Test accessing from innermost scope
	val, exists := inner.Get("global_var")
	if !exists {
		t.Error("Expected 'global_var' to be accessible from inner scope")
	}
	if val != "global" {
		t.Errorf("Expected 'global', got %v", val)
	}

	val, exists = inner.Get("middle_var")
	if !exists {
		t.Error("Expected 'middle_var' to be accessible from inner scope")
	}
	if val != "middle" {
		t.Errorf("Expected 'middle', got %v", val)
	}

	val, exists = inner.Get("inner_var")
	if !exists {
		t.Error("Expected 'inner_var' to exist in inner scope")
	}
	if val != "inner" {
		t.Errorf("Expected 'inner', got %v", val)
	}

	// Test that variables are not accessible in wrong direction
	val, exists = global.Get("inner_var")
	if exists {
		t.Error("Expected 'inner_var' to not be accessible from global scope")
	}
	if val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestEnvironment_GetStore(t *testing.T) {
	env := NewEnvironment()
	env.Set("test1", "value1")
	env.Set("test2", "value2")

	store := env.GetStore()
	if store == nil {
		t.Fatal("Expected store to not be nil")
	}

	if len(store) != 2 {
		t.Errorf("Expected store to have 2 items, got %d", len(store))
	}

	if store["test1"] != "value1" {
		t.Errorf("Expected 'value1', got %v", store["test1"])
	}

	if store["test2"] != "value2" {
		t.Errorf("Expected 'value2', got %v", store["test2"])
	}

	// Test with nil environment
	var nilEnv *Environment
	nilStore := nilEnv.GetStore()
	if nilStore != nil {
		t.Error("Expected nil store for nil environment")
	}
}

func TestEnvironment_DirectoryInheritance(t *testing.T) {
	outer := NewEnvironment()
	outer.Dir = "/path/to/project"

	inner := NewInnerEnv(outer)

	if inner.Dir != outer.Dir {
		t.Errorf("Expected inner Dir to be %q, got %q", outer.Dir, inner.Dir)
	}

	// Test that changing inner Dir doesn't affect outer
	inner.Dir = "/different/path"

	if outer.Dir != "/path/to/project" {
		t.Errorf("Expected outer Dir to remain %q, got %q", "/path/to/project", outer.Dir)
	}

	if inner.Dir != "/different/path" {
		t.Errorf("Expected inner Dir to be %q, got %q", "/different/path", inner.Dir)
	}
}

func TestEnvironment_VariableTypes(t *testing.T) {
	env := NewEnvironment()

	// Test different types of variables
	testCases := []struct {
		name  string
		value interface{}
	}{
		{"string_var", "hello"},
		{"int_var", 42},
		{"float_var", 3.14},
		{"bool_var", true},
		{"nil_var", nil},
		{"slice_var", []int{1, 2, 3}},
		{"map_var", map[string]int{"a": 1, "b": 2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env.Set(tc.name, tc.value)

			val, exists := env.Get(tc.name)
			if !exists {
				t.Errorf("Expected %q to exist", tc.name)
			}

			// Use deep comparison for complex types
			switch expected := tc.value.(type) {
			case []int:
				actual, ok := val.([]int)
				if !ok {
					t.Errorf("Expected []int, got %T", val)
				} else if len(actual) != len(expected) {
					t.Errorf("Expected slice length %d, got %d", len(expected), len(actual))
				} else {
					for i, v := range expected {
						if actual[i] != v {
							t.Errorf("Expected slice[%d] = %v, got %v", i, v, actual[i])
						}
					}
				}
			case map[string]int:
				actual, ok := val.(map[string]int)
				if !ok {
					t.Errorf("Expected map[string]int, got %T", val)
				} else if len(actual) != len(expected) {
					t.Errorf("Expected map length %d, got %d", len(expected), len(actual))
				} else {
					for k, v := range expected {
						if actual[k] != v {
							t.Errorf("Expected map[%q] = %v, got %v", k, v, actual[k])
						}
					}
				}
			default:
				if val != tc.value {
					t.Errorf("Expected %v, got %v", tc.value, val)
				}
			}
		})
	}
}

// Note: Environment is not designed for concurrent access
// Maps in Go are not thread-safe, so concurrent access would require synchronization

func TestEnvironment_ScopeIsolation(t *testing.T) {
	parent := NewEnvironment()
	child1 := NewInnerEnv(parent)
	child2 := NewInnerEnv(parent)

	// Set variables in different scopes
	parent.Set("parent", "parent_value")
	child1.Set("child1", "child1_value")
	child2.Set("child2", "child2_value")

	// Test that children can access parent
	val, exists := child1.Get("parent")
	if !exists || val != "parent_value" {
		t.Error("Child1 should be able to access parent variables")
	}

	val, exists = child2.Get("parent")
	if !exists || val != "parent_value" {
		t.Error("Child2 should be able to access parent variables")
	}

	// Test that children cannot access each other
	_, exists = child1.Get("child2")
	if exists {
		t.Error("Child1 should not be able to access child2 variables")
	}

	_, exists = child2.Get("child1")
	if exists {
		t.Error("Child2 should not be able to access child1 variables")
	}

	// Test that parent cannot access children
	_, exists = parent.Get("child1")
	if exists {
		t.Error("Parent should not be able to access child1 variables")
	}

	_, exists = parent.Get("child2")
	if exists {
		t.Error("Parent should not be able to access child2 variables")
	}
}

// Benchmark tests
func BenchmarkEnvironment_Set(b *testing.B) {
	env := NewEnvironment()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		env.Set(fmt.Sprintf("var_%d", i), i)
	}
}

func BenchmarkEnvironment_Get_Local(b *testing.B) {
	env := NewEnvironment()
	env.Set("test", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		env.Get("test")
	}
}

func BenchmarkEnvironment_Get_Nested(b *testing.B) {
	// Create deep nesting (10 levels)
	env := NewEnvironment()
	current := env
	for i := 0; i < 10; i++ {
		current = NewInnerEnv(current)
	}

	env.Set("deep_var", "value")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		current.Get("deep_var")
	}
}

func BenchmarkEnvironment_Get_NonExistent(b *testing.B) {
	// Create deep nesting (10 levels)
	env := NewEnvironment()
	current := env
	for i := 0; i < 10; i++ {
		current = NewInnerEnv(current)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		current.Get("nonexistent")
	}
}
