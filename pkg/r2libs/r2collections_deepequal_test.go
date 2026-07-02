package r2libs

import (
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func collectionsModuleForTest(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterCollections(env)
	modObj, ok := env.Get("collections")
	if !ok {
		t.Fatal("collections module not found")
	}
	return modObj.(map[string]interface{})
}

func TestDeepEqualPrimitivesAndNesting(t *testing.T) {
	mod := collectionsModuleForTest(t)
	deepEqual := mod["deepEqual"].(r2core.BuiltinFunction)

	cases := []struct {
		name     string
		a, b     interface{}
		expected bool
	}{
		{"equal numbers", 1.0, 1.0, true},
		{"different numbers", 1.0, 2.0, false},
		{"equal strings", "hi", "hi", true},
		{"different strings", "hi", "bye", false},
		{"nil vs nil", nil, nil, true},
		{"nil vs value", nil, 1.0, false},
		{
			"equal nested maps",
			map[string]interface{}{"a": 1.0, "b": []interface{}{1.0, 2.0, map[string]interface{}{"c": true}}},
			map[string]interface{}{"a": 1.0, "b": []interface{}{1.0, 2.0, map[string]interface{}{"c": true}}},
			true,
		},
		{
			"different nested maps",
			map[string]interface{}{"a": 1.0, "b": []interface{}{1.0, 2.0}},
			map[string]interface{}{"a": 1.0, "b": []interface{}{1.0, 3.0}},
			false,
		},
		{
			"different map lengths",
			map[string]interface{}{"a": 1.0},
			map[string]interface{}{"a": 1.0, "b": 2.0},
			false,
		},
		{
			"array vs InterfaceSlice equal",
			[]interface{}{1.0, 2.0, map[string]interface{}{"x": 1.0}},
			r2core.InterfaceSlice{1.0, 2.0, map[string]interface{}{"x": 1.0}},
			true,
		},
		{
			"array vs InterfaceSlice different",
			[]interface{}{1.0, 2.0},
			r2core.InterfaceSlice{1.0, 3.0},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := deepEqual(tc.a, tc.b).(bool)
			if got != tc.expected {
				t.Errorf("deepEqual(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}

func TestDeepEqualSelfReferentialDoesNotStackOverflow(t *testing.T) {
	mod := collectionsModuleForTest(t)
	deepEqual := mod["deepEqual"].(r2core.BuiltinFunction)

	a := make([]interface{}, 1)
	a[0] = a
	b := make([]interface{}, 1)
	b[0] = b

	done := make(chan bool, 1)
	go func() {
		done <- deepEqual(a, b).(bool)
	}()

	select {
	case result := <-done:
		if !result {
			t.Errorf("deepEqual on identically-shaped cyclic arrays = false, want true")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("deepEqual did not terminate on self-referential arrays")
	}

	mapA := map[string]interface{}{}
	mapA["self"] = mapA
	mapB := map[string]interface{}{}
	mapB["self"] = mapB

	done2 := make(chan bool, 1)
	go func() {
		done2 <- deepEqual(mapA, mapB).(bool)
	}()

	select {
	case result := <-done2:
		if !result {
			t.Errorf("deepEqual on identically-shaped cyclic maps = false, want true")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("deepEqual did not terminate on self-referential maps")
	}

	mapC := map[string]interface{}{}
	mapC["self"] = mapC
	notCyclic := map[string]interface{}{"self": map[string]interface{}{"other": 1.0}}

	done3 := make(chan bool, 1)
	go func() {
		done3 <- deepEqual(mapC, notCyclic).(bool)
	}()

	select {
	case result := <-done3:
		if result {
			t.Errorf("deepEqual(cyclic, non-cyclic) = true, want false")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("deepEqual did not terminate comparing cyclic vs non-cyclic maps")
	}
}

func TestDeepCloneAliasesDeepCopy(t *testing.T) {
	mod := collectionsModuleForTest(t)
	deepClone := mod["deepClone"].(r2core.BuiltinFunction)
	deepEqual := mod["deepEqual"].(r2core.BuiltinFunction)

	original := r2core.InterfaceSlice{
		1.0,
		map[string]interface{}{"nested": []interface{}{1.0, 2.0, r2core.InterfaceSlice{3.0, 4.0}}},
	}

	cloned := deepClone(original)

	if !deepEqual([]interface{}(original), cloned).(bool) {
		t.Fatal("deepClone result is not deepEqual to original")
	}

	clonedSlice := cloned.([]interface{})
	clonedInner := clonedSlice[1].(map[string]interface{})
	clonedInner["nested"].([]interface{})[0] = 999.0

	originalInner := original[1].(map[string]interface{})
	if originalInner["nested"].([]interface{})[0] == 999.0 {
		t.Fatal("mutating clone mutated original: deepClone is not a deep copy")
	}
}
