package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestRandInterfaceSlice(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterRand(env)

	randModuleObj, ok := env.Get("rand")
	if !ok {
		t.Fatal("rand module not found")
	}
	randModule := randModuleObj.(map[string]interface{})

	randChoiceFunc := randModule["randChoice"].(r2core.BuiltinFunction)
	shuffleFunc := randModule["shuffle"].(r2core.BuiltinFunction)
	sampleFunc := randModule["sample"].(r2core.BuiltinFunction)

	// r2core.InterfaceSlice is what array methods like .map()/.filter()
	// return; randChoice/shuffle/sample must accept it just like a plain
	// []interface{} literal, not panic with a type error.
	t.Run("randChoice accepts InterfaceSlice", func(t *testing.T) {
		arr := r2core.InterfaceSlice{1.0, 2.0, 3.0}
		result := randChoiceFunc(arr)
		found := false
		for _, v := range arr {
			if v == result {
				found = true
			}
		}
		if !found {
			t.Errorf("randChoice: result %v not found in source array", result)
		}
	})

	t.Run("shuffle accepts InterfaceSlice and mutates in place", func(t *testing.T) {
		arr := r2core.InterfaceSlice{1.0, 2.0, 3.0, 4.0, 5.0}
		shuffleFunc(arr)
		sum := 0.0
		for _, v := range arr {
			sum += v.(float64)
		}
		if sum != 15.0 {
			t.Errorf("shuffle: expected same elements summing to 15, got sum %v (%v)", sum, arr)
		}
	})

	t.Run("sample accepts InterfaceSlice", func(t *testing.T) {
		arr := r2core.InterfaceSlice{1.0, 2.0, 3.0, 4.0, 5.0}
		result := sampleFunc(arr, 3.0).([]interface{})
		if len(result) != 3 {
			t.Errorf("sample: expected 3 elements, got %d", len(result))
		}
	})
}
