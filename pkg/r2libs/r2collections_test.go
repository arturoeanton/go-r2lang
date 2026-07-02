package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestCollectionsFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterCollections(env)

	collectionsModuleObj, ok := env.Get("collections")
	if !ok {
		t.Fatal("collections module not found")
	}
	collectionsModule := collectionsModuleObj.(map[string]interface{})

	mapFunc := collectionsModule["map"].(r2core.BuiltinFunction)
	filterFunc := collectionsModule["filter"].(r2core.BuiltinFunction)
	reduceFunc := collectionsModule["reduce"].(r2core.BuiltinFunction)
	sortFunc := collectionsModule["sort"].(r2core.BuiltinFunction)
	findFunc := collectionsModule["find"].(r2core.BuiltinFunction)
	containsFunc := collectionsModule["contains"].(r2core.BuiltinFunction)

	// Dummy user function for testing: func identity(x) { return x }
	identityFuncCode := "func identity(x) { return x }"
	identityParser := r2core.NewParser(identityFuncCode)
	identityProgram := identityParser.ParseProgram()
	identityProgram.Eval(env) // Evaluate to register the function
	var identityUserFunc *r2core.UserFunction
	{
		val, ok := env.Get("identity")
		if !ok {
			t.Fatal("identity function not found in environment")
		}
		identityUserFunc = val.(*r2core.UserFunction)
	}

	// Dummy comparison function for sort: func lessThan(a, b) { return a < b }
	compFuncCode := "func lessThan(a, b) { return a < b }"
	compParser := r2core.NewParser(compFuncCode)
	compProgram := compParser.ParseProgram()
	compProgram.Eval(env) // Evaluate to register the function
	var compUserFunc *r2core.UserFunction
	{
		val, ok := env.Get("lessThan")
		if !ok {
			t.Fatal("lessThan function not found in environment")
		}
		compUserFunc = val.(*r2core.UserFunction)
	}

	// Dummy sum function for reduce: func sum(a, b) { return a + b }
	sumFuncCode := "func sum(a, b) { return a + b }"
	sumParser := r2core.NewParser(sumFuncCode)
	sumProgram := sumParser.ParseProgram()
	sumProgram.Eval(env) // Evaluate to register the function
	var sumUserFunc *r2core.UserFunction
	{
		val, ok := env.Get("sum")
		if !ok {
			t.Fatal("sum function not found in environment")
		}
		sumUserFunc = val.(*r2core.UserFunction)
	}

	tests := []struct {
		name     string
		function r2core.BuiltinFunction
		args     []interface{}
		expected interface{}
	}{
		{
			"map with numbers",
			mapFunc,
			[]interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				identityUserFunc,
			},
			[]interface{}{1.0, 2.0, 3.0},
		},
		{
			"filter with numbers",
			filterFunc,
			[]interface{}{
				[]interface{}{1.0, 2.0, 3.0, 4.0},
				identityUserFunc,
			},
			[]interface{}{1.0, 2.0, 3.0, 4.0},
		},
		{
			"reduce with numbers",
			reduceFunc,
			[]interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				sumUserFunc,
				0.0,
			},
			6.0,
		},
		{
			"sort with numbers (default)",
			sortFunc,
			[]interface{}{
				[]interface{}{3.0, 1.0, 2.0},
			},
			[]interface{}{1.0, 2.0, 3.0},
		},
		{
			"sort with numbers (custom func)",
			sortFunc,
			[]interface{}{
				[]interface{}{3.0, 1.0, 2.0},
				compUserFunc,
			},
			[]interface{}{1.0, 2.0, 3.0},
		},
		{
			"find with numbers",
			findFunc,
			[]interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				identityUserFunc,
			},
			1.0,
		},
		{
			"contains with numbers",
			containsFunc,
			[]interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				2.0,
			},
			true,
		},
		{
			"contains not found",
			containsFunc,
			[]interface{}{
				[]interface{}{1.0, 2.0, 3.0},
				4.0,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.function(tt.args...)
			// For slice comparisons, need to iterate
			if expectedSlice, ok := tt.expected.([]interface{}); ok {
				if resultArray, ok := result.([]interface{}); ok {
					if len(expectedSlice) != len(resultArray) {
						t.Errorf("expected length %d, got %d", len(expectedSlice), len(resultArray))
						return
					}
					for i := range expectedSlice {
						if expectedSlice[i] != resultArray[i] {
							t.Errorf("expected %v, got %v at index %d", expectedSlice[i], resultArray[i], i)
							return
						}
					}
				} else {
					t.Errorf("expected slice, got %T", result)
				}
			} else if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCollectionsExtendedFunctions(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterCollections(env)

	collectionsModuleObj, ok := env.Get("collections")
	if !ok {
		t.Fatal("collections module not found")
	}
	collectionsModule := collectionsModuleObj.(map[string]interface{})

	indexOfFunc := collectionsModule["indexOf"].(r2core.BuiltinFunction)
	uniqueFunc := collectionsModule["unique"].(r2core.BuiltinFunction)
	compactFunc := collectionsModule["compact"].(r2core.BuiltinFunction)
	flattenFunc := collectionsModule["flatten"].(r2core.BuiltinFunction)
	chunkFunc := collectionsModule["chunk"].(r2core.BuiltinFunction)
	partitionFunc := collectionsModule["partition"].(r2core.BuiltinFunction)
	zipFunc := collectionsModule["zip"].(r2core.BuiltinFunction)
	groupByFunc := collectionsModule["groupBy"].(r2core.BuiltinFunction)
	sortByFunc := collectionsModule["sortBy"].(r2core.BuiltinFunction)

	// isEven(x) { return x % 2 == 0 }
	isEvenCode := "func isEven(x) { return x % 2 == 0 }"
	isEvenParser := r2core.NewParser(isEvenCode)
	isEvenParser.ParseProgram().Eval(env)
	isEvenVal, ok := env.Get("isEven")
	if !ok {
		t.Fatal("isEven function not found in environment")
	}
	isEvenUserFunc := isEvenVal.(*r2core.UserFunction)

	// negate(x) { return -x } used as a sort key extractor
	negateCode := "func negate(x) { return -x }"
	negateParser := r2core.NewParser(negateCode)
	negateParser.ParseProgram().Eval(env)
	negateVal, ok := env.Get("negate")
	if !ok {
		t.Fatal("negate function not found in environment")
	}
	negateUserFunc := negateVal.(*r2core.UserFunction)

	t.Run("indexOf found", func(t *testing.T) {
		result := indexOfFunc([]interface{}{10.0, 20.0, 30.0}, 20.0)
		if result != 1.0 {
			t.Errorf("expected 1, got %v", result)
		}
	})

	t.Run("indexOf not found", func(t *testing.T) {
		result := indexOfFunc([]interface{}{10.0, 20.0, 30.0}, 99.0)
		if result != -1.0 {
			t.Errorf("expected -1, got %v", result)
		}
	})

	t.Run("unique removes duplicates preserving order", func(t *testing.T) {
		result := uniqueFunc([]interface{}{1.0, 2.0, 1.0, 3.0, 2.0}).([]interface{})
		expected := []interface{}{1.0, 2.0, 3.0}
		if len(result) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(result))
		}
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("expected %v at %d, got %v", expected[i], i, result[i])
			}
		}
	})

	t.Run("compact removes nil and falsy", func(t *testing.T) {
		result := compactFunc([]interface{}{1.0, nil, 0.0, "a", false, "", 2.0}).([]interface{})
		expected := []interface{}{1.0, "a", 2.0}
		if len(result) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(result))
		}
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("expected %v at %d, got %v", expected[i], i, result[i])
			}
		}
	})

	t.Run("flatten default depth 1", func(t *testing.T) {
		nested := []interface{}{1.0, []interface{}{2.0, 3.0}, []interface{}{4.0, []interface{}{5.0}}}
		result := flattenFunc(nested).([]interface{})
		expected := []interface{}{1.0, 2.0, 3.0, 4.0, []interface{}{5.0}}
		if len(result) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(result))
		}
	})

	t.Run("flatten depth 2 flattens fully", func(t *testing.T) {
		nested := []interface{}{1.0, []interface{}{2.0, []interface{}{3.0}}}
		result := flattenFunc(nested, 2.0).([]interface{})
		expected := []interface{}{1.0, 2.0, 3.0}
		if len(result) != len(expected) {
			t.Fatalf("expected length %d, got %d", len(expected), len(result))
		}
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("expected %v at %d, got %v", expected[i], i, result[i])
			}
		}
	})

	t.Run("flatten rejects depth beyond safety cap on self-referential array", func(t *testing.T) {
		selfRef := make([]interface{}, 1)
		selfRef[0] = selfRef
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for depth beyond maxFlattenDepth")
			}
		}()
		flattenFunc(selfRef, float64(maxFlattenDepth+1))
	})

	t.Run("flatten allows self-referential array up to the safety cap", func(t *testing.T) {
		selfRef := make([]interface{}, 1)
		selfRef[0] = selfRef
		result := flattenFunc(selfRef, float64(maxFlattenDepth)).([]interface{})
		if len(result) != 1 {
			t.Fatalf("expected length 1, got %d", len(result))
		}
	})

	t.Run("chunk splits into groups", func(t *testing.T) {
		result := chunkFunc([]interface{}{1.0, 2.0, 3.0, 4.0, 5.0}, 2.0).([]interface{})
		if len(result) != 3 {
			t.Fatalf("expected 3 chunks, got %d", len(result))
		}
		last := result[2].([]interface{})
		if len(last) != 1 || last[0] != 5.0 {
			t.Errorf("expected last chunk to be [5], got %v", last)
		}
	})

	t.Run("chunk size zero panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for size 0")
			}
		}()
		chunkFunc([]interface{}{1.0}, 0.0)
	})

	t.Run("partition splits matching and non-matching", func(t *testing.T) {
		result := partitionFunc([]interface{}{1.0, 2.0, 3.0, 4.0}, isEvenUserFunc).([]interface{})
		if len(result) != 2 {
			t.Fatalf("expected 2 groups, got %d", len(result))
		}
		matched := result[0].([]interface{})
		rest := result[1].([]interface{})
		if len(matched) != 2 || len(rest) != 2 {
			t.Errorf("expected 2/2 split, got %d/%d", len(matched), len(rest))
		}
	})

	t.Run("zip pairs elements up to shortest length", func(t *testing.T) {
		result := zipFunc([]interface{}{1.0, 2.0, 3.0}, []interface{}{"a", "b"}).([]interface{})
		if len(result) != 2 {
			t.Fatalf("expected 2 pairs, got %d", len(result))
		}
		pair0 := result[0].([]interface{})
		if pair0[0] != 1.0 || pair0[1] != "a" {
			t.Errorf("expected [1, a], got %v", pair0)
		}
	})

	t.Run("groupBy groups by key function", func(t *testing.T) {
		result := groupByFunc([]interface{}{1.0, 2.0, 3.0, 4.0}, isEvenUserFunc).(map[string]interface{})
		evens := result["true"].([]interface{})
		odds := result["false"].([]interface{})
		if len(evens) != 2 || len(odds) != 2 {
			t.Errorf("expected 2/2 split, got %d/%d", len(evens), len(odds))
		}
	})

	t.Run("sortBy sorts ascending by key", func(t *testing.T) {
		result := sortByFunc([]interface{}{1.0, 3.0, 2.0}, negateUserFunc).([]interface{})
		expected := []interface{}{3.0, 2.0, 1.0}
		for i := range expected {
			if result[i] != expected[i] {
				t.Errorf("expected %v at %d, got %v", expected[i], i, result[i])
			}
		}
	})
}

