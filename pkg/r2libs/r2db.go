package r2libs

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"

	// Database drivers - imported anonymously
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// r2db.go: Database connectivity functions for R2Lang

// Global map to store database connections
var dbConnections = make(map[string]*sql.DB)

func RegisterDB(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"dbConnect": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("dbConnect needs (driver, dataSourceName)")
			}
			driver := toString(args[0])
			dsn := toString(args[1])

			// Validate driver
			supportedDrivers := []string{"sqlite3", "postgres", "mysql"}
			isSupported := false
			for _, d := range supportedDrivers {
				if d == driver {
					isSupported = true
					break
				}
			}
			if !isSupported {
				panic(fmt.Sprintf("dbConnect: unsupported driver '%s'. Supported: %v", driver, supportedDrivers))
			}

			db, err := sql.Open(driver, dsn)
			if err != nil {
				panic(fmt.Sprintf("dbConnect: failed to open database: %v", err))
			}

			err = db.Ping()
			if err != nil {
				db.Close()
				panic(fmt.Sprintf("dbConnect: failed to ping database: %v", err))
			}

			// Generate connection ID
			connId := fmt.Sprintf("conn_%d", len(dbConnections))
			dbConnections[connId] = db

			return connId
		}),

		"dbQuery": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("dbQuery needs (connectionId, query, ...args)")
			}
			connId := toString(args[0])
			query := toString(args[1])

			db, exists := dbConnections[connId]
			if !exists {
				panic(fmt.Sprintf("dbQuery: connection '%s' not found", connId))
			}

			// Prepare arguments for query
			queryArgs := make([]interface{}, len(args)-2)
			for i, arg := range args[2:] {
				queryArgs[i] = arg
			}

			rows, err := db.Query(query, queryArgs...)
			if err != nil {
				panic(fmt.Sprintf("dbQuery: %v", err))
			}
			defer rows.Close()

			// Get column names
			columns, err := rows.Columns()
			if err != nil {
				panic(fmt.Sprintf("dbQuery: failed to get columns: %v", err))
			}

			// Prepare result slice
			var results []interface{}

			for rows.Next() {
				// Create slice for scanning
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				err := rows.Scan(valuePtrs...)
				if err != nil {
					panic(fmt.Sprintf("dbQuery: failed to scan row: %v", err))
				}

				// Create map for this row
				rowMap := make(map[string]interface{})
				for i, col := range columns {
					// Convert []uint8 to string if needed
					if b, ok := values[i].([]uint8); ok {
						rowMap[col] = string(b)
					} else {
						rowMap[col] = values[i]
					}
				}

				results = append(results, rowMap)
			}

			if err = rows.Err(); err != nil {
				panic(fmt.Sprintf("dbQuery: row iteration error: %v", err))
			}

			return results
		}),

		"dbExec": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("dbExec needs (connectionId, query, ...args)")
			}
			connId := toString(args[0])
			query := toString(args[1])

			db, exists := dbConnections[connId]
			if !exists {
				panic(fmt.Sprintf("dbExec: connection '%s' not found", connId))
			}

			// Prepare arguments for query
			queryArgs := make([]interface{}, len(args)-2)
			for i, arg := range args[2:] {
				queryArgs[i] = arg
			}

			result, err := db.Exec(query, queryArgs...)
			if err != nil {
				panic(fmt.Sprintf("dbExec: %v", err))
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				panic(fmt.Sprintf("dbExec: failed to get rows affected: %v", err))
			}

			return rowsAffected
		}),

		"dbClose": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbClose needs (connectionId)")
			}
			connId := toString(args[0])

			db, exists := dbConnections[connId]
			if !exists {
				panic(fmt.Sprintf("dbClose: connection '%s' not found", connId))
			}

			err := db.Close()
			if err != nil {
				panic(fmt.Sprintf("dbClose: %v", err))
			}

			delete(dbConnections, connId)
			return true
		}),

		"dbBegin": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbBegin needs (connectionId)")
			}
			connId := toString(args[0])

			db, exists := dbConnections[connId]
			if !exists {
				panic(fmt.Sprintf("dbBegin: connection '%s' not found", connId))
			}

			_, err := db.Begin()
			if err != nil {
				panic(fmt.Sprintf("dbBegin: %v", err))
			}

			// For now, return a simple transaction ID
			// In a real implementation, you would store the *sql.Tx separately
			txId := fmt.Sprintf("tx_%s_%d", connId, len(dbConnections))
			return txId
		}),

		"dbLastInsertId": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("dbLastInsertId needs (connectionId, query, ...args)")
			}
			connId := toString(args[0])
			query := toString(args[1])

			db, exists := dbConnections[connId]
			if !exists {
				panic(fmt.Sprintf("dbLastInsertId: connection '%s' not found", connId))
			}

			// Prepare arguments for query
			queryArgs := make([]interface{}, len(args)-2)
			for i, arg := range args[2:] {
				queryArgs[i] = arg
			}

			result, err := db.Exec(query, queryArgs...)
			if err != nil {
				panic(fmt.Sprintf("dbLastInsertId: %v", err))
			}

			lastId, err := result.LastInsertId()
			if err != nil {
				panic(fmt.Sprintf("dbLastInsertId: failed to get last insert id: %v", err))
			}

			return lastId
		}),

		"dbPing": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbPing needs (connectionId)")
			}
			connId := toString(args[0])

			db, exists := dbConnections[connId]
			if !exists {
				panic(fmt.Sprintf("dbPing: connection '%s' not found", connId))
			}

			err := db.Ping()
			return err == nil
		}),

		"dbEscape": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbEscape needs (value)")
			}
			value := toString(args[0])
			// Simple SQL escape - replace single quotes with double single quotes
			return strings.ReplaceAll(value, "'", "''")
		}),

		"dbGetConnections": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			var connIds []interface{}
			for connId := range dbConnections {
				if !strings.HasPrefix(connId, "tx_") { // Exclude transaction IDs
					connIds = append(connIds, connId)
				}
			}
			return connIds
		}),
	}

	RegisterModule(env, "db", functions)
}

// Helper function to convert interface{} to string
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case int:
		return strconv.Itoa(val)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
