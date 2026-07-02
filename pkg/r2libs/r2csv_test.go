package r2libs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func csvModule(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterCSV(env)
	mod, ok := env.Get("csv")
	if !ok {
		t.Fatal("csv module not registered")
	}
	return mod.(map[string]interface{})
}

func TestCSVParseStripsBOM(t *testing.T) {
	mod := csvModule(t)
	parseFunc := mod["parse"].(r2core.BuiltinFunction)

	csvStr := "\uFEFFID,Name\n1,Alice\n"
	result := parseFunc(csvStr)
	rows, ok := result.([]interface{})
	if !ok || len(rows) != 1 {
		t.Fatalf("expected 1 row, got %#v", result)
	}
	row, ok := rows[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected row to be a map, got %#v", rows[0])
	}
	if _, exists := row["ID"]; !exists {
		t.Errorf("expected header 'ID' without BOM contamination, got keys: %#v", row)
	}
	if v := row["ID"]; v != 1.0 {
		t.Errorf("expected ID=1, got %#v", v)
	}
}

func TestCSVReadFileStripsBOM(t *testing.T) {
	mod := csvModule(t)
	readFileFunc := mod["readFile"].(r2core.BuiltinFunction)

	dir := t.TempDir()
	path := filepath.Join(dir, "bom.csv")
	content := "\xEF\xBB\xBFID,Name\n1,Alice\n2,Bob\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result := readFileFunc(path)
	rows, ok := result.([]interface{})
	if !ok || len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %#v", result)
	}
	row, ok := rows[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected row to be a map, got %#v", rows[0])
	}
	if v, exists := row["ID"]; !exists || v != 1.0 {
		t.Errorf("expected header 'ID'=1 without BOM contamination, got %#v", row)
	}
}

func TestCSVStringifyHeaderOrderIsDeterministic(t *testing.T) {
	mod := csvModule(t)
	stringifyFunc := mod["stringify"].(r2core.BuiltinFunction)

	data := []interface{}{
		map[string]interface{}{"a": 1.0, "b": 2.0, "c": 3.0, "d": 4.0, "e": 5.0, "f": 6.0, "g": 7.0, "h": 8.0},
	}

	const expected = "a,b,c,d,e,f,g,h\n1,2,3,4,5,6,7,8\n"
	for i := 0; i < 10; i++ {
		result := stringifyFunc(data)
		out, ok := result.(string)
		if !ok {
			t.Fatalf("expected string result, got %#v", result)
		}
		if out != expected {
			t.Fatalf("run %d: header order not deterministic, got %q, want %q", i, out, expected)
		}
	}
}

func TestCSVGetHeadersIsDeterministic(t *testing.T) {
	mod := csvModule(t)
	getHeadersFunc := mod["getHeaders"].(r2core.BuiltinFunction)

	data := []interface{}{
		map[string]interface{}{"z": 1.0, "y": 2.0, "x": 3.0, "w": 4.0},
	}

	result := getHeadersFunc(data)
	headers, ok := result.([]interface{})
	if !ok || len(headers) != 4 {
		t.Fatalf("expected 4 headers, got %#v", result)
	}
	expected := []string{"w", "x", "y", "z"}
	for i, h := range headers {
		if h != expected[i] {
			t.Errorf("headers not sorted deterministically: got %#v, want %#v", headers, expected)
			break
		}
	}
}
