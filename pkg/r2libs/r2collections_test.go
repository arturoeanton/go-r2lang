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
