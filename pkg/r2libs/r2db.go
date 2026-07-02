package r2libs

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"

	// Database drivers - imported anonymously
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// r2db.go: Database connectivity functions for R2Lang

// dbConn bundles a live *sql.DB with the driver name it was opened with, so
// callers can adapt driver-specific behavior (e.g. postgres uses $1/$2/...
// placeholders instead of the ?  placeholders mysql/sqlite3 accept).
type dbConn struct {
	db     *sql.DB
	driver string
}

// Global map to store database connections. R2Lang scripts can call db
// builtins concurrently from multiple "r2"/goroutines, so access is guarded
// by dbConnectionsMu.
var (
	dbConnections   = make(map[string]*dbConn)
	dbConnectionsMu sync.RWMutex
	dbConnCounter   int
)

// adaptPlaceholders rewrites `?` positional placeholders into the
// driver-specific syntax a driver expects. mysql and sqlite3 accept `?`
// as-is; postgres (lib/pq) only recognizes `$1`, `$2`, ... and otherwise
// reports NumInput() == 0, which makes database/sql reject any query args
// with "sql: expected 0 arguments, got N". `?` inside single-quoted string
// literals is left untouched.
func adaptPlaceholders(driver, query string) string {
	if driver != "postgres" || !strings.Contains(query, "?") {
		return query
	}

	var sb strings.Builder
	sb.Grow(len(query) + 8)
	inString := false
	argN := 0
	for i := 0; i < len(query); i++ {
		c := query[i]
		switch {
		case c == '\'':
			inString = !inString
			sb.WriteByte(c)
		case c == '?' && !inString:
			argN++
			sb.WriteByte('$')
			sb.WriteString(strconv.Itoa(argN))
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// sqlValueToR2 converts a value scanned from a *sql.Rows into R2Lang's
// native value model (float64 for numbers, string, bool, nil, or
// *r2core.DateValue for dates/timestamps). Left unconverted, driver-returned
// types such as int64 or time.Time slip past R2Lang's arithmetic helpers
// (which only recognize float64/int/bool/string/nil), causing numeric
// operations on query results to either panic or silently fall back to
// string concatenation.
func sqlValueToR2(v interface{}) interface{} {
	switch val := v.(type) {
	case nil:
		return nil
	case []byte:
		return string(val)
	case int64:
		return float64(val)
	case int32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case time.Time:
		return r2core.NewDateValue(val)
	default:
		return val
	}
}

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

			dbConnectionsMu.Lock()
			// Use a monotonic counter, not len(dbConnections): once
			// connections are closed, len() can produce an id that
			// collides with one still in the map, silently overwriting
			// (and leaking) that live *sql.DB.
			dbConnCounter++
			connId := fmt.Sprintf("conn_%d", dbConnCounter)
			dbConnections[connId] = &dbConn{db: db, driver: driver}
			dbConnectionsMu.Unlock()

			return connId
		}),

		"dbQuery": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("dbQuery needs (connectionId, query, ...args)")
			}
			connId := toString(args[0])
			query := toString(args[1])

			dbConnectionsMu.RLock()
			conn, exists := dbConnections[connId]
			dbConnectionsMu.RUnlock()
			if !exists {
				panic(fmt.Sprintf("dbQuery: connection '%s' not found", connId))
			}

			// Prepare arguments for query
			queryArgs := make([]interface{}, len(args)-2)
			for i, arg := range args[2:] {
				queryArgs[i] = arg
			}

			rows, err := conn.db.Query(adaptPlaceholders(conn.driver, query), queryArgs...)
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
					rowMap[col] = sqlValueToR2(values[i])
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

			dbConnectionsMu.RLock()
			conn, exists := dbConnections[connId]
			dbConnectionsMu.RUnlock()
			if !exists {
				panic(fmt.Sprintf("dbExec: connection '%s' not found", connId))
			}

			// Prepare arguments for query
			queryArgs := make([]interface{}, len(args)-2)
			for i, arg := range args[2:] {
				queryArgs[i] = arg
			}

			result, err := conn.db.Exec(adaptPlaceholders(conn.driver, query), queryArgs...)
			if err != nil {
				panic(fmt.Sprintf("dbExec: %v", err))
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				panic(fmt.Sprintf("dbExec: failed to get rows affected: %v", err))
			}

			return sqlValueToR2(rowsAffected)
		}),

		"dbClose": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbClose needs (connectionId)")
			}
			connId := toString(args[0])

			dbConnectionsMu.RLock()
			conn, exists := dbConnections[connId]
			dbConnectionsMu.RUnlock()
			if !exists {
				panic(fmt.Sprintf("dbClose: connection '%s' not found", connId))
			}

			err := conn.db.Close()
			if err != nil {
				panic(fmt.Sprintf("dbClose: %v", err))
			}

			dbConnectionsMu.Lock()
			delete(dbConnections, connId)
			dbConnectionsMu.Unlock()
			return true
		}),

		"dbBegin": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbBegin needs (connectionId)")
			}
			connId := toString(args[0])

			dbConnectionsMu.RLock()
			conn, exists := dbConnections[connId]
			dbConnectionsMu.RUnlock()
			if !exists {
				panic(fmt.Sprintf("dbBegin: connection '%s' not found", connId))
			}

			tx, err := conn.db.Begin()
			if err != nil {
				panic(fmt.Sprintf("dbBegin: %v", err))
			}

			// The returned id is a placeholder (there's no dbCommit/dbRollback
			// yet to reference it), so roll back right away instead of
			// leaking the pooled connection the Tx holds onto forever.
			if err := tx.Rollback(); err != nil {
				panic(fmt.Sprintf("dbBegin: %v", err))
			}

			dbConnectionsMu.RLock()
			n := len(dbConnections)
			dbConnectionsMu.RUnlock()
			txId := fmt.Sprintf("tx_%s_%d", connId, n)
			return txId
		}),

		"dbLastInsertId": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("dbLastInsertId needs (connectionId, query, ...args)")
			}
			connId := toString(args[0])
			query := toString(args[1])

			dbConnectionsMu.RLock()
			conn, exists := dbConnections[connId]
			dbConnectionsMu.RUnlock()
			if !exists {
				panic(fmt.Sprintf("dbLastInsertId: connection '%s' not found", connId))
			}

			// Prepare arguments for query
			queryArgs := make([]interface{}, len(args)-2)
			for i, arg := range args[2:] {
				queryArgs[i] = arg
			}

			result, err := conn.db.Exec(adaptPlaceholders(conn.driver, query), queryArgs...)
			if err != nil {
				panic(fmt.Sprintf("dbLastInsertId: %v", err))
			}

			lastId, err := result.LastInsertId()
			if err != nil {
				panic(fmt.Sprintf("dbLastInsertId: failed to get last insert id: %v", err))
			}

			return sqlValueToR2(lastId)
		}),

		"dbPing": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("dbPing needs (connectionId)")
			}
			connId := toString(args[0])

			dbConnectionsMu.RLock()
			conn, exists := dbConnections[connId]
			dbConnectionsMu.RUnlock()
			if !exists {
				panic(fmt.Sprintf("dbPing: connection '%s' not found", connId))
			}

			err := conn.db.Ping()
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
			dbConnectionsMu.RLock()
			defer dbConnectionsMu.RUnlock()

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
