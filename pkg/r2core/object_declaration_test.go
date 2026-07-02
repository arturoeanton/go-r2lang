package r2core

import "testing"

// TestObjectDeclaration_FieldDefaultInitializer guards against a bug where
// a class field's declared default value (e.g. "let value = 0;") was
// silently discarded — every instance's field started as nil regardless of
// its initializer, until something explicitly assigned it. Any class
// relying on a default field value without an explicit constructor
// assignment was silently broken.
func TestObjectDeclaration_FieldDefaultInitializer(t *testing.T) {
	code := `
class Counter {
    let value = 0;
    let label = "counter";
    let items = [1, 2, 3];

    func getValue() { return this.value; }
    func getLabel() { return this.label; }
    func getItems() { return this.items; }
}
let c = Counter();
let value = c.getValue();
let label = c.getLabel();
let items = c.getItems();
`
	env := NewEnvironment()
	program := NewParser(code).ParseProgram()
	program.Eval(env)

	value, _ := env.Get("value")
	if value != float64(0) {
		t.Errorf("expected default field value 0, got %v (%T)", value, value)
	}

	label, _ := env.Get("label")
	if label != "counter" {
		t.Errorf("expected default field label %q, got %v (%T)", "counter", label, label)
	}

	itemsVal, _ := env.Get("items")
	items, ok := itemsVal.([]interface{})
	if !ok || len(items) != 3 {
		t.Errorf("expected default field items to be a 3-element array, got %v (%T)", itemsVal, itemsVal)
	}
}

// TestObjectDeclaration_FieldWithoutInitializerStaysNil confirms a field
// declared without an initializer ("let x;") still starts as nil, so the
// fix above only changes behavior for fields that actually declare a
// default value.
func TestObjectDeclaration_FieldWithoutInitializerStaysNil(t *testing.T) {
	code := `
class Empty {
    let value;
    func getValue() { return this.value; }
}
let e = Empty();
let value = e.getValue();
`
	env := NewEnvironment()
	program := NewParser(code).ParseProgram()
	program.Eval(env)

	value, _ := env.Get("value")
	if value != nil {
		t.Errorf("expected field with no initializer to stay nil, got %v (%T)", value, value)
	}
}

// TestObjectDeclaration_ConstructorOverridesDefault confirms an explicit
// constructor assignment still overrides the field's declared default, so
// this fix doesn't regress the (already-working) explicit-assignment path.
func TestObjectDeclaration_ConstructorOverridesDefault(t *testing.T) {
	code := `
class WithConstructor {
    let value = 0;
    func constructor(v) { this.value = v; }
    func getValue() { return this.value; }
}
let w = WithConstructor(42);
let value = w.getValue();
`
	env := NewEnvironment()
	program := NewParser(code).ParseProgram()
	program.Eval(env)

	value, _ := env.Get("value")
	if value != float64(42) {
		t.Errorf("expected constructor to override default, got %v (%T)", value, value)
	}
}
