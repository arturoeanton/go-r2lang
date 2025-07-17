package r2libs

import (
	"os"
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func TestDatabaseFunctions(t *testing.T) {
	// Create a test environment
	env := r2core.NewEnvironment()
	RegisterDB(env)

	// Test database connection
	t.Run("TestDBConnect", func(t *testing.T) {
		// Test SQLite connection (in-memory)
		connIdFuncVal, ok := env.Get("dbConnect")
		if !ok {
			t.Fatal("dbConnect function not found")
		}
		connIdFunc := connIdFuncVal.(r2core.BuiltinFunction)
		result := connIdFunc("sqlite3", ":memory:")

		if result == nil {
			t.Fatal("dbConnect should return a connection ID")
		}

		connId := result.(string)
		if connId == "" {
			t.Fatal("dbConnect should return a non-empty connection ID")
		}
	})

	t.Run("TestDBConnectInvalidDriver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("dbConnect should panic with invalid driver")
			}
		}()

		connIdFuncVal, _ := env.Get("dbConnect")
		connIdFunc := connIdFuncVal.(r2core.BuiltinFunction)
		connIdFunc("invalid_driver", ":memory:")
	})

	t.Run("TestDBOperations", func(t *testing.T) {
		// Connect to SQLite
		connIdFuncVal, _ := env.Get("dbConnect")
		connIdFunc := connIdFuncVal.(r2core.BuiltinFunction)
		connId := connIdFunc("sqlite3", ":memory:").(string)

		// Test dbPing
		pingFuncVal, _ := env.Get("dbPing")
		pingFunc := pingFuncVal.(r2core.BuiltinFunction)
		pingResult := pingFunc(connId)
		if pingResult != true {
			t.Fatal("dbPing should return true for valid connection")
		}

		// Test table creation
		execFuncVal, _ := env.Get("dbExec")
		execFunc := execFuncVal.(r2core.BuiltinFunction)
		createTableSQL := `CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE,
			age INTEGER
		)`

		result := execFunc(connId, createTableSQL)
		if result.(int64) != 0 {
			t.Fatal("CREATE TABLE should affect 0 rows")
		}

		// Test insert with dbLastInsertId
		lastInsertIdFuncVal, _ := env.Get("dbLastInsertId")
		lastInsertIdFunc := lastInsertIdFuncVal.(r2core.BuiltinFunction)
		insertSQL := "INSERT INTO users (name, email, age) VALUES (?, ?, ?)"
		lastId := lastInsertIdFunc(connId, insertSQL, "John Doe", "john@example.com", 30)
		if lastId.(int64) != 1 {
			t.Fatalf("First insert should have ID 1, got %v", lastId)
		}

		// Test multiple inserts
		execFunc(connId, insertSQL, "Jane Smith", "jane@example.com", 25)
		execFunc(connId, insertSQL, "Bob Johnson", "bob@example.com", 35)

		// Test query
		queryFuncVal, _ := env.Get("dbQuery")
		queryFunc := queryFuncVal.(r2core.BuiltinFunction)
		selectSQL := "SELECT id, name, email, age FROM users ORDER BY id"
		queryResult := queryFunc(connId, selectSQL)

		rows := queryResult.([]interface{})
		if len(rows) != 3 {
			t.Fatalf("Expected 3 rows, got %d", len(rows))
		}

		// Verify first row
		firstRow := rows[0].(map[string]interface{})
		if firstRow["name"] != "John Doe" {
			t.Fatalf("Expected name 'John Doe', got %v", firstRow["name"])
		}
		if firstRow["email"] != "john@example.com" {
			t.Fatalf("Expected email 'john@example.com', got %v", firstRow["email"])
		}

		// Test parameterized query
		paramSQL := "SELECT * FROM users WHERE age > ? ORDER BY age"
		paramResult := queryFunc(connId, paramSQL, 28)
		paramRows := paramResult.([]interface{})
		if len(paramRows) != 2 {
			t.Fatalf("Expected 2 rows with age > 28, got %d", len(paramRows))
		}

		// Test update
		updateSQL := "UPDATE users SET age = ? WHERE name = ?"
		updateResult := execFunc(connId, updateSQL, 31, "John Doe")
		if updateResult.(int64) != 1 {
			t.Fatalf("UPDATE should affect 1 row, got %v", updateResult)
		}

		// Verify update
		verifySQL := "SELECT age FROM users WHERE name = ?"
		verifyResult := queryFunc(connId, verifySQL, "John Doe")
		verifyRows := verifyResult.([]interface{})
		if len(verifyRows) != 1 {
			t.Fatal("Expected 1 row after update")
		}
		age := verifyRows[0].(map[string]interface{})["age"]
		if age.(int64) != 31 {
			t.Fatalf("Expected age 31, got %v", age)
		}

		// Test delete
		deleteSQL := "DELETE FROM users WHERE name = ?"
		deleteResult := execFunc(connId, deleteSQL, "Bob Johnson")
		if deleteResult.(int64) != 1 {
			t.Fatalf("DELETE should affect 1 row, got %v", deleteResult)
		}

		// Verify deletion
		countSQL := "SELECT COUNT(*) as count FROM users"
		countResult := queryFunc(connId, countSQL)
		countRows := countResult.([]interface{})
		count := countRows[0].(map[string]interface{})["count"]
		if count.(int64) != 2 {
			t.Fatalf("Expected 2 users remaining, got %v", count)
		}

		// Test dbClose
		closeFuncVal, _ := env.Get("dbClose")
		closeFunc := closeFuncVal.(r2core.BuiltinFunction)
		closeResult := closeFunc(connId)
		if closeResult != true {
			t.Fatal("dbClose should return true")
		}
	})

	t.Run("TestDBUtilities", func(t *testing.T) {
		// Test dbEscape
		escapeFuncVal, _ := env.Get("dbEscape")
		escapeFunc := escapeFuncVal.(r2core.BuiltinFunction)
		testString := "O'Connor's Data"
		escaped := escapeFunc(testString)
		expected := "O''Connor''s Data"
		if escaped != expected {
			t.Fatalf("Expected '%s', got '%s'", expected, escaped)
		}

		// Test dbGetConnections with active connections
		connIdFuncVal, _ := env.Get("dbConnect")
		connIdFunc := connIdFuncVal.(r2core.BuiltinFunction)
		connId1 := connIdFunc("sqlite3", ":memory:").(string)
		connId2 := connIdFunc("sqlite3", ":memory:").(string)

		getConnsFuncVal, _ := env.Get("dbGetConnections")
		getConnsFunc := getConnsFuncVal.(r2core.BuiltinFunction)
		connections := getConnsFunc().([]interface{})

		if len(connections) < 2 {
			t.Fatalf("Expected at least 2 connections, got %d", len(connections))
		}

		// Clean up
		closeFuncVal, _ := env.Get("dbClose")
		closeFunc := closeFuncVal.(r2core.BuiltinFunction)
		closeFunc(connId1)
		closeFunc(connId2)
	})

	t.Run("TestDBErrors", func(t *testing.T) {
		// Test operations on non-existent connection
		queryFuncVal, _ := env.Get("dbQuery")
		queryFunc := queryFuncVal.(r2core.BuiltinFunction)

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("dbQuery should panic with invalid connection")
			}
		}()

		queryFunc("invalid_conn", "SELECT 1")
	})

	t.Run("TestDBFileOperations", func(t *testing.T) {
		// Test with a file-based SQLite database
		testDB := "/tmp/test_r2lang.db"
		defer os.Remove(testDB) // Clean up after test

		connIdFuncVal, _ := env.Get("dbConnect")
		connIdFunc := connIdFuncVal.(r2core.BuiltinFunction)
		connId := connIdFunc("sqlite3", testDB).(string)

		execFuncVal, _ := env.Get("dbExec")
		execFunc := execFuncVal.(r2core.BuiltinFunction)
		queryFuncVal, _ := env.Get("dbQuery")
		queryFunc := queryFuncVal.(r2core.BuiltinFunction)

		// Create and populate table
		execFunc(connId, "CREATE TABLE test_table (id INTEGER PRIMARY KEY, value TEXT)")
		execFunc(connId, "INSERT INTO test_table (value) VALUES (?)", "test_value")

		// Query data
		result := queryFunc(connId, "SELECT * FROM test_table")
		rows := result.([]interface{})
		if len(rows) != 1 {
			t.Fatalf("Expected 1 row, got %d", len(rows))
		}

		row := rows[0].(map[string]interface{})
		if row["value"] != "test_value" {
			t.Fatalf("Expected 'test_value', got %v", row["value"])
		}

		closeFuncVal, _ := env.Get("dbClose")
		closeFunc := closeFuncVal.(r2core.BuiltinFunction)
		closeFunc(connId)

		// Verify file was created
		if _, err := os.Stat(testDB); os.IsNotExist(err) {
			t.Fatal("Database file should exist")
		}
	})

	t.Run("TestDBTypesAndConversions", func(t *testing.T) {
		connIdFuncVal, _ := env.Get("dbConnect")
		connIdFunc := connIdFuncVal.(r2core.BuiltinFunction)
		connId := connIdFunc("sqlite3", ":memory:").(string)

		execFuncVal, _ := env.Get("dbExec")
		execFunc := execFuncVal.(r2core.BuiltinFunction)
		queryFuncVal, _ := env.Get("dbQuery")
		queryFunc := queryFuncVal.(r2core.BuiltinFunction)

		// Create table with different data types
		createSQL := `CREATE TABLE type_test (
			id INTEGER PRIMARY KEY,
			int_val INTEGER,
			real_val REAL,
			text_val TEXT,
			bool_val INTEGER
		)`
		execFunc(connId, createSQL)

		// Insert test data
		insertSQL := "INSERT INTO type_test (int_val, real_val, text_val, bool_val) VALUES (?, ?, ?, ?)"
		execFunc(connId, insertSQL, 42, 3.14, "hello world", 1)

		// Query and verify types
		result := queryFunc(connId, "SELECT * FROM type_test")
		rows := result.([]interface{})
		row := rows[0].(map[string]interface{})

		if row["int_val"].(int64) != 42 {
			t.Fatalf("Expected int_val 42, got %v", row["int_val"])
		}

		if row["real_val"].(float64) != 3.14 {
			t.Fatalf("Expected real_val 3.14, got %v", row["real_val"])
		}

		if row["text_val"].(string) != "hello world" {
			t.Fatalf("Expected text_val 'hello world', got %v", row["text_val"])
		}

		closeFuncVal, _ := env.Get("dbClose")
		closeFunc := closeFuncVal.(r2core.BuiltinFunction)
		closeFunc(connId)
	})
}