// TestCollectionsSlice_EndInclusiveOfLastElement guards against an
// off-by-one: end is an exclusive upper bound (matching Go slice-expression
// semantics), so end == len(arr) is valid and means "through the last
// element". The previous "end >= len(arr)" range check rejected that case,
// making it impossible to ever slice through an array's last element.
func TestCollectionsSlice_EndInclusiveOfLastElement(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterCollections(env)
	collectionsModule := mustGetModule(t, env, "collections")
	sliceFunc := collectionsModule["slice"].(r2core.BuiltinFunction)

	arr := []interface{}{1.0, 2.0, 3.0, 4.0, 5.0}

	result := sliceFunc(arr, 0.0, 5.0).([]interface{})
	if len(result) != 5 || result[4] != 5.0 {
		t.Fatalf("slice(arr, 0, 5) should include the last element, got %v", result)
	}

	result = sliceFunc(arr, 3.0, 5.0).([]interface{})
	if len(result) != 2 || result[0] != 4.0 || result[1] != 5.0 {
		t.Fatalf("slice(arr, 3, 5) should return the last two elements, got %v", result)
	}

	result = sliceFunc(arr, 5.0, 5.0).([]interface{})
	if len(result) != 0 {
		t.Fatalf("slice(arr, 5, 5) should return an empty array, got %v", result)
	}

	result = sliceFunc(r2core.InterfaceSlice{1.0, 2.0, 3.0}, 0.0, 3.0).([]interface{})
	if len(result) != 3 {
		t.Fatalf("slice should accept r2core.InterfaceSlice, got %v", result)
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected slice(arr, 0, 6) to panic (end out of range)")
			}
		}()
		sliceFunc(arr, 0.0, 6.0)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected slice(arr, 3, 1) to panic (start > end)")
			}
		}()
		sliceFunc(arr, 3.0, 1.0)
	}()
}

func mustGetModule(t *testing.T, env *r2core.Environment, name string) map[string]interface{} {
	t.Helper()
	obj, ok := env.Get(name)
	if !ok {
		t.Fatalf("module %q not found", name)
	}
	module, ok := obj.(map[string]interface{})
	if !ok {
		t.Fatalf("%q is not a module map", name)
	}
	return module
}
