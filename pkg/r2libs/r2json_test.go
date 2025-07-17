package r2libs

import (
	"strings"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestJSONParse(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	parseFunc := module["parse"].(r2core.BuiltinFunction)

	// Test parsing object
	result := parseFunc(`{"name": "John", "age": 30}`)
	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected object from JSON parse")
	}
	if obj["name"] != "John" || obj["age"] != 30.0 {
		t.Error("JSON parse did not correctly parse object")
	}

	// Test parsing array
	result = parseFunc(`[1, 2, 3]`)
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatal("Expected array from JSON parse")
	}
	if len(arr) != 3 || arr[0] != 1.0 || arr[1] != 2.0 || arr[2] != 3.0 {
		t.Error("JSON parse did not correctly parse array")
	}
}

func TestJSONStringify(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	stringifyFunc := module["stringify"].(r2core.BuiltinFunction)

	// Test stringifying object
	obj := map[string]interface{}{
		"name": "John",
		"age":  30.0,
	}
	result := stringifyFunc(obj)
	jsonStr, ok := result.(string)
	if !ok {
		t.Fatal("Expected string from JSON stringify")
	}
	if jsonStr != `{"age":30,"name":"John"}` {
		t.Errorf("JSON stringify returned unexpected result: %s", jsonStr)
	}

	// Test stringifying array
	arr := []interface{}{1.0, 2.0, 3.0}
	result = stringifyFunc(arr)
	jsonStr, ok = result.(string)
	if !ok {
		t.Fatal("Expected string from JSON stringify")
	}
	if jsonStr != `[1,2,3]` {
		t.Errorf("JSON stringify returned unexpected result: %s", jsonStr)
	}
}

func TestJSONValidate(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	validateFunc := module["validate"].(r2core.BuiltinFunction)

	// Test valid JSON
	result := validateFunc(`{"name": "John", "age": 30}`)
	if result != true {
		t.Error("Expected valid JSON to return true")
	}

	// Test invalid JSON
	result = validateFunc(`{"name": "John", "age": 30`)
	if result != false {
		t.Error("Expected invalid JSON to return false")
	}
}

func TestJSONGetKeys(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	getKeysFunc := module["getKeys"].(r2core.BuiltinFunction)

	result := getKeysFunc(`{"name": "John", "age": 30}`)
	keys, ok := result.([]interface{})
	if !ok {
		t.Fatal("Expected array from getKeys")
	}
	if len(keys) != 2 {
		t.Error("Expected 2 keys")
	}

	keysSet := make(map[string]bool)
	for _, key := range keys {
		keysSet[key.(string)] = true
	}

	if !keysSet["name"] || !keysSet["age"] {
		t.Error("Expected keys 'name' and 'age'")
	}
}

func TestJSONGetValue(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	getValueFunc := module["getValue"].(r2core.BuiltinFunction)

	result := getValueFunc(`{"name": "John", "age": 30}`, "name")
	if result != "John" {
		t.Error("Expected getValue to return 'John'")
	}

	result = getValueFunc(`{"name": "John", "age": 30}`, "age")
	if result != 30.0 {
		t.Error("Expected getValue to return 30")
	}

	result = getValueFunc(`{"name": "John", "age": 30}`, "missing")
	if result != nil {
		t.Error("Expected getValue to return nil for missing key")
	}
}

func TestJSONSetValue(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	setValueFunc := module["setValue"].(r2core.BuiltinFunction)
	parseFunc := module["parse"].(r2core.BuiltinFunction)

	result := setValueFunc(`{"name": "John", "age": 30}`, "city", "New York")
	updatedJSON, ok := result.(string)
	if !ok {
		t.Fatal("Expected string from setValue")
	}

	parsed := parseFunc(updatedJSON)
	obj, ok := parsed.(map[string]interface{})
	if !ok {
		t.Fatal("Expected object from parsed JSON")
	}

	if obj["city"] != "New York" {
		t.Error("Expected city to be 'New York'")
	}
}

func TestJSONHasKey(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	hasKeyFunc := module["hasKey"].(r2core.BuiltinFunction)

	result := hasKeyFunc(`{"name": "John", "age": 30}`, "name")
	if result != true {
		t.Error("Expected hasKey to return true for existing key")
	}

	result = hasKeyFunc(`{"name": "John", "age": 30}`, "missing")
	if result != false {
		t.Error("Expected hasKey to return false for missing key")
	}
}

