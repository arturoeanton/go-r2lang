package assertions

import (
	"testing"
)

func TestAssert_Equals(t *testing.T) {
	assert := NewAssert("test")

	// Test equal values
	assert.Equals(5, 5)
	assert.Equals("hello", "hello")
	assert.Equals(true, true)

	// Test unequal values should panic
	defer func() {
		if r := recover(); r != nil {
			if ae, ok := r.(*AssertionError); ok {
				if ae.TestName != "test" {
					t.Errorf("Expected test name 'test', got '%s'", ae.TestName)
				}
			} else {
				t.Errorf("Expected AssertionError, got %T", r)
			}
		} else {
			t.Error("Expected panic for unequal values")
		}
	}()

	assert.Equals(5, 10) // This should panic
}

func TestAssert_NotEquals(t *testing.T) {
	assert := NewAssert("test")

	// Test unequal values
	assert.NotEquals(5, 10)
	assert.NotEquals("hello", "world")
	assert.NotEquals(true, false)

	// Test equal values should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for equal values")
		}
	}()

	assert.NotEquals(5, 5) // This should panic
}

func TestAssert_True(t *testing.T) {
	assert := NewAssert("test")

	// Test truthy values
	assert.True(true)
	assert.True(1)
	assert.True("hello")
	assert.True([]interface{}{1, 2, 3})

	// Test falsy value should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for false value")
		}
	}()

	assert.True(false) // This should panic
}

func TestAssert_False(t *testing.T) {
	assert := NewAssert("test")

	// Test falsy values
	assert.False(false)
	assert.False(0)
	assert.False("")
	assert.False([]interface{}{})

	// Test truthy value should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for true value")
		}
	}()

	assert.False(true) // This should panic
}

func TestAssert_Nil(t *testing.T) {
	assert := NewAssert("test")

	// Test nil value
	assert.Nil(nil)

	// Test non-nil value should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-nil value")
		}
	}()

	assert.Nil("not nil") // This should panic
}

func TestAssert_NotNil(t *testing.T) {
	assert := NewAssert("test")

	// Test non-nil values
	assert.NotNil("hello")
	assert.NotNil(5)
	assert.NotNil(true)

	// Test nil value should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil value")
		}
	}()

	assert.NotNil(nil) // This should panic
}

func TestAssert_Contains(t *testing.T) {
	assert := NewAssert("test")

	// Test string contains
	assert.Contains("hello world", "world")
	assert.Contains("testing", "test")

	// Test string doesn't contain should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for string that doesn't contain substring")
		}
	}()

	assert.Contains("hello", "world") // This should panic
}

func TestAssert_NotContains(t *testing.T) {
	assert := NewAssert("test")

	// Test string doesn't contain
	assert.NotContains("hello", "world")
	assert.NotContains("testing", "xyz")

	// Test string contains should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for string that contains substring")
		}
	}()

	assert.NotContains("hello world", "world") // This should panic
}

func TestAssert_Greater(t *testing.T) {
	assert := NewAssert("test")

	// Test greater values
	assert.Greater(10, 5)
	assert.Greater(5.5, 5.0)

	// Test not greater should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for not greater value")
		}
	}()

	assert.Greater(5, 10) // This should panic
}

func TestAssert_Less(t *testing.T) {
	assert := NewAssert("test")

	// Test less values
	assert.Less(5, 10)
	assert.Less(5.0, 5.5)

	// Test not less should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for not less value")
		}
	}()

	assert.Less(10, 5) // This should panic
}

func TestAssert_HasLength(t *testing.T) {
	assert := NewAssert("test")

	// Test correct lengths
	assert.HasLength("hello", 5)
	assert.HasLength([]interface{}{1, 2, 3}, 3)
	assert.HasLength(map[string]interface{}{"a": 1, "b": 2}, 2)

	// Test wrong length should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for wrong length")
		}
	}()

	assert.HasLength("hello", 10) // This should panic
}

func TestAssert_Empty(t *testing.T) {
	assert := NewAssert("test")

	// Test empty collections
	assert.Empty("")
	assert.Empty([]interface{}{})
	assert.Empty(map[string]interface{}{})

	// Test non-empty should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-empty collection")
		}
	}()

	assert.Empty("hello") // This should panic
}

func TestAssert_NotEmpty(t *testing.T) {
	assert := NewAssert("test")

	// Test non-empty collections
	assert.NotEmpty("hello")
	assert.NotEmpty([]interface{}{1})
	assert.NotEmpty(map[string]interface{}{"a": 1})

	// Test empty should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty collection")
		}
	}()

	assert.NotEmpty("") // This should panic
}

func TestAssert_Panics(t *testing.T) {
	assert := NewAssert("test")

	// Test function that panics
	assert.Panics(func() {
		panic("test panic")
	})

	// Test function that doesn't panic should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for function that doesn't panic")
		}
	}()

	assert.Panics(func() {
		// This function doesn't panic
	})
}

func TestAssert_NotPanics(t *testing.T) {
	assert := NewAssert("test")

	// Test function that doesn't panic
	assert.NotPanics(func() {
		// This function doesn't panic
	})

	// Test function that panics should panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for function that panics")
		}
	}()

	assert.NotPanics(func() {
		panic("test panic")
	})
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{true, true},
		{false, false},
		{1, true},
		{0, false},
		{-1, true},
		{"hello", true},
		{"", false},
		{[]interface{}{1}, true},
		{[]interface{}{}, false},
		{map[string]interface{}{"a": 1}, true},
		{map[string]interface{}{}, false},
		{nil, false},
	}

	for _, test := range tests {
		result := isTruthy(test.value)
		if result != test.expected {
			t.Errorf("isTruthy(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestGetLength(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected int
	}{
		{"hello", 5},
		{"", 0},
		{[]interface{}{1, 2, 3}, 3},
		{[]interface{}{}, 0},
		{map[string]interface{}{"a": 1, "b": 2}, 2},
		{map[string]interface{}{}, 0},
		{nil, 0},
	}

	for _, test := range tests {
		result := getLength(test.value)
		if result != test.expected {
			t.Errorf("getLength(%v) = %v, expected %v", test.value, result, test.expected)
		}
	}
}

func TestCompareNumbers(t *testing.T) {
	tests := []struct {
		a, b     interface{}
		expected int
	}{
		{5, 3, 1},
		{3, 5, -1},
		{5, 5, 0},
		{5.5, 5.0, 1},
		{5.0, 5.5, -1},
		{5.0, 5.0, 0},
	}

	for _, test := range tests {
		result := compareNumbers(test.a, test.b)
		if result != test.expected {
			t.Errorf("compareNumbers(%v, %v) = %v, expected %v", test.a, test.b, result, test.expected)
		}
	}
}