func TestJSONMerge(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	mergeFunc := module["merge"].(r2core.BuiltinFunction)
	parseFunc := module["parse"].(r2core.BuiltinFunction)

	result := mergeFunc(`{"name": "John"}`, `{"age": 30}`)
	merged, ok := result.(string)
	if !ok {
		t.Fatal("Expected string from merge")
	}

	parsed := parseFunc(merged)
	obj, ok := parsed.(map[string]interface{})
	if !ok {
		t.Fatal("Expected object from parsed JSON")
	}

	if obj["name"] != "John" || obj["age"] != 30.0 {
		t.Error("Expected merged object to contain both name and age")
	}
}

func TestJSONQuery(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	queryFunc := module["query"].(r2core.BuiltinFunction)

	jsonStr := `{"user": {"name": "John", "details": {"age": 30}}}`

	result := queryFunc(jsonStr, "user.name")
	if result != "John" {
		t.Error("Expected query to return 'John'")
	}

	result = queryFunc(jsonStr, "user.details.age")
	if result != 30.0 {
		t.Error("Expected query to return 30")
	}
}

func TestJSONType(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	typeFunc := module["type"].(r2core.BuiltinFunction)

	result := typeFunc(`{"name": "John"}`)
	if result != "object" {
		t.Error("Expected type to return 'object'")
	}

	result = typeFunc(`[1, 2, 3]`)
	if result != "array" {
		t.Error("Expected type to return 'array'")
	}

	result = typeFunc(`"hello"`)
	if result != "string" {
		t.Error("Expected type to return 'string'")
	}

	result = typeFunc(`42`)
	if result != "number" {
		t.Error("Expected type to return 'number'")
	}

	result = typeFunc(`true`)
	if result != "boolean" {
		t.Error("Expected type to return 'boolean'")
	}

	result = typeFunc(`null`)
	if result != "null" {
		t.Error("Expected type to return 'null'")
	}
}

func TestJSONSize(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	sizeFunc := module["size"].(r2core.BuiltinFunction)

	result := sizeFunc(`{"name": "John", "age": 30}`)
	if result != 2.0 {
		t.Error("Expected size to return 2 for object with 2 keys")
	}

	result = sizeFunc(`[1, 2, 3, 4]`)
	if result != 4.0 {
		t.Error("Expected size to return 4 for array with 4 elements")
	}

	result = sizeFunc(`"hello"`)
	if result != 5.0 {
		t.Error("Expected size to return 5 for string with 5 characters")
	}
}

func TestJSONFlatten(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	flattenFunc := module["flatten"].(r2core.BuiltinFunction)
	parseFunc := module["parse"].(r2core.BuiltinFunction)

	result := flattenFunc(`{"user": {"name": "John", "age": 30}}`)
	flattened, ok := result.(string)
	if !ok {
		t.Fatal("Expected string from flatten")
	}

	parsed := parseFunc(flattened)
	obj, ok := parsed.(map[string]interface{})
	if !ok {
		t.Fatal("Expected object from parsed JSON")
	}

	if obj["user.name"] != "John" || obj["user.age"] != 30.0 {
		t.Error("Expected flattened object to have dotted keys")
	}
}

func TestJSONPretty(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	prettyFunc := module["pretty"].(r2core.BuiltinFunction)

	result := prettyFunc(`{"name":"John","age":30}`)
	pretty, ok := result.(string)
	if !ok {
		t.Fatal("Expected string from pretty")
	}

	if !strings.Contains(pretty, "\n") {
		t.Error("Expected pretty JSON to contain newlines")
	}
}

func TestJSONMinify(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterJSON(env)

	jsonModule, _ := env.Get("json")
	module := jsonModule.(map[string]interface{})
	minifyFunc := module["minify"].(r2core.BuiltinFunction)

	result := minifyFunc(`{
		"name": "John",
		"age": 30
	}`)
	minified, ok := result.(string)
	if !ok {
		t.Fatal("Expected string from minify")
	}

	if strings.Contains(minified, "\n") || strings.Contains(minified, "  ") {
		t.Error("Expected minified JSON to not contain whitespace")
	}
}
